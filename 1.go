package main

import "log"

func stepOne(votes []Vote, validVotersGraduate, validVotersUndergraduate, validVotersUndefined *[]string) ([]Vote, []Vote, Summary) {
	var initialSize int = len(votes)

	var validVotes []Vote
	var invalidVotes []Vote

	for _, v := range votes {
		if contains(validVotersGraduate, v.ONID) {
			validVotes = append(validVotes, v)
		} else if contains(validVotersUndergraduate, v.ONID) {
			validVotes = append(validVotes, v)
		} else if contains(validVotersUndefined, v.ONID) {
			validVotes = append(validVotes, v)
		} else {
			invalidVotes = append(invalidVotes, v)
		}
	}

	if len(validVotes)+len(invalidVotes) != initialSize {
		log.Fatal("Step 1 vote counts don't match")
	}

	log.Println("Step 1: Invalid votes:", len(invalidVotes))
	log.Println("Step 1: Valid votes:", len(validVotes))

	return validVotes, invalidVotes, Summary{processed: len(validVotes) + len(invalidVotes), valid: len(validVotes), invalid: len(invalidVotes)}
}
