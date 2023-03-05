package main

import (
	"log"
	"os"
	"sort"

	"ethohampton.com/BallotCleaner/internal/steps"
	"ethohampton.com/BallotCleaner/internal/util"
)

// get the ONID emails for all eligible students that haven't voted yet

const START = 0
const END = 100 //no harm in going overboard here

func main() {
	var dataFile string
	// in the form of `program <file_to_process>`
	if len(os.Args) == 2 {
		dataFile = os.Args[1]
	} else {
		log.Fatal("Need to specify a file to process")
	}

	_, err := os.Stat("output")
	if os.IsNotExist(err) && os.Mkdir("output", 0755) != nil {
		log.Fatal("Could not create output directory")
	}

	// Load the valid voters
	log.Println("Loading valid voters...")
	validVoters := util.LoadValidVoters("data/validVoters.csv")

	log.Printf("There are %d valid student voters\n", len(validVoters))

	// Load the votes
	log.Println("Loading votes...")
	votes := util.LoadVotesCSV("data/ballots/"+dataFile, START, END)
	log.Printf("%d votes loaded\n", len(votes))

	// reuse step two to get the ONID emails for all eligible students that have already voted
	// this won't be 100% accurate because it will also include non-corvallis/non-students,
	// but it's good enough because they won't ever make it into the final list
	_, _, alreadyVoted, _ := steps.StepTwo(votes, &[]string{})
	log.Printf("There are %d people who have already voted\n", len(alreadyVoted))

	// load people who don't want to be counted and consider them to have already voted
	doNotCount := util.LoadStringArrayFile("data/doNotRemind.csv")
	alreadyVoted = append(alreadyVoted, doNotCount...)
	log.Printf("There are %d people who don't want to be reminded\n", len(doNotCount))

	//now loop through all valid voters and see if they have already voted
	var onidEmails []string
	for _, v := range validVoters {
		if !util.Contains(&alreadyVoted, v) {
			onidEmails = append(onidEmails, v)
		}
	}

	log.Printf("There are %d students who haven't voted yet\n", len(onidEmails))

	//sort the emails because I am a nice person
	sort.Strings(onidEmails)

	//write the emails to a file
	util.StoreNotYetVoted(onidEmails, "haveNotVoted.csv")
}
