package main

import "log"

//returns valid, invaild, onids that voted and summary
func stepTwo(votes []Vote, alreadyVoted *[]string) ([]Vote, []Vote, []string, Summary) {
	var initialSize int = len(votes)

	var validVotes []Vote
	var invalidVotes []Vote
	var votedToday []string

	for _, v := range votes {
		if contains(alreadyVoted, v.ONID) || contains(&votedToday, v.ONID) {
			invalidVotes = append(invalidVotes, v)
		} else {
			validVotes = append(validVotes, v)
			votedToday = append(votedToday, v.ONID)
			//TODO make note of submission ID and log it for future reference
		}
	}

	if len(validVotes)+len(invalidVotes) != initialSize {
		log.Fatal("Step 2 vote counts don't match")
	}

	log.Println("Step 2: Invalid votes:", len(invalidVotes))
	log.Println("Step 2: Valid votes:", len(validVotes))

	return validVotes, invalidVotes, votedToday, Summary{len(validVotes) + len(invalidVotes), len(validVotes), len(invalidVotes)}
}
