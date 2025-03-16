package util

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ELECTION_TIMEZONE, timezoneErr = time.LoadLocation("America/Los_Angeles")

// TODO use this to also load new votes csv (add ONID and logging options)
func LoadVotesCSV(countingConfig *CountingConfig, fileName string, startDay, endDay int) []Vote {
	// make sure our timezone and epoch are valid
	if timezoneErr != nil {
		log.Fatal(timezoneErr)
	}
	var EPOCH, epochErr = time.ParseInLocation("2006-Jan-02 03:04:05", countingConfig.ElectionEpoch, ELECTION_TIMEZONE)
	if epochErr != nil {
		log.Fatal(epochErr)
	}

	var validStartTime = EPOCH.Add(time.Duration(startDay) * 24 * time.Hour)
	var validEndTime = EPOCH.Add(time.Duration(endDay+1) * 24 * time.Hour) // add one day to end day

	// validate the start and end time
	if startDay == 0 && endDay == countingConfig.ElectionNumDays {
		newStartTime, err := time.ParseInLocation(countingConfig.BallotTimeFormat, countingConfig.ElectionStartTime, ELECTION_TIMEZONE)
		if err != nil {
			log.Fatal(err)
		}
		validStartTime = newStartTime

		newEndTime, err := time.ParseInLocation(countingConfig.BallotTimeFormat, countingConfig.ElectionEndTime, ELECTION_TIMEZONE)
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

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for idx, rec := range records {
		if idx == 1 {
			err := candidateColumnIndexValidation(countingConfig, rec)
			if err != nil {
				log.Fatal(err)
			}
		}
		//skip the first few rows which are headers
		if rec[countingConfig.ImportTimestamp] == "EndDate" || rec[countingConfig.ImportTimestamp] == "End Date" || strings.Contains(rec[countingConfig.ImportTimestamp], "ImportId") {
			continue
		}
		if rec[countingConfig.ImportTimestamp] == "StartDate" || rec[countingConfig.ImportTimestamp] == "Start Date" {
			continue
		}

		if rec[countingConfig.ImportType] == "Survey Preview" {
			log.Println("Skipping survey preview response")
			continue
		}

		timestamp, err := time.ParseInLocation(countingConfig.BallotTimeFormat, rec[countingConfig.ImportTimestamp], ELECTION_TIMEZONE) //"1/2/2006 15:04" //2/14/2022 9:10
		if err != nil {
			log.Fatal(err)
		}

		//make sure it is only reading the correct day
		if timestamp.Before(validStartTime) || timestamp.After(validEndTime) {
			//log.Printf("Response before or after valid times: %+v\n", rec)
			outOfTimeVotes++
			continue
		}

		ONID := rec[countingConfig.ImportONID]
		//sanity check to make sure the ONID looks like an email
		if !strings.Contains(ONID, "@oregonstate.edu") {
			log.Fatalf("ONID is not an email address: %s, vote id: %s\n", ONID, rec[countingConfig.ImportID])
		}

		if strings.Contains(strings.Split(ONID, "@")[0], ".") {
			log.Fatalf("ONID should not contain a dot: %s\n", ONID)
		}

		//make sure it is a complete row
		if !strings.EqualFold(rec[countingConfig.ImportComplete], "TRUE") {
			//log.Printf("Vote is not complete: %+v\n", rec)
			log.Printf("Vote is not complete from %s: %+v\n", rec[countingConfig.ImportONID], rec[0:countingConfig.ImportComplete+2])
			incompleteVotes++
			continue
		}

		id := rec[countingConfig.ImportID]
		if !strings.HasPrefix(rec[countingConfig.ImportID], "R_") {
			log.Fatalf("Response ID is not valid: %+v\n", rec)
		}

		//append rec to votes
		votes = append(votes, Vote{Raw: rec, Timestamp: timestamp, ONID: ONID, ID: id})
	}

	log.Printf("%d votes were incomplete, and not counted\n", incompleteVotes)
	log.Printf("%d votes were out of time, and not counted\n", outOfTimeVotes)

	return votes
}

var ErrCandidateOrder = errors.New("candidates in wrong order")

func CandidateOrderError(got, expected string) error {
	return fmt.Errorf("candidates in wrong order %w. Got '%s', expected '%s'", ErrCandidateOrder, got, expected)
}

// Make sure the candidate names for president/vice-president and sfc-chair are in the right order
func candidateColumnIndexValidation(countingConfig *CountingConfig, rec []string) error {
	// Validate president/vice-president order
	initialPresidentOffset := countingConfig.TallyPresidentOptionsIndex
	for idx, can := range countingConfig.CandidatesPresident {
		offset := initialPresidentOffset + idx
		entry := rec[offset]
		if !strings.HasSuffix(strings.TrimSpace(entry), can) {
			return CandidateOrderError(entry, can)
		}
	}

	initialSfcOffset := countingConfig.TallySFCChairOptionsIndex
	for idx, can := range countingConfig.CandidatesSFCChair {
		offset := initialSfcOffset + idx
		entry := rec[offset]
		if !strings.HasSuffix(strings.TrimSpace(entry), can) {
			return CandidateOrderError(entry, can)
		}
	}

	// No errors found
	return nil
}

func LoadValidVoters(countingConfig *CountingConfig, indicator string) []string {
	var voters []string

	//open csv file
	f, err := os.Open(countingConfig.ValidVotersFile)
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
	if err != nil || strings.Contains(first[countingConfig.ValidVotersEmailIndex], "@") {
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
		if rec[countingConfig.ValidVotersStatusIndex] == indicator {
			//confirm it is an email
			if !strings.Contains(rec[countingConfig.ValidVotersEmailIndex], "@") {
				log.Fatalf("ONID is not an email address: %s\n", rec[countingConfig.ValidVotersEmailIndex])
			}

			voters = append(voters, rec[countingConfig.ValidVotersEmailIndex])
		}

	}

	return voters
}

func LoadAlreadyVoted(countingConfig *CountingConfig, upToDay int64) []string {
	var alreadyVoted []string

	folder := countingConfig.AlreadyVotedDir

	//make sure folder ends with a slash
	if !strings.HasSuffix(folder, "/") {
		folder += "/"
	}

	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		log.Fatalf("%s doesn't exist", folder)
	}

	files, err := os.ReadDir(folder)
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
	_, err = fmt.Fprintf(f, "Processed: %d\n", summary.Processed)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(f, "Valid: %d\n", summary.Valid)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(f, "Invalid: %d\n", summary.Invalid)
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
func StoreStringArrayFile(onids []string, filename string, sortList bool) {
	f, err := os.Create("output/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()

	if sortList {
		//sort the alreadyVoted slice
		sort.Strings(onids)
	}

	for _, record := range onids {
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
