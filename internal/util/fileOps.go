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

var ELECTION_TIMEZONE, timezoneErr = time.LoadLocation("America/Los_Angeles")
var EPOCH, epochErr = time.ParseInLocation("2006-Jan-02 03:04:05", "2024-Feb-19 00:00:01", ELECTION_TIMEZONE)

const ELECTION_NUM_DAYS = 18
const ELECTION_START_TIME = "2/19/2024 12:00"
const ELECTION_END_TIME = "3/1/2024 12:00"

const BALLOT_TIME_FORMAT = "1/2/2006 15:04"

// values to use when importing from csv
const IMPORT_TIMESTAMP = 1 //using end date so it is consistent across submission times
const IMPORT_TYPE = 2
const IMPORT_ONID = 39
const IMPORT_COMPLETE = 6
const IMPORT_ID = 8

// TODO use this to also load new votes csv (add ONID and logging options)
func LoadVotesCSV(fileName string, startDay, endDay, ONIDIndex int64) []Vote {
	// make sure our timezone and epoch are valid
	if timezoneErr != nil {
		log.Fatal(timezoneErr)
	}
	if epochErr != nil {
		log.Fatal(epochErr)
	}

	var validStartTime = EPOCH.Add(time.Duration(startDay) * 24 * time.Hour)
	var validEndTime = EPOCH.Add(time.Duration(endDay+1) * 24 * time.Hour) // add one day to end day

	// validate the start and end time
	if startDay == 0 && endDay == ELECTION_NUM_DAYS {
		newStartTime, err := time.ParseInLocation(BALLOT_TIME_FORMAT, ELECTION_START_TIME, ELECTION_TIMEZONE)
		if err != nil {
			log.Fatal(err)
		}
		validStartTime = newStartTime

		newEndTime, err := time.ParseInLocation(BALLOT_TIME_FORMAT, ELECTION_END_TIME, ELECTION_TIMEZONE)
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
	var outOfTimeVotes int = 0

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		//skip the first few rows which are headers
		if rec[IMPORT_TIMESTAMP] == "EndDate" || rec[IMPORT_TIMESTAMP] == "End Date" || strings.Contains(rec[IMPORT_TIMESTAMP], "ImportId") {
			continue
		}

		if rec[IMPORT_TYPE] == "Survey Preview" {
			log.Println("Skipping survey preview response")
			continue
		}

		timestamp, err := time.ParseInLocation(BALLOT_TIME_FORMAT, rec[IMPORT_TIMESTAMP], ELECTION_TIMEZONE) //"1/2/2006 15:04" //2/14/2022 9:10
		if err != nil {
			log.Fatal(err)
		}

		//make sure it is only reading the correct day
		if timestamp.Before(validStartTime) || timestamp.After(validEndTime) {
			//log.Printf("Response before or after valid times: %+v\n", rec)
			outOfTimeVotes++
			continue
		}

		ONID := rec[ONIDIndex]
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
			log.Printf("Vote is not complete from %s: %+v\n", rec[ONIDIndex], rec[0:IMPORT_COMPLETE+2])
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
	log.Printf("%d votes were out of time, and not counted\n", outOfTimeVotes)

	return votes
}

const VALID_ONID_EMAIL = 2
const VALID_STATUS = 4 //G_UG_STATUS TODO

func LoadValidVoters(fileName string, indicator string) []string {
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

	//skip the first row which is headers
	// and check it doesn't contain an @ (which could be an email)
	first, err := csvReader.Read()
	if err != nil || strings.Contains(first[VALID_ONID_EMAIL], "@") {
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

		//check and see if the indicator (G_UG_STATUS) is valid for who we are trying to process
		if rec[VALID_STATUS] == indicator {
			//confirm it is an email
			if !strings.Contains(rec[VALID_ONID_EMAIL], "@") {
				log.Fatalf("ONID is not an email address: %s\n", rec[VALID_ONID_EMAIL])
			}

			voters = append(voters, rec[VALID_ONID_EMAIL])
		}

	}

	return voters
}

func LoadAlreadyVoted(folder string, upToDay int64) []string {
	var alreadyVoted []string

	//make sure folder ends with a slash
	if !strings.HasSuffix(folder, "/") {
		folder += "/"
	}

	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		log.Fatalf("%s doesn't exist", folder)
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		// note we are grabbing the end date instead of the start date since that guarantees no overlap
		// most files will have the same start and end, so this likely doesn't matter
		//expecting alreadyVoted-<day>-<endDay>.csv
		fileDay, err := strconv.ParseInt(strings.TrimSuffix(strings.Split(file.Name(), "-")[2], ".csv"), 10, 64)
		if err != nil {
			log.Fatalln("Already voted file name formated incorrectly", err)
		}

		// ignore files that are past today since they won't yield helpful results
		if fileDay >= upToDay {
			continue
		}

		alreadyVoted = append(alreadyVoted, LoadStringArrayFile(folder+file.Name())...)
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
func StoreStringArrayFile(alreadyVoted []string, filename string) {
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
