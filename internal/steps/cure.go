package steps

import (
	"log"
	"strings"
	"time"

	"ethohampton.com/BallotCleaner/internal/util"
)

const START_TIME_INDEX = 0
const MAX_POTENTIAL_CURED_BALLOTS = 45 //this is based on previous analysis of start times and senate vote counts. Should never exceed this
const SENATE_VOTES_ONID_INDEX = 36

//returns valid, invaild, cured ballot summary and summary
func StepCure(votes []util.Vote, senateVotesFile string) ([]util.Vote, []util.Vote, []string, util.Summary) {
	var initialSize int = len(votes)

	var logMessages []string
	var validVotes []util.Vote
	var invalidVotes []util.Vote

	var BEFORE_TIME_CURE = time.Date(2023, time.February, 20, 17, 32, 0, 0, time.Local)

	var curedBallotSummary []string
	var ballotsBeforeTime = 0

	// load senate votes
	newSenateVotes := util.LoadVotesCSV(senateVotesFile, 0, 100, SENATE_VOTES_ONID_INDEX)

	//find all voters who started voting in the first 31 minutes the poll was open and voted for 6 senators
	//see if they have voted using the new senator ballot, and if so, replace their old senate ballot with the new one
	for _, v := range votes {
		t, err := time.ParseInLocation(util.BALLOT_TIME_FORMAT, v.Raw[START_TIME_INDEX], time.Local)
		if err != nil {
			log.Fatal(err)
		}

		// make sure we don't overwrite sfc votes
		sfcVote := v.Raw[TALLY_SFCATLARGE_OPTIONS]

		//only check for people who started voting before the cure time
		if t.Before(BEFORE_TIME_CURE) {
			//now we need to check for 6 senate votes
			var originalSenateEntry = v.Raw[TALLY_SENATE_OPTIONS]
			var originalSenateVotes = strings.Split(originalSenateEntry, ",")
			if len(originalSenateVotes) == 6 {
				ballotsBeforeTime++
				//now we replace their old senate ballot with the new one if it exists
				if indexOfVote(newSenateVotes, v.ONID) != -1 {
					//if they voted with new senate ballot, then replace their old senate choices with their new ones
					invalidVotes = append(invalidVotes, v)
					logMessages = append(logMessages, "Initial vote from "+v.ONID+" with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05")+" replaced because they voted with the new senate ballot")
					curedBallotSummary = append(curedBallotSummary, v.ONID+",yes") //did change because they voted in the senate ballot

					//now we need to replace their old senate ballot with the new one
					for i := TALLY_SENATE_OPTIONS; i <= TALLY_SENATE_WRITEINS; i++ {
						v.Raw[i] = newSenateVotes[indexOfVote(newSenateVotes, v.ONID)].Raw[i]
					}
					validVotes = append(validVotes, v)

					//make sure we didn't change their sfc vote
					if v.Raw[TALLY_SFCATLARGE_OPTIONS] != sfcVote {
						log.Fatal("SFC vote changed")
					}
				} else {
					validVotes = append(validVotes, v)
					curedBallotSummary = append(curedBallotSummary, v.ONID+",no") //did not change because they did't vote in the senate ballot
				}
			} else {
				validVotes = append(validVotes, v)
			}
		} else {
			validVotes = append(validVotes, v)
		}
	}

	log.Println("Cure step ballots before time and 6 votes:", ballotsBeforeTime)
	log.Println("New senate votes:", len(newSenateVotes))

	if len(curedBallotSummary) > MAX_POTENTIAL_CURED_BALLOTS || len(invalidVotes) > MAX_POTENTIAL_CURED_BALLOTS {
		log.Fatal("Cure step invalid count ", len(curedBallotSummary), len(invalidVotes))
	}

	if len(invalidVotes) != len(newSenateVotes) {
		log.Fatal("Cure step count doesn't match new senate votes")
	}

	if len(validVotes) != initialSize {
		log.Fatal("Cure step vote counts don't match")
	}

	return validVotes, invalidVotes, curedBallotSummary, util.Summary{
		StepInfo:  "Step Cure: Remove initial votes from people curing their ballots",
		Processed: len(validVotes) + len(invalidVotes),
		Valid:     len(validVotes),
		Invalid:   len(invalidVotes),
		Log:       logMessages}
}

func indexOfVote(votes []util.Vote, onid string) int {
	for i, v := range votes {
		if v.ONID == onid {
			return i
		}
	}
	return -1
}
