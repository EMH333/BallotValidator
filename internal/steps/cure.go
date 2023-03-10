package steps

import (
	"log"
	"strings"
	"time"

	"ethohampton.com/BallotCleaner/internal/util"
)

const START_TIME_INDEX = 0
const MAX_POTENTIAL_CURED_BALLOTS = 45 //this is based on previous analysis of start times and senate vote counts. Should never exceed this

//returns valid, invaild, cured ballot summary and summary
func StepCure(votes []util.Vote) ([]util.Vote, []util.Vote, []string, util.Summary) {
	var initialSize int = len(votes)

	var logMessages []string
	var validVotes []util.Vote
	var invalidVotes []util.Vote

	var BEFORE_TIME_CURE = time.Date(2023, time.February, 20, 17, 32, 0, 0, time.Local)

	var curedBallotSummary []string
	var ballotsBeforeTime = 0

	//find all voters who started voting in the first 31 minutes the poll was open and voted for 6 senators
	//if they have voted more than once, then remove their first vote so that they can vote again and have all 18 senate options avaliable

	for _, v := range votes {
		t, err := time.ParseInLocation(util.BALLOT_TIME_FORMAT, v.Raw[START_TIME_INDEX], time.Local)
		if err != nil {
			log.Fatal(err)
		}

		//only check for people who started voting before the cure time
		if t.Before(BEFORE_TIME_CURE) {
			ballotsBeforeTime++
			//now we need to check for 6 senate votes
			var senateEntry = v.Raw[TALLY_SENATE_OPTIONS]
			var senateVotes = strings.Split(senateEntry, ",")
			if len(senateVotes) == 6 {
				//now we need to check if they have voted more than once
				if votedMoreThanOnce(votes, v.ONID) {
					//if they have voted more than once, then remove their first vote so that they can vote again and have all 18 senate options avaliable
					invalidVotes = append(invalidVotes, v)
					logMessages = append(logMessages, "Initial vote from "+v.ONID+" with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05")+" removed because they voted again")
					curedBallotSummary = append(curedBallotSummary, v.ONID+",yes") //did remove from list because they voted more than once, so count the next one
				} else {
					validVotes = append(validVotes, v)
					curedBallotSummary = append(curedBallotSummary, v.ONID+",no") //did not remove from list because they voted only once
				}
			} else {
				validVotes = append(validVotes, v)
			}
		} else {
			validVotes = append(validVotes, v)
		}
	}

	log.Println("Cure step ballots before time:", ballotsBeforeTime)

	if len(curedBallotSummary) > MAX_POTENTIAL_CURED_BALLOTS || len(invalidVotes) > MAX_POTENTIAL_CURED_BALLOTS {
		log.Fatal("Cure step invalid count ", len(curedBallotSummary), len(invalidVotes))
	}

	if len(validVotes)+len(invalidVotes) != initialSize {
		log.Fatal("Cure step vote counts don't match")
	}

	return validVotes, invalidVotes, curedBallotSummary, util.Summary{
		StepInfo:  "Step Cure: Remove initial votes from people curing their ballots",
		Processed: len(validVotes) + len(invalidVotes),
		Valid:     len(validVotes),
		Invalid:   len(invalidVotes),
		Log:       logMessages}
}

func votedMoreThanOnce(votes []util.Vote, onid string) bool {
	var count int = 0
	for _, v := range votes {
		if v.ONID == onid {
			count++
		}
	}
	return count > 1
}
