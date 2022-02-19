package steps

import (
	"fmt"
	"hash/fnv"
	"math/rand"

	"ethohampton.com/BallotCleaner/internal/util"
)

func StepFour(votes []util.Vote, seed string, numberToPick int) ([]util.Vote, []string, util.Summary) {
	var initialSize int = len(votes)
	var winners []string
	var logMessages []string

	var hashAlgo = fnv.New64a()
	hashAlgo.Write([]byte(seed))

	rand.Seed(int64(hashAlgo.Sum64())) //seed based on the given seed (which should already include the start and end dates)

	for i := 0; i < numberToPick; i++ {
		randomVal := rand.Intn(initialSize)
		winner := votes[randomVal] // pick winner
		if !util.Contains(&winners, winner.ONID) {
			winners = append(winners, winner.ONID)
			logMessages = append(logMessages, "Winner: "+winner.ONID+" with response ID "+winner.ID+" chosen with random value "+fmt.Sprint(randomVal))
		} else {
			i-- //try again
		}
	}

	return votes, winners, util.Summary{
		StepInfo:  "Step 4: Incentives",
		Processed: initialSize,
		Valid:     initialSize,
		Invalid:   0,
		Log:       logMessages}
}
