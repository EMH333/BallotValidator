package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const TALLY_MEASURE = 17
const TALLY_SENATE_OPTIONS = 18
const TALLY_SENATE_WRITEINS = 6
const TALLY_SFCATLARGE_OPTIONS = TALLY_SENATE_OPTIONS + TALLY_SENATE_WRITEINS + 1
const TALLY_SFCATLARGE_WRITEINS = 4
const TALLY_UNDERGRADREPS_OPTIONS = TALLY_SFCATLARGE_OPTIONS + TALLY_SFCATLARGE_WRITEINS + 2 // one more for the additional grad/undergrad question
const TALLY_UNDERGRADREPS_WRITEINS = 20
const TALLY_GRADREPS_OPTIONS = TALLY_UNDERGRADREPS_OPTIONS + TALLY_UNDERGRADREPS_WRITEINS + 1
const TALLY_GRADREPS_WRITEINS = 5

// designed to do all the counting and output a nice little summary
func stepFourtyTwo(votes []Vote, outputDirname string) {
	var ballotYes int = 0
	var ballotNo int = 0

	var senate map[string]int = make(map[string]int)
	var sfcAtLarge map[string]int = make(map[string]int)
	var undergradReps map[string]int = make(map[string]int)
	var gradReps map[string]int = make(map[string]int)

	for _, vote := range votes {
		if vote.Raw[TALLY_MEASURE] == "Yes" {
			ballotYes++
		} else if vote.Raw[TALLY_MEASURE] == "No" {
			ballotNo++
		}

		///////////////////SENATE/////////////////////////////
		countPopularityVote(&vote, &senate, TALLY_SENATE_OPTIONS, TALLY_SENATE_WRITEINS)

		///////////////////SFC AT LARGE/////////////////////////////
		countPopularityVote(&vote, &sfcAtLarge, TALLY_SFCATLARGE_OPTIONS, TALLY_SFCATLARGE_WRITEINS)

		///////////////////UNDERGRAD REPS/////////////////////////////
		countPopularityVote(&vote, &undergradReps, TALLY_UNDERGRADREPS_OPTIONS, TALLY_UNDERGRADREPS_WRITEINS)

		///////////////////GRAD REPS/////////////////////////////
		countPopularityVote(&vote, &gradReps, TALLY_GRADREPS_OPTIONS, TALLY_GRADREPS_WRITEINS)
	}

	_, err := os.Stat(outputDirname)
	if os.IsNotExist(err) && os.Mkdir(outputDirname, 0755) != nil {
		log.Fatal("Could not create output directory", outputDirname)
	}

	//write to ballot measure file
	f, err := os.Create(outputDirname + "/ballot-measure.csv")
	if err != nil {
		log.Fatal(err)
	}

	f.WriteString(fmt.Sprint("Ballot Measure Yes,", ballotYes, "\n"))
	f.WriteString(fmt.Sprint("Ballot Measure No,", ballotNo, "\n"))
	f.Sync()
	f.Close()

	//write to senate file
	f, err = os.Create(outputDirname + "/senate.csv")
	if err != nil {
		log.Fatal(err)
	}

	f.WriteString("Candidate,Votes\n")
	for vote, count := range senate {
		f.WriteString("\"" + vote + "\"" + "," + fmt.Sprint(count) + "\n")
	}
	f.Sync()
	f.Close()

	//write to SFC At-large file
	f, err = os.Create(outputDirname + "/sfc-at-large.csv")
	if err != nil {
		log.Fatal(err)
	}

	f.WriteString("Candidate,Votes\n")
	for vote, count := range sfcAtLarge {
		f.WriteString("\"" + vote + "\"" + "," + fmt.Sprint(count) + "\n")
	}
	f.Sync()
	f.Close()

	//write to Undergrad Reps file
	f, err = os.Create(outputDirname + "/undergrad-reps.csv")
	if err != nil {
		log.Fatal(err)
	}

	f.WriteString("Candidate,Votes\n")
	for vote, count := range undergradReps {
		f.WriteString("\"" + vote + "\"" + "," + fmt.Sprint(count) + "\n")
	}
	f.Sync()
	f.Close()

	//write to Grad Reps file
	f, err = os.Create(outputDirname + "/grad-reps.csv")
	if err != nil {
		log.Fatal(err)
	}

	f.WriteString("Candidate,Votes\n")
	for vote, count := range gradReps {
		f.WriteString("\"" + vote + "\"" + "," + fmt.Sprint(count) + "\n")
	}
	f.Sync()
	f.Close()
}

func countPopularityVote(vote *Vote, position *map[string]int, initialPosition int, numWriteins int) {
	votes := strings.Split(vote.Raw[initialPosition], ",")
	for i := 0; i < numWriteins; i++ {
		wi := vote.Raw[initialPosition+1+i]
		if wi != "" {
			votes = append(votes, wi)
		}
	}

	var cleanVotes []string
	for _, v := range votes {
		if v != "" && v != "Write-in:" {
			v = strings.TrimSpace(v)
			v = strings.ToUpper(v)
			cleanVotes = append(cleanVotes, v)
		}
	}

	cleanVotes = removeDuplicateStr(cleanVotes)

	for _, v := range cleanVotes {
		(*position)[v]++
	}
}
