package steps

import (
	"log"
	"maps"

	"ethohampton.com/BallotValidator/internal/util"
)

func StepOne(votes []util.Vote, validVotersGraduate, validVotersUndergraduate, validVotersUndefined *[]string) ([]util.Vote, []util.Vote, util.Summary) {
	var initialSize = len(votes)

	var messageLog []string

	validVotes := make([]util.Vote, 0, initialSize)
	var invalidVotes []util.Vote

	// put into a map for performance (plus its just more natural)
	allValidVoters := make(map[string]struct{}, len(*validVotersUndergraduate)+len(*validVotersGraduate)+len(*validVotersUndefined))
	maps.Insert(allValidVoters, util.StringSliceToMap(*validVotersUndergraduate))
	maps.Insert(allValidVoters, util.StringSliceToMap(*validVotersGraduate))
	maps.Insert(allValidVoters, util.StringSliceToMap(*validVotersUndefined))

	for _, v := range votes {
		if _, ok := allValidVoters[v.ONID]; ok {
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
