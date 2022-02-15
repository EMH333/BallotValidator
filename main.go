package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"time"
)

// Raw is the raw row
// Timestamp is the timestamp of the row (parsed)
// ONID is the ONID of the row (taken from the raw)
type Vote struct {
	Raw       []string
	Timestamp time.Time
	ONID      string
}

type Summary struct {
	processed int
	valid     int
	invalid   int
}

var validVotersGraduate []string
var validVotersUndergrad []string
var validVotersUndefined []string

var alreadyVotedPrevious []string
var alreadyVotedToday []string

func main() {
	var startDay int // what day are we starting on to process votes
	var endDay int   // what day are we ending on to process votes

	log.Printf("Selected start day: %d, Selected end day: %d\n", startDay, endDay)

	// Load the valid voters
	log.Println("Loading valid voters...")
	validVotersGraduate = loadValidVoters("data/validVoters.csv", "G")
	validVotersUndergrad = loadValidVoters("data/validVoters.csv", "UG")
	validVotersUndefined = loadValidVoters("data/validVoters.csv", "Self Identified on Ballot")

	// Load the already voted
	log.Printf("Loading already voted up to day %d...\n", startDay)
	alreadyVotedPrevious = loadAlreadyVoted("TODO: folder", startDay)
	alreadyVotedToday = make([]string, 0, 100)

	// Load the votes
	log.Println("Loading votes...")
	votes := loadVotesCSV("TODO: filename")

	// step one: valid voter
	log.Println("Step 1: Valid voter")
	validPostOne, invalidPostOne, oneSummary := stepOne(votes, &validVotersGraduate, &validVotersUndergrad, &validVotersUndefined)
	storeVotes(validPostOne, "TODO: filename")
	storeVotes(invalidPostOne, "TODO: filename")
	storeSummary(oneSummary, "TODO: filename")

}

const VALID_STATUS = 4
const VALID_ONID_EMAIL = 2

func loadValidVoters(fileName string, indicator string) []string {
	var voters []string
	//return []string{"TODO"}
	//load csv file
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// read csv values using csv.Reader
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

func loadAlreadyVoted(folderName string, upToDay int) []string {
	return []string{"TODO"}
}

func storeVotes(votes []Vote, filename string) {
	//store the vote.raw in csv format under filename
	//TODO
}

func storeSummary(summary Summary, filename string) {
	//store the summary
	//TODO
}
