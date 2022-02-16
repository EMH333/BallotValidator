package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Raw is the raw row
// Timestamp is the timestamp of the row (parsed)
// ONID is the ONID of the row (taken from the raw)
type Vote struct {
	Raw       []string
	Timestamp time.Time
	ONID      string
	ID        string
}

type Summary struct {
	processed int
	valid     int
	invalid   int
	log       []string
}

var validVotersGraduate []string
var validVotersUndergrad []string
var validVotersUndefined []string

var alreadyVotedPrevious []string

const numToPick int = 10 // how many winners to pick //TODO tie into command line

func main() {
	var startDay int64 = 0      // what day are we starting on to process votes
	var endDay int64 = startDay // what day are we ending on to process votes

	// in the form of `program <day>`
	if len(os.Args) == 2 {
		day, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			log.Fatal("Couldn't parse argument")
		}
		startDay = day
		endDay = day
	}

	// in the form of `program <start_day> <end_day>`
	if len(os.Args) == 3 {
		day1, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			log.Fatal("Couldn't parse argument 1")
		}

		day2, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			log.Fatal("Couldn't parse argument 2")
		}

		startDay = day1
		endDay = day2
	}

	log.Printf("Selected start day: %d, Selected end day: %d\n", startDay, endDay)
	var dayToDayFormat = fmt.Sprint(startDay) + "-" + fmt.Sprint(endDay)

	//"random" seed so winners are deterministic
	var seed = loadSeed() + "-" + dayToDayFormat //include days in picking so it is unique

	_, err := os.Stat("output")
	if os.IsNotExist(err) && os.Mkdir("output", 0755) != nil {
		log.Fatal("Could not create output directory")
	}

	// Load the valid voters
	log.Println("Loading valid voters...")
	validVotersGraduate = loadValidVoters("data/validVoters.csv", "G")
	validVotersUndergrad = loadValidVoters("data/validVoters.csv", "UG")
	validVotersUndefined = loadValidVoters("data/validVoters.csv", "Self Identified on Ballot")

	log.Printf("There are %d valid voters for graduate students\n", len(validVotersGraduate))
	log.Printf("There are %d valid voters for undergrad students\n", len(validVotersUndergrad))
	log.Printf("There are %d valid voters for undefined students\n", len(validVotersUndefined))

	// Load the already voted
	log.Printf("Loading already voted up to day %d...\n", startDay)
	alreadyVotedPrevious = loadAlreadyVoted("alreadyVoted", int64(startDay))

	log.Printf("%d students have already voted\n", len(alreadyVotedPrevious))

	// Load the votes
	log.Println("Loading votes...")
	votes := loadVotesCSV("data/ballots/0-Feb-14-complete.csv", startDay, endDay) //TODO allow for flexibility in the filename via command line args
	log.Printf("%d votes loaded for day %d through %d\n", len(votes), startDay, endDay)

	// step one: valid voter
	log.Println()
	log.Println("Step 1: Valid voter")
	validPostOne, invalidPostOne, oneSummary := stepOne(votes, &validVotersGraduate, &validVotersUndergrad, &validVotersUndefined)
	storeVotes(validPostOne, "1-valid-"+dayToDayFormat+".csv")
	storeVotes(invalidPostOne, "1-invalid-"+dayToDayFormat+".csv")
	storeSummary(oneSummary, "1-summary-"+dayToDayFormat+".txt")
	log.Println("Step 1: Invalid votes:", oneSummary.invalid)
	log.Println("Step 1: Valid votes:", oneSummary.valid)

	// step two: valid voter
	log.Println()
	log.Println("Step 2: Dedupe")
	validPostTwo, invalidPostTwo, alreadyVotedToday, twoSummary := stepTwo(validPostOne, &alreadyVotedPrevious)
	storeVotes(validPostTwo, "2-valid-"+dayToDayFormat+".csv")
	storeVotes(invalidPostTwo, "2-invalid-"+dayToDayFormat+".csv")
	storeSummary(twoSummary, "2-summary-"+dayToDayFormat+".txt")
	storeAlreadyVoted(alreadyVotedToday, "alreadyVoted-"+dayToDayFormat+".csv")
	log.Println("Step 2: Invalid votes:", twoSummary.invalid)
	log.Println("Step 2: Valid votes:", twoSummary.valid)

	// step three: grad/undergrad
	log.Println()
	log.Println("Step 3: Grad/undergrad")
	validPostThree, invalidPostThree, threeSummary := stepThree(validPostTwo, &validVotersGraduate, &validVotersUndergrad)
	storeVotes(validPostThree, "3-valid-"+dayToDayFormat+".csv")
	storeVotes(invalidPostThree, "3-modified-"+dayToDayFormat+".csv")
	storeSummary(threeSummary, "3-summary-"+dayToDayFormat+".txt")
	log.Println("Step 3: Modified votes:", threeSummary.invalid)
	log.Println("Step 3: Valid votes:", threeSummary.valid)

	// step four: Incentives
	log.Println()
	log.Println("Step 4: Incentives")
	postFour, winners, fourSummary := stepFour(validPostThree, seed, numToPick)
	storeVotes(postFour, "4-valid-"+dayToDayFormat+".csv")
	storeSummary(fourSummary, "4-summary-"+dayToDayFormat+".txt")
	storeAlreadyVoted(winners, "winners-"+dayToDayFormat+".csv")
	log.Println("Step 4: Valid votes:", threeSummary.valid)
}

