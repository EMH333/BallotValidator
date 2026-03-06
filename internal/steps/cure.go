package steps

import (
	"log"
	"slices"
	"strings"

	"ethohampton.com/BallotValidator/internal/util"
)

const NEW_SFCAL_INDEX = 17
const MAX_POTENTIAL_CURED_BALLOTS = 50 //this is based on previous analysis of start times and sfc at-large vote counts. Should never exceed this

// returns valid, invaild, cured ballot summary and summary
func StepCure(countingConfig *util.CountingConfig, votes []util.Vote, sfcatlargefile, allowedVotersFile string) ([]util.Vote, []util.Vote, []string, util.Summary) {
	var initialSize = len(votes)

	var logMessages []string
	var validVotes []util.Vote
	var invalidVotes []util.Vote

	var curedBallotSummary []string
	var ballotsModified = 0

	// we can't use normal config ONID index, because it is at the end
	sfcalConfig := *countingConfig
	sfcalConfig.ImportONID = 18

	// load votes
	sfcalVotes := util.LoadVotesCSV(&sfcalConfig, sfcatlargefile, 0, 100)
	allowedToOverwrite := util.LoadStringArrayFile(allowedVotersFile)

	for _, v := range votes {
		// only change for those approved to fix vote
		if !slices.Contains(allowedToOverwrite, v.ONID) {
			validVotes = append(validVotes, v)
			continue
		}

		// continue if they did not fix vote
		fixedVoteIndex := indexOfVote(sfcalVotes, v.ONID)
		if fixedVoteIndex == -1 {
			validVotes = append(validVotes, v)
			continue
		}

		originalSFCALEntry := v.Raw[countingConfig.TallySFCAtLargeOptionsIndex]
		originalSFCALVotes := strings.Split(originalSFCALEntry, ",")
		if len(originalSFCALVotes) == 1 {
			ballotsModified++
			//if they voted with new sfc at large ballot, then replace their old sfc at large choices with their new ones
			invalidVotes = append(invalidVotes, v)
			logMessages = append(logMessages, "Initial vote from "+v.ONID+" with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05")+" replaced because they voted with the new SFCAL ballot")
			curedBallotSummary = append(curedBallotSummary, v.ONID+",yes") //did change because they voted in the sfc at large ballot

			// log.Println(v.ONID)

			//now we need to replace their old sfc at large ballot with the new one
			v.Raw[countingConfig.TallySFCAtLargeOptionsIndex] = sfcalVotes[fixedVoteIndex].Raw[NEW_SFCAL_INDEX]
			validVotes = append(validVotes, v)
		} else {
			logMessages = append(logMessages, "WARNING: Initial vote from "+v.ONID+" with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05")+" voted in new SFCAL ballot, but selected multiple candidates")
			validVotes = append(validVotes, v)
		}
	}

	// log.Println()
	// for _, v := range sfcalVotes {
	// 	log.Println(v.ONID)
	// }

	log.Println("Cure step ballots modified:", ballotsModified)
	log.Println("New SFCAL votes:", len(sfcalVotes))

	if len(curedBallotSummary) > MAX_POTENTIAL_CURED_BALLOTS || len(invalidVotes) > MAX_POTENTIAL_CURED_BALLOTS {
		log.Fatal("Cure step invalid count ", len(curedBallotSummary), len(invalidVotes))
	}

	if len(invalidVotes) != len(sfcalVotes) {
		log.Fatal("Cure step count doesn't match new SFCAL votes")
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
