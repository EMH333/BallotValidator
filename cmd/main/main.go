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

	"ethohampton.com/BallotCleaner/internal/steps"
	"ethohampton.com/BallotCleaner/internal/util"
)

var validVotersGraduate []string
var validVotersUndergrad []string
var validVotersUndefined []string

var alreadyVotedPrevious []string

const numToPick int = 10 // how many winners to pick //TODO tie into command line

func main() {
	var startDay int64 = 0      // what day are we starting on to process votes
	var endDay int64 = startDay // what day are we ending on to process votes
	var dataFile string = "ballotData.csv"

	// in the form of `program <day> <file_to_process>`
	if len(os.Args) == 3 {
		day, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			log.Fatal("Couldn't parse argument")
		}
		startDay = day
		endDay = day
		dataFile = os.Args[2]
	}

	// in the form of `program <start_day> <end_day> <file_to_process>`
	if len(os.Args) == 4 {
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
		dataFile = os.Args[3]
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
	votes := util.LoadVotesCSV("data/ballots/"+dataFile, startDay, endDay)
	log.Printf("%d votes loaded for day %d through %d\n", len(votes), startDay, endDay)

	// step one: valid voter
	log.Println()
	log.Println("Step 1: Valid voter")
	validPostOne, invalidPostOne, oneSummary := steps.StepOne(votes, &validVotersGraduate, &validVotersUndergrad, &validVotersUndefined)
	storeVotes(validPostOne, "1-valid-"+dayToDayFormat+".csv")
	storeVotes(invalidPostOne, "1-invalid-"+dayToDayFormat+".csv")
	storeSummary(oneSummary, "1-summary-"+dayToDayFormat+".txt")
	log.Println("Step 1: Invalid votes:", oneSummary.Invalid)
	log.Println("Step 1: Valid votes:", oneSummary.Valid)

	// step two: valid voter
	log.Println()
	log.Println("Step 2: Dedupe")
	validPostTwo, invalidPostTwo, alreadyVotedToday, twoSummary := steps.StepTwo(validPostOne, &alreadyVotedPrevious)
	storeVotes(validPostTwo, "2-valid-"+dayToDayFormat+".csv")
	storeVotes(invalidPostTwo, "2-invalid-"+dayToDayFormat+".csv")
	storeSummary(twoSummary, "2-summary-"+dayToDayFormat+".txt")
	storeAlreadyVoted(alreadyVotedToday, "alreadyVoted-"+dayToDayFormat+".csv")
	log.Println("Step 2: Invalid votes:", twoSummary.Invalid)
	log.Println("Step 2: Valid votes:", twoSummary.Valid)

	// step three: grad/undergrad
	log.Println()
	log.Println("Step 3: Grad/undergrad")
	validPostThree, invalidPostThree, threeSummary := steps.StepThree(validPostTwo, &validVotersGraduate, &validVotersUndergrad)
	storeVotes(validPostThree, "3-valid-"+dayToDayFormat+".csv")
	storeVotes(invalidPostThree, "3-modified-"+dayToDayFormat+".csv")
	storeSummary(threeSummary, "3-summary-"+dayToDayFormat+".txt")
	log.Println("Step 3: Modified votes:", threeSummary.Invalid)
	log.Println("Step 3: Valid votes:", threeSummary.Valid)

	// step four: Incentives
	log.Println()
	log.Println("Step 4: Incentives")
	postFour, winners, fourSummary := steps.StepFour(validPostThree, seed, numToPick)
	storeVotes(postFour, "4-valid-"+dayToDayFormat+".csv")
	storeSummary(fourSummary, "4-summary-"+dayToDayFormat+".txt")
	storeAlreadyVoted(winners, "incentive-winners-"+dayToDayFormat+".csv")
	log.Println("Step 4: Valid votes:", threeSummary.Valid)
	log.Println("Step 4: Selected winners:", len(winners))

	//experimental
	steps.StepFourtyTwo(postFour, "output/results")
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

	alreadyVoted = util.RemoveDuplicateStr(alreadyVoted) //make sure we don't have any duplicates (though it doesn't really matter)

	return alreadyVoted
}

func storeVotes(votes []util.Vote, filename string) {
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

func storeSummary(summary util.Summary, filename string) {
	f, err := os.Create("output/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()

	_, err = f.WriteString(summary.StepInfo + "\n")
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(fmt.Sprintf("Processed: %d\n", summary.Processed))
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(fmt.Sprintf("Valid: %d\n", summary.Valid))
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(fmt.Sprintf("Invalid: %d\n", summary.Invalid))
	if err != nil {
		log.Fatal(err)
	}
	if len(summary.Log) != 0 {
		_, err = f.WriteString("\n\nLog Messages:\n")
		if err != nil {
			log.Fatal(err)
		}
		for _, message := range summary.Log {
			_, err = f.WriteString(message + "\n")
			if err != nil {
				log.Fatal(err)
			}
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
		_, err = f.WriteString(record + "\n")
		if err != nil {
			log.Fatal(err)
		}
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
