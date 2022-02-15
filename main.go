package main

import (
	"log"
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
var alreadyVotedPrevious []string
var alreadyVotedToday []string

func main() {
	var startDay int // what day are we starting on to process votes
	var endDay int   // what day are we ending on to process votes

	log.Printf("Selected start day: %d, Selected end day: %d\n", startDay, endDay)

	// Load the valid voters
	log.Println("Loading valid voters...")
	validVotersGraduate = loadValidVoters("TODO")
	validVotersUndergrad = loadValidVoters("TODO")

	// Load the already voted
	log.Printf("Loading already voted up to day %d...\n", startDay)
	alreadyVotedPrevious = loadAlreadyVoted("TODO: folder", startDay)
	alreadyVotedToday = make([]string, 0, 100)

	// Load the votes
	log.Println("Loading votes...")
	votes := loadVotesCSV("TODO: filename")

	// step one: valid voter
	log.Println("Step 1: Valid voter")
	validPostOne, invalidPostOne, oneSummary := stepOne(votes, validVotersGraduate, validVotersUndergrad)
	storeVotes(validPostOne, "TODO: filename")
	storeVotes(invalidPostOne, "TODO: filename")
	storeSummary(oneSummary, "TODO: filename")

}

func loadValidVoters(fileName string) []string {
	return []string{"TODO"}
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
