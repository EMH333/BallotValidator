package steps

import (
	"fmt"
	"log"
	"os"
	"strings"

	"ethohampton.com/BallotCleaner/internal/util"
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

const TALLY_SPEAKER_OPTIONS_START = TALLY_GRADREPS_OPTIONS + TALLY_GRADREPS_WRITEINS + 1
const TALLY_SPEAKER_OPTIONS_NUMBER = 6                                                          //doesn't include writein
const TALLY_PRES_OPTIONS_START = TALLY_SPEAKER_OPTIONS_START + TALLY_SPEAKER_OPTIONS_NUMBER + 2 //2 for the write-in
const TALLY_PRES_OPTIONS_NUMBER = 3
const TALLY_SFCCHAIR_OPTIONS_START = TALLY_PRES_OPTIONS_START + TALLY_PRES_OPTIONS_NUMBER + 2 //2 because of the writeins
const TALLY_SFCCHAIR_OPTIONS_NUMBER = 1

// designed to do all the counting and output a nice little summary
func StepFourtyTwo(votes []util.Vote, outputDirname string) {
	if len(votes) == 0 {
		log.Fatal("No votes to find results for")
	}

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

	//speaker of the house
	speakerResults := util.RunIRV(votes, []string{"Abril Uribe", "Jordan Porter", "Madelyn Neuschwander", "Jared Pratt", "Rooney Ferguson", "Jarrett Alto"}, TALLY_SPEAKER_OPTIONS_NUMBER, TALLY_SPEAKER_OPTIONS_START)

	//presidental ticket
	presidentResults := util.RunIRV(votes, []string{"Calvin Anderman for President & Braeden Howard for Vice President", "Alexander Kerner for President & Isabella Griffiths for Vice President", "Matteo Paola for President & Sierra Young for Vice President"}, TALLY_PRES_OPTIONS_NUMBER, TALLY_PRES_OPTIONS_START)

	//SFC chair
	sfcChairResults := util.RunIRV(votes, []string{"Joe Page"}, TALLY_SFCCHAIR_OPTIONS_NUMBER, TALLY_SFCCHAIR_OPTIONS_START)

	_, err := os.Stat(outputDirname)
	if os.IsNotExist(err) && os.Mkdir(outputDirname, 0755) != nil {
		log.Fatal("Could not create output directory", outputDirname)
	}

	//write to ballot measure file
	f, err := os.Create(outputDirname + "/ballot-measure.csv")
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString(fmt.Sprint("Ballot Measure Yes,", ballotYes, "\n"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(fmt.Sprint("Ballot Measure No,", ballotNo, "\n"))
	if err != nil {
		log.Fatal(err)
	}
	err = f.Sync()
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	//write to senate file
	writeMultipleVoteResults(&senate, outputDirname+"/senate.csv")

	//write to SFC At-large file
	writeMultipleVoteResults(&sfcAtLarge, outputDirname+"/sfc-at-large.csv")

	//write to Undergrad Reps file
	writeMultipleVoteResults(&undergradReps, outputDirname+"/undergrad-reps.csv")

	//write to Grad Reps file
	writeMultipleVoteResults(&gradReps, outputDirname+"/grad-reps.csv")

	//write to speaker of the house file
	writeIRVResults(speakerResults, outputDirname+"/speaker-of-the-house.txt")

	//write to president file
	writeIRVResults(presidentResults, outputDirname+"/president.txt")

	//write to SFC chair file
	writeIRVResults(sfcChairResults, outputDirname+"/sfc-chair.txt")
}

func countPopularityVote(vote *util.Vote, position *map[string]int, initialPosition int, numWriteins int) {
	votes := strings.Split(vote.Raw[initialPosition], ",")
	// clean up the write in entries
	for i, vote := range votes {
		if vote == "Write in:" || vote == "Write-in:" {
			votes[i] = ""
		}
	}

	for i := 0; i < numWriteins; i++ {
		wi := vote.Raw[initialPosition+1+i]
		if wi != "" {
			votes = append(votes, util.CleanVote(wi))
		}
	}

	votes = util.RemoveDuplicateStr(votes)

	// can't pick more than the max for these positions
	if len(votes) > numWriteins {
		log.Println("WARNING: More than the max number of votes for", vote.ONID)
		return
	}

	for _, v := range votes {
		if v != "" {
			(*position)[v]++
		}
	}
}

func writeMultipleVoteResults(results *map[string]int, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString("Candidate,Votes\n")
	if err != nil {
		log.Fatal(err)
	}

	for vote, count := range *results {
		_, err = f.WriteString("\"" + vote + "\"" + "," + fmt.Sprint(count) + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	err = f.Sync()
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}

func writeIRVResults(results []string, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range results {
		_, err = f.WriteString(v + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	err = f.Sync()
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}
