package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"ethohampton.com/BallotCleaner/internal/steps"
	"ethohampton.com/BallotCleaner/internal/util"
)

var validVotersGraduate []string
var validVotersUndergrad []string
var validVotersUndefined []string

var alreadyVotedPrevious []string

const numToPick int = 15 // how many winners to pick //TODO tie into command line

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
	var seed = util.LoadSeed() + "-" + dayToDayFormat //include days in picking so it is unique

	_, err := os.Stat("output")
	if os.IsNotExist(err) && os.Mkdir("output", 0755) != nil {
		log.Fatal("Could not create output directory")
	}

	// Load the valid voters
	log.Println("Loading valid voters...")
	validVotersGraduate = util.LoadValidVoters("data/validVoters.csv", "G")
	validVotersUndergrad = util.LoadValidVoters("data/validVoters.csv", "UG")
	validVotersUndefined = util.LoadValidVoters("data/validVoters.csv", "Self Identified on Ballot")

	log.Printf("There are %d valid voters for graduate students\n", len(validVotersGraduate))
	log.Printf("There are %d valid voters for undergrad students\n", len(validVotersUndergrad))
	log.Printf("There are %d valid voters for undefined students\n", len(validVotersUndefined))

	// Load the already voted
	log.Printf("Loading already voted up to day %d...\n", startDay)
	alreadyVotedPrevious = util.LoadAlreadyVoted("alreadyVoted", int64(startDay))

	log.Printf("%d students have already voted\n", len(alreadyVotedPrevious))

	// Load the votes
	log.Println("Loading votes...")
	votes := util.LoadVotesCSV("data/ballots/"+dataFile, startDay, endDay)
	log.Printf("%d votes loaded for day %d through %d\n", len(votes), startDay, endDay)

	// step one: valid voter
	log.Println()
	log.Println("Step 1: Valid voter")
	validPostOne, invalidPostOne, oneSummary := steps.StepOne(votes, &validVotersGraduate, &validVotersUndergrad, &validVotersUndefined)
	util.StoreVotes(validPostOne, "1-valid-"+dayToDayFormat+".csv")
	util.StoreVotes(invalidPostOne, "1-invalid-"+dayToDayFormat+".csv")
	util.StoreSummary(oneSummary, "1-summary-"+dayToDayFormat+".txt")
	log.Println("Step 1: Invalid votes:", oneSummary.Invalid)
	log.Println("Step 1: Valid votes:", oneSummary.Valid)

	// step two: valid voter
	log.Println()
	log.Println("Step 2: Dedupe")
	validPostTwo, invalidPostTwo, alreadyVotedToday, twoSummary := steps.StepTwo(validPostOne, &alreadyVotedPrevious)
	util.StoreVotes(validPostTwo, "2-valid-"+dayToDayFormat+".csv")
	util.StoreVotes(invalidPostTwo, "2-invalid-"+dayToDayFormat+".csv")
	util.StoreSummary(twoSummary, "2-summary-"+dayToDayFormat+".txt")
	util.StoreAlreadyVoted(alreadyVotedToday, "alreadyVoted-"+dayToDayFormat+".csv")
	log.Println("Step 2: Invalid votes:", twoSummary.Invalid)
	log.Println("Step 2: Valid votes:", twoSummary.Valid)

	// step three: grad/undergrad
	log.Println()
	log.Println("Step 3: Grad/undergrad")
	validPostThree, invalidPostThree, threeSummary := steps.StepThree(validPostTwo, &validVotersGraduate, &validVotersUndergrad)
	util.StoreVotes(validPostThree, "3-valid-"+dayToDayFormat+".csv")
	util.StoreVotes(invalidPostThree, "3-modified-"+dayToDayFormat+".csv")
	util.StoreSummary(threeSummary, "3-summary-"+dayToDayFormat+".txt")
	log.Println("Step 3: Modified votes:", threeSummary.Invalid)
	log.Println("Step 3: Valid votes:", threeSummary.Valid)

	// step four: Incentives
	log.Println()
	log.Println("Step 4: Incentives")
	postFour, winners, fourSummary := steps.StepFour(alreadyVotedPrevious, validPostThree, seed, numToPick)
	util.StoreVotes(postFour, "4-valid-"+dayToDayFormat+".csv")
	util.StoreSummary(fourSummary, "4-summary-"+dayToDayFormat+".txt")
	util.StoreAlreadyVoted(winners, "incentive-winners-"+dayToDayFormat+".csv")
	log.Println("Step 4: Valid votes:", threeSummary.Valid)
	log.Println("Step 4: Selected winners:", len(winners))

	//only figure out the winners if we are across multiple days
	if startDay != endDay {
		//experimental
		steps.StepFourtyTwo(postFour, "output/results")
	}
}
