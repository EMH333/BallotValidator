package steps

import (
	"log"

	"ethohampton.com/BallotValidator/internal/util"
)

func StepOne(votes []util.Vote, validVoters *[]string) ([]util.Vote, []util.Vote, util.Summary) {
	var initialSize int = len(votes)

	var messageLog []string

	var validVotes []util.Vote
	var invalidVotes []util.Vote

	for _, v := range votes {
		if util.Contains(validVoters, v.ONID) {
			validVotes = append(validVotes, v)
		} else {
			invalidVotes = append(invalidVotes, v)
			messageLog = append(messageLog, "Invalid vote from "+v.ONID+" with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05"))
		}
	}

	if len(validVotes)+len(invalidVotes) != initialSize {
		log.Fatal("Step 1 vote counts don't match")
	}

	return validVotes, invalidVotes, util.Summary{
		StepInfo:  "Step 1: Valid voter",
		Processed: len(validVotes) + len(invalidVotes),
		Valid:     len(validVotes),
		Invalid:   len(invalidVotes),
		Log:       messageLog}
}
