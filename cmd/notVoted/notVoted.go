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

	// Load the valid voters
	log.Println("Loading valid voters...")
	validVotersGraduate := util.LoadValidVoters("data/validVoters.csv", "G")
	validVotersUndergrad := util.LoadValidVoters("data/validVoters.csv", "UG")
	validVotersUndefined := util.LoadValidVoters("data/validVoters.csv", "Self Identified on Ballot")

	log.Printf("There are %d valid voters for graduate students\n", len(validVotersGraduate))
	log.Printf("There are %d valid voters for undergrad students\n", len(validVotersUndergrad))
	log.Printf("There are %d valid voters for undefined students\n", len(validVotersUndefined))

	// Load the votes
	log.Println("Loading votes...")
	votes := util.LoadVotesCSV("data/ballots/"+dataFile, START, END)
	log.Printf("%d votes loaded\n", len(votes))

	// reuse step two to get the ONID emails for all eligible students that have already voted
	// this won't be 100% accurate because it will also include non-corvallis/non-students,
	// but it's good enough because they won't ever make it into the final list
	_, _, alreadyVoted, _ := steps.StepTwo(votes, &[]string{})
	log.Printf("There are %d people who have already voted\n", len(alreadyVoted))

	//now loop through all valid voters and see if they have already voted
	var onidEmails []string
	for _, v := range validVotersGraduate {
		if !util.Contains(&alreadyVoted, v) {
			onidEmails = append(onidEmails, v)
		}
	}
	for _, v := range validVotersUndergrad {
		if !util.Contains(&alreadyVoted, v) {
			onidEmails = append(onidEmails, v)
		}
	}
	for _, v := range validVotersUndefined {
		if !util.Contains(&alreadyVoted, v) {
			onidEmails = append(onidEmails, v)
		}
	}

	log.Printf("There are %d students who haven't voted yet\n", len(onidEmails))

	//sort the emails because I am a nice person
	sort.Strings(onidEmails)

	//write the emails to a file
	util.StoreAlreadyVoted(onidEmails, "haveNotVoted.csv")
}
