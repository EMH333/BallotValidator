package steps

import (
	"fmt"
	"hash/fnv"
	"log"
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

	//seed based on the given seed (which should already include the start and end dates)
	numberGenerator := rand.New(rand.NewSource(int64(hashAlgo.Sum64())))

	for i := 0; i < numberToPick; i++ {
		randomVal := numberGenerator.Intn(initialSizeCurrentVotes + initialSizePreviousVotes)
		//randomVal := numberGenerator.Intn(initialSizeCurrentVotes) // for the weekend draw
		//randomVal := numberGenerator.Intn(initialSizeCurrentVotes + initialSizePreviousVotes + 4) // for the 4 people who tagged us

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
		} else if randomVal < initialSizeCurrentVotes+initialSizePreviousVotes {
			winner := previousVotes[randomVal-initialSizeCurrentVotes] // pick winner
			if !util.Contains(&winners, winner) {
				winners = append(winners, winner)
				logMessages = append(logMessages, "Winner: "+winner+" chosen from past responses with random value "+fmt.Sprint(randomVal))
			} else {
				i-- //try again
			}
		} else {
			log.Fatal("Random value is out of bounds. This should never happen.")
			//winners = append(winners, "tagged "+fmt.Sprint(randomVal-initialSizeCurrentVotes-initialSizePreviousVotes))
			//logMessages = append(logMessages, "Winner: tagged "+fmt.Sprint(randomVal-initialSizeCurrentVotes-initialSizePreviousVotes)+" chosen with random value "+fmt.Sprint(randomVal))
		}
	}

	return votes, winners, util.Summary{
		StepInfo:  "Step 4: Incentives",
		Processed: initialSizeCurrentVotes,
		Valid:     initialSizeCurrentVotes,
		Invalid:   0,
		Log:       logMessages}
}
