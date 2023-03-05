package util

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var EPOCH, epochErr = time.Parse("2006-Jan-02 03:04:05", "2023-Feb-20 00:00:01")

//TODO set correct import values
// values to use when importing from csv
const IMPORT_TIMESTAMP = 1 //using end date so it is consistent across submission times
const IMPORT_TYPE = 2
const IMPORT_ONID = 50
const IMPORT_COMPLETE = 6
const IMPORT_ID = 8

func LoadVotesCSV(fileName string, startDay, endDay int64) []Vote {
	// make sure our epoch is valid
	if epochErr != nil {
		log.Fatal(epochErr)
	}

	var validStartTime = EPOCH.Add(time.Duration(startDay) * 24 * time.Hour)
	var validEndTime = EPOCH.Add(time.Duration(endDay+1) * 24 * time.Hour) // add one day to end day

	//TODO verify this is correct
	if endDay == 18 {
		newEndTime, err := time.Parse("2006-Jan-02 15:04:05", "2023-Mar-10 17:00:59")
		if err != nil {
			log.Fatal(err)
		}
		validEndTime = newEndTime
	}

	var votes []Vote

	//load csv file
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// read csv values using csv.Reader
	//with modifications to handle the specifics of the valid votes list
	csvReader := csv.NewReader(f)
	csvReader.Comma = ','
	csvReader.TrimLeadingSpace = true

	var incompleteVotes int = 0

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		//skip the first few rows
		if rec[IMPORT_TIMESTAMP] == "EndDate" || rec[IMPORT_TIMESTAMP] == "End Date" || strings.Contains(rec[IMPORT_TIMESTAMP], "ImportId") {
			continue
		}

		if rec[IMPORT_TYPE] == "Survey Preview" {
			log.Println("Skipping survey preview response")
			continue
		}

		timestamp, err := time.Parse("2006-01-02 15:04:05", rec[IMPORT_TIMESTAMP]) //"1/2/2006 15:04" //2/14/2022 9:10
		if err != nil {
			log.Fatal(err)
		}

		//make sure it is only reading the correct day
		if timestamp.Before(validStartTime) || timestamp.After(validEndTime) {
			continue
		}

		ONID := rec[IMPORT_ONID]
		//sanity check to make sure the ONID looks like an email
		if !strings.Contains(ONID, "@oregonstate.edu") {
			log.Fatalf("ONID is not an email address: %s\n", ONID)
		}

		if strings.Contains(strings.Split(ONID, "@")[0], ".") {
			log.Fatalf("ONID should not contain a dot: %s\n", ONID)
		}

		//make sure it is a complete row
		if strings.ToUpper(rec[IMPORT_COMPLETE]) != "TRUE" {
			//log.Printf("Vote is not complete: %+v\n", rec)
			log.Printf("Vote is not complete from %s: %+v\n", rec[IMPORT_ONID], rec[0:IMPORT_COMPLETE+2])
			incompleteVotes++
			continue
		}

		id := rec[IMPORT_ID]
		if !strings.HasPrefix(rec[IMPORT_ID], "R_") {
			log.Fatalf("Response ID is not valid: %+v\n", rec)
		}

		//append rec to votes
		votes = append(votes, Vote{Raw: rec, Timestamp: timestamp, ONID: ONID, ID: id})
	}

	log.Printf("%d votes were incomplete, and not counted\n", incompleteVotes)

	return votes
}

const VALID_ONID_EMAIL = 2

func LoadValidVoters(fileName string) []string {
	var voters []string

	//open csv file
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//with modifications to handle the specifics of the valid votes list
	csvReader := csv.NewReader(f)
	//csvReader.Comma = '\t'
	csvReader.TrimLeadingSpace = true

	//skip the first row
	_, err = csvReader.Read()
	if err != nil {
		log.Fatal(err)
	}

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		voters = append(voters, rec[VALID_ONID_EMAIL])

	}

	return voters
}

func LoadAlreadyVoted(folderName string, upToDay int64) []string {
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

	alreadyVoted = RemoveDuplicateStr(alreadyVoted) //make sure we don't have any duplicates (though it doesn't really matter)

	return alreadyVoted
}

func StoreVotes(votes []Vote, filename string) {
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

func StoreSummary(summary Summary, filename string) {
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

// simple store list of onids
func StoreAlreadyVoted(alreadyVoted []string, filename string) {
	f, err := os.Create("output/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()

	//sort the alreadyVoted slice
	sort.Strings(alreadyVoted)

	for _, record := range alreadyVoted {
		_, err = f.WriteString(record + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

// store in similar format to valid voters, but with fake data other than the onid email
func StoreNotYetVoted(notYetVoted []string, filename string) {
	f, err := os.Create("output/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()

	//sort the notYetVoted slice
	sort.Strings(notYetVoted)

	_, err = f.WriteString("First Name,Last Name,Email,ONID\n")
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range notYetVoted {
		_, err = f.WriteString("OSU,Student," + record + ",osustudent\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func LoadSeed() string {
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

func LoadStringArrayFile(fileName string) []string {
	var strings []string
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Error opening file", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strings = append(strings, scanner.Text())
	}

	return strings
}
