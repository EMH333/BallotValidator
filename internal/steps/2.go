package steps

import (
	"log"
	"maps"
	"slices"

	"ethohampton.com/BallotValidator/internal/util"
)

// returns valid, invaild, onids that voted and summary
func StepTwo(votes []util.Vote, alreadyVoted *[]string) ([]util.Vote, []util.Vote, []string, util.Summary) {
	var initialSize = len(votes)

	var logMessages []string

	validVotes := make([]util.Vote, 0, initialSize)
	var invalidVotes []util.Vote
	votedToday := make(map[string]bool, len(votes)-len(*alreadyVoted))

	for _, v := range votes {
		if votedToday[v.ONID] || slices.Contains(*alreadyVoted, v.ONID) {
			invalidVotes = append(invalidVotes, v)
			logMessages = append(logMessages, "Invalid vote from "+v.ONID+" with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05"))
		} else {
			validVotes = append(validVotes, v)
			votedToday[v.ONID] = true
		}
	}

	if len(validVotes)+len(invalidVotes) != initialSize {
		log.Fatal("Step 2 vote counts don't match")
	}

	return validVotes, invalidVotes, slices.Collect(maps.Keys(votedToday)), util.Summary{
		StepInfo:  "Step 2: Dedupe",
		Processed: len(validVotes) + len(invalidVotes),
		Valid:     len(validVotes),
		Invalid:   len(invalidVotes),
		Log:       logMessages}
}
