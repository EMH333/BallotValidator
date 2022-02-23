package steps

import (
	"fmt"
	"hash/fnv"
	"math/rand"

	"ethohampton.com/BallotCleaner/internal/util"
)

func StepFour(previousVotes []string, votes []util.Vote, seed string, numberToPick int) ([]util.Vote, []string, util.Summary) {
	var initialSizeCurrentVotes int = len(votes)
	var initialSizePreviousVotes int = len(previousVotes)
	var winners []string
	var logMessages []string

	var hashAlgo = fnv.New64a()
	hashAlgo.Write([]byte(seed))

	rand.Seed(int64(hashAlgo.Sum64())) //seed based on the given seed (which should already include the start and end dates)

	for i := 0; i < numberToPick; i++ {
		randomVal := rand.Intn(initialSizeCurrentVotes + initialSizePreviousVotes)
		//if the random val is in the current votes, then grab the vote and add it to the winners
		//otherwise it is in the previous votes, so use that
		if randomVal < initialSizeCurrentVotes {
			winner := votes[randomVal] // pick winner
			if !util.Contains(&winners, winner.ONID) {
				winners = append(winners, winner.ONID)
				logMessages = append(logMessages, "Winner: "+winner.ONID+" with response ID "+winner.ID+" chosen with random value "+fmt.Sprint(randomVal))
			} else {
				i-- //try again
			}
		} else {
			winner := previousVotes[randomVal-initialSizeCurrentVotes] // pick winner
			if !util.Contains(&winners, winner) {
				winners = append(winners, winner)
				logMessages = append(logMessages, "Winner: "+winner+" chosen from past responses with random value "+fmt.Sprint(randomVal))
			} else {
				i-- //try again
			}
		}
	}

	return votes, winners, util.Summary{
		StepInfo:  "Step 4: Incentives",
		Processed: initialSizeCurrentVotes,
		Valid:     initialSizeCurrentVotes,
		Invalid:   0,
		Log:       logMessages}
}
