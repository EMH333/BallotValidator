package main

import "log"

func stepOne(votes []Vote, validVotersGraduate, validVotersUndergraduate []string) ([]Vote, []Vote, Summary) {
	var initialSize int = len(votes)

	var validVotes []Vote
	var invalidVotes []Vote

	for _, v := range votes {
		if contains(validVotersGraduate, v.ONID) {
			validVotes = append(validVotes, v)
		} else if contains(validVotersUndergraduate, v.ONID) {
			validVotes = append(validVotes, v)
		} else {
			invalidVotes = append(invalidVotes, v)
		}
	}

	if len(validVotes)+len(invalidVotes) != initialSize {
		log.Println("Step 1: Invalid votes:", len(invalidVotes))
		log.Println("Step 2: Valid votes:", len(validVotes))
		log.Fatal("Step 1 vote counts don't match")
	}

	return validVotes, invalidVotes, Summary{len(validVotes) + len(invalidVotes), len(validVotes), len(invalidVotes)}
}
