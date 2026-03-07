package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"ethohampton.com/BallotValidator/internal/steps"
	"ethohampton.com/BallotValidator/internal/util"
)

var validVotersGraduate []string
var validVotersUndergrad []string
var validVotersUndefined []string

var alreadyVotedPrevious []string

const numToPick int = 28 // how many winners to pick

func main() {
	var startDay = 0      // what day are we starting on to process votes
	var endDay = startDay // what day are we ending on to process votes
	var dataFile = "ballotData.csv"
	var countingConfigFile = "counting_config.json"

	// in the form of `program <day> <file_to_process>`
	if len(os.Args) == 3 {
		day, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			log.Fatal("Couldn't parse argument")
		}
		startDay = int(day)
		endDay = int(day)
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

		startDay = int(day1)
		endDay = int(day2)
		dataFile = os.Args[3]
	}

	log.Printf("Selected start day: %d, Selected end day: %d\n", startDay, endDay)
	var dayToDayFormat = fmt.Sprint(startDay) + "-" + fmt.Sprint(endDay)

	//"random" seed so winners are deterministic
	var seed = util.LoadSeed() + "-" + dayToDayFormat //include days in picking so it is unique
	log.Println("Seed:", seed)

	log.Printf("Loading counting config from: %s", countingConfigFile)
	countingConfig := util.LoadCountingConfig(countingConfigFile)

	_, err := os.Stat("output")
	if os.IsNotExist(err) && os.Mkdir("output", 0o755) != nil {
		log.Fatal("Could not create output directory")
	}

	// Load the valid voters
	log.Println("Loading valid voters...")
	//TODO eventually handle this better
	validVotersFile := util.LoadFileToReader(countingConfig.ValidVotersFile)
	validVotersFileBytes, err := io.ReadAll(validVotersFile)
	if err != nil {
		log.Fatal(err)
	}
	validVotersFileContents := string(validVotersFileBytes)
	validVotersGraduate = util.LoadValidVoters(&countingConfig, bytes.NewBufferString(validVotersFileContents), "G")
	validVotersUndergrad = util.LoadValidVoters(&countingConfig, bytes.NewBufferString(validVotersFileContents), "UG")
	validVotersUndefined = util.LoadValidVoters(&countingConfig, bytes.NewBufferString(validVotersFileContents), "INTO non-UG/G")
	validVotersTotal := len(util.LoadValidVoters(&countingConfig, bytes.NewBufferString(validVotersFileContents), ""))
	validVotersTotalActual := len(validVotersGraduate) + len(validVotersUndergrad) + len(validVotersUndefined)
	log.Printf("There are %d valid voters for graduate students\n", len(validVotersGraduate))
	log.Printf("There are %d valid voters for undergrad students\n", len(validVotersUndergrad))
	log.Printf("There are %d valid voters for undefined students\n", len(validVotersUndefined))
	if validVotersTotal != validVotersTotalActual {
		log.Printf("There are %d valid voters total, which should add up to %d\n", validVotersTotalActual, validVotersTotal)
	}

	// Load the already voted
	log.Printf("Loading already voted up to day %d...\n", startDay)
	alreadyVotedPrevious = util.LoadAlreadyVoted(&countingConfig, int64(startDay))
	log.Printf("%d students have already voted\n", len(alreadyVotedPrevious))
	// print warning to make sure results are accurate
	if startDay != endDay && len(alreadyVotedPrevious) > 0 {
		log.Println("Warning: already voted data is being used across multiple days, this should not be done for the final results")
	}

	// Load the votes
	log.Println("Loading votes...")
	votes := util.LoadVotesCSV(&countingConfig, "data/ballots/"+dataFile, startDay, endDay)
	log.Printf("%d votes loaded for day %d through %d\n", len(votes), startDay, endDay)
	util.StoreVotes(votes, "original-"+dayToDayFormat+".csv")

	// step one: valid voter
	log.Println()
	log.Println("Step 1: Valid voter")
	validPostOne, invalidPostOne, oneSummary := steps.StepOne(votes, &validVotersGraduate, &validVotersUndergrad, &validVotersUndefined)
	util.StoreVotes(validPostOne, "1-valid-"+dayToDayFormat+".csv")
	util.StoreVotes(invalidPostOne, "1-invalid-"+dayToDayFormat+".csv")
	util.StoreSummary(oneSummary, "1-summary-"+dayToDayFormat+".txt")
	log.Println("Step 1: Invalid votes:", oneSummary.Invalid)
	log.Println("Step 1: Valid votes:", oneSummary.Valid)
	votes = []util.Vote{} // nolint:ineffassign,staticcheck // should not be used beyond this point

	// step two: valid voter
	log.Println()
	log.Println("Step 2: Dedupe")
	validPostTwo, invalidPostTwo, alreadyVotedToday, twoSummary := steps.StepTwo(validPostOne, &alreadyVotedPrevious)
	util.StoreVotes(validPostTwo, "2-valid-"+dayToDayFormat+".csv")
	util.StoreVotes(invalidPostTwo, "2-invalid-"+dayToDayFormat+".csv")
	util.StoreSummary(twoSummary, "2-summary-"+dayToDayFormat+".txt")
	util.StoreStringArrayFile(alreadyVotedToday, "alreadyVoted-"+dayToDayFormat+".csv", true)
	log.Println("Step 2: Invalid votes:", twoSummary.Invalid)
	log.Println("Step 2: Valid votes:", twoSummary.Valid)
	validPostOne = []util.Vote{} // nolint:ineffassign,staticcheck // should not be used beyond this point

	// step three: grad/undergrad
	log.Println()
	log.Println("Step 3: Grad/undergrad")
	validPostThree, invalidPostThree, threeSummary := steps.StepThree(&countingConfig, validPostTwo, &validVotersGraduate, &validVotersUndergrad)
	util.StoreVotes(validPostThree, "3-valid-"+dayToDayFormat+".csv")
	util.StoreVotes(invalidPostThree, "3-modified-"+dayToDayFormat+".csv")
	util.StoreSummary(threeSummary, "3-summary-"+dayToDayFormat+".txt")
	log.Println("Step 3: Modified votes:", threeSummary.Invalid)
	log.Println("Step 3: Valid votes:", threeSummary.Valid)
	validPostTwo = []util.Vote{} // nolint:ineffassign,staticcheck // should not be used beyond this point

	// step four: Incentives
	log.Println()
	log.Println("Step 4: Incentives")
	postFour, winners, fourSummary := steps.StepFour(alreadyVotedPrevious, validPostThree, seed, numToPick)
	util.StoreVotes(postFour, "4-valid-"+dayToDayFormat+".csv")
	util.StoreSummary(fourSummary, "4-summary-"+dayToDayFormat+".txt")
	util.StoreStringArrayFile(winners, "incentive-winners-"+dayToDayFormat+".csv", false)
	log.Println("Step 4: Valid votes:", twoSummary.Valid)
	log.Println("Step 4: Selected winners:", len(winners))
	validPostThree = []util.Vote{} // nolint:ineffassign,staticcheck // should not be used beyond this point

	// cure
	// For students only given the option to select one SFC At-Large candidate
	log.Println()
	log.Println("Step 5: Cure SFC At-Large")
	postCure, cureInvalidVotes, curedBallotSummary, cureSummary := steps.StepCure(&countingConfig, postFour, "data/ballots/SFCAL-Mar-6-Final.csv", "data/validVotersSfcal.csv")
	util.StoreVotes(postCure, "c-valid-"+dayToDayFormat+".csv")
	util.StoreVotes(cureInvalidVotes, "c-modified-"+dayToDayFormat+".csv")
	util.StoreStringArrayFile(curedBallotSummary, "c-summary-"+dayToDayFormat+".csv", false)
	util.StoreSummary(cureSummary, "c-summary-"+dayToDayFormat+".txt")
	log.Println("Step 5: Modified votes:", cureSummary.Invalid)
	log.Println("Step 5: Valid votes:", cureSummary.Valid)
	postFour = []util.Vote{} // nolint:ineffassign,staticcheck // should not be used beyond this point

	//only figure out the winners if we are across multiple days
	if startDay != endDay {
		steps.StepFourtyTwo(&countingConfig, postCure, "output/results")
	} else {
		log.Println("Not running step 42, only one day")
		log.Println("Adding already voted to the already voted data directory")
		util.StoreStringArrayFile(alreadyVotedToday, "../data/alreadyVoted/alreadyVoted-"+dayToDayFormat+".csv", true)
	}

	log.Println("Done")
}
