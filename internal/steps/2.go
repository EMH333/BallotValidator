package steps

import (
	"log"

	"ethohampton.com/BallotValidator/internal/util"
)

//returns valid, invaild, onids that voted and summary
func StepTwo(votes []util.Vote, alreadyVoted *[]string) ([]util.Vote, []util.Vote, []string, util.Summary) {
	var initialSize int = len(votes)

	var logMessages []string

	var validVotes []util.Vote
	var invalidVotes []util.Vote
	var votedToday []string

	for _, v := range votes {
		if util.Contains(alreadyVoted, v.ONID) || util.Contains(&votedToday, v.ONID) {
			invalidVotes = append(invalidVotes, v)
			logMessages = append(logMessages, "Invalid vote from "+v.ONID+" with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05"))
		} else {
			validVotes = append(validVotes, v)
			votedToday = append(votedToday, v.ONID)
		}
	}

	if len(validVotes)+len(invalidVotes) != initialSize {
		log.Fatal("Step 2 vote counts don't match")
	}

	return validVotes, invalidVotes, votedToday, util.Summary{
		StepInfo:  "Step 2: Dedupe",
		Processed: len(validVotes) + len(invalidVotes),
		Valid:     len(validVotes),
		Invalid:   len(invalidVotes),
		Log:       logMessages}
}
