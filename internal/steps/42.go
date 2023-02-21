package steps

import (
	"fmt"
	"log"
	"os"
	"strings"

	"ethohampton.com/BallotCleaner/internal/util"
)

//TODO set correct index values

const TALLY_SENATE_OPTIONS = 17 //index
const TALLY_SENATE_WRITEINS = 18
const TALLY_SFCATLARGE_OPTIONS = TALLY_SENATE_OPTIONS + TALLY_SENATE_WRITEINS + 1 //index
const TALLY_SFCATLARGE_WRITEINS = 5

const TALLY_PRES_OPTIONS_START = TALLY_SFCATLARGE_OPTIONS + TALLY_SFCATLARGE_WRITEINS + 2 //2 for the write-in //index
const TALLY_PRES_OPTIONS_NUMBER = 2
const TALLY_SFCCHAIR_OPTIONS_START = TALLY_PRES_OPTIONS_START + TALLY_PRES_OPTIONS_NUMBER + 2 //2 because of the writeins //index
const TALLY_SFCCHAIR_OPTIONS_NUMBER = 2

// designed to do all the counting and output a nice little summary
func StepFourtyTwo(votes []util.Vote, outputDirname string) {
	if len(votes) == 0 {
		log.Fatal("No votes to find results for")
	}

	var senate map[string]int = make(map[string]int)
	var sfcAtLarge map[string]int = make(map[string]int)

	//TODO confirm number of seats for each
	for _, vote := range votes {
		///////////////////SENATE/////////////////////////////
		countPopularityVote(&vote, &senate, TALLY_SENATE_OPTIONS, TALLY_SENATE_WRITEINS, 18)

		///////////////////SFC AT LARGE/////////////////////////////
		countPopularityVote(&vote, &sfcAtLarge, TALLY_SFCATLARGE_OPTIONS, TALLY_SFCATLARGE_WRITEINS, 5)
	}
	//log.Println("Counted Popularity Votes")

	//TODO set correct names
	//presidental ticket
	presidentResults := util.RunIRV(votes, []string{"Calvin Anderman for President & Braeden Howard for Vice President", "Alexander Kerner for President & Isabella Griffiths for Vice President"}, TALLY_PRES_OPTIONS_NUMBER, TALLY_PRES_OPTIONS_START)

	//TODO set correct names
	//SFC chair
	sfcChairResults := util.RunIRV(votes, []string{"Matteo Paola", "Noa Stoll"}, TALLY_SFCCHAIR_OPTIONS_NUMBER, TALLY_SFCCHAIR_OPTIONS_START)

	_, err := os.Stat(outputDirname)
	if os.IsNotExist(err) && os.Mkdir(outputDirname, 0755) != nil {
		log.Fatal("Could not create output directory", outputDirname)
	}

	//write to senate file
	writeMultipleVoteResults(&senate, outputDirname+"/senate.csv")

	//write to SFC At-large file
	writeMultipleVoteResults(&sfcAtLarge, outputDirname+"/sfc-at-large.csv")

	//write to president file
	writeIRVResults(presidentResults, outputDirname+"/president.txt")

	//write to SFC chair file
	writeIRVResults(sfcChairResults, outputDirname+"/sfc-chair.txt")
}

func countPopularityVote(vote *util.Vote, position *map[string]int, initialPosition int, numWriteins int, maxVotes int) {
	rawVotes := strings.Split(vote.Raw[initialPosition], ",")
	var votes []string
	// clean up the write in entries
	for _, vote := range rawVotes {
		vote = strings.TrimSpace(vote)
		if !(vote == "Write in:" || vote == "Write-in:" || vote == "Write-In") {
			votes = append(votes, vote)
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
	if len(votes) > maxVotes {
		log.Println("WARNING: More than the max number of votes for", vote.ONID)
		return
	}

	for _, v := range votes {
		if v != "" {
			(*position)[v]++
		}
	}
}

//TODO order by votes
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