const VALID_STATUS = 4
const VALID_ONID_EMAIL = 2

func loadValidVoters(fileName string, indicator string) []string {
	var voters []string

	//open csv file
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//with modifications to handle the specifics of the valid votes list
	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	csvReader.TrimLeadingSpace = true

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		//check and see if the indicator (G_UG_STATUS) is valid for who we are trying to process
		if rec[VALID_STATUS] == indicator {
			voters = append(voters, rec[VALID_ONID_EMAIL])
		}
	}

	return voters
}

func loadAlreadyVoted(folderName string, upToDay int64) []string {
	var alreadyVoted []string

	var folder = "data/" + folderName + "/"

	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		log.Fatal("alreadyVoted doesn't exist within data")
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileDay, err := strconv.ParseInt(strings.Split(file.Name(), "-")[1], 10, 64) //expecting alreadyVoted-<day>-<endDay>.csv
		if err != nil {
			log.Fatalln("Already voted file name formated incorrectly", err)
		}

		// ignore files that are past today since they won't yield helpful results
		if fileDay >= upToDay {
			continue
		}

		file, err := os.Open(folder + file.Name())
		if err != nil {
			log.Fatalln("Error opening already voted file", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			alreadyVoted = append(alreadyVoted, scanner.Text())
		}
	}

	alreadyVoted = removeDuplicateStr(alreadyVoted) //make sure we don't have any duplicates (though it doesn't really matter)

	return alreadyVoted
}

func storeVotes(votes []Vote, filename string) {
	//store the vote.raw in csv format under filename
	f, err := os.Create("output/" + filename)
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range votes {
		if err := w.Write(record.Raw); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

	w.Flush() // make sure we flush before closing file
}

func storeSummary(summary Summary, filename string) {
	f, err := os.Create("output/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()

	f.WriteString("Step 1: Valid voter\n")
	f.WriteString(fmt.Sprintf("Processed: %d\n", summary.processed))
	f.WriteString(fmt.Sprintf("Valid: %d\n", summary.valid))
	f.WriteString(fmt.Sprintf("Invalid: %d\n", summary.invalid))
	if len(summary.log) != 0 {
		f.WriteString("\n\nLog Messages:\n")
		for _, message := range summary.log {
			f.WriteString(message + "\n")
		}
	}
}

func storeAlreadyVoted(alreadyVoted []string, filename string) {
	f, err := os.Create("output/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()

	for _, record := range alreadyVoted {
		f.WriteString(record + "\n")
	}
}

func loadSeed() string {
	//load the seed from the seed.txt file
	f, err := os.Open("data/seed.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	seed := scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	if seed == "" {
		log.Fatal("seed.txt is empty")
	}
	return seed
}
