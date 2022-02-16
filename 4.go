package main

import (
	"fmt"
	"hash/fnv"
	"math/rand"
)

func stepFour(votes []Vote, seed string, numberToPick int) ([]Vote, []string, Summary) {
	var initialSize int = len(votes)
	var winners []string
	var logMessages []string

	var hashAlgo = fnv.New64a()
	hashAlgo.Write([]byte(seed))

	rand.Seed(int64(hashAlgo.Sum64())) //seed based on the given seed (which should already include the start and end dates)

	for i := 0; i < numberToPick; i++ {
		randomVal := rand.Intn(initialSize)
		winner := votes[randomVal] // pick winner
		if !contains(&winners, winner.ONID) {
			winners = append(winners, winner.ONID)
			logMessages = append(logMessages, "Winner: "+winner.ONID+" with response ID "+winner.ID+" chosen with random value "+fmt.Sprint(randomVal))
		} else {
			i-- //try again
		}
	}

	return votes, winners, Summary{processed: initialSize, valid: initialSize, invalid: 0, log: logMessages}
}
