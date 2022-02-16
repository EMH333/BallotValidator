package main

import "log"

func stepOne(votes []Vote, validVotersGraduate, validVotersUndergraduate, validVotersUndefined *[]string) ([]Vote, []Vote, Summary) {
	var initialSize int = len(votes)

	var messageLog []string

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
			messageLog = append(messageLog, "Invalid vote from "+v.ONID+" with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05"))
		}
	}

	if len(validVotes)+len(invalidVotes) != initialSize {
		log.Fatal("Step 1 vote counts don't match")
	}

	return validVotes, invalidVotes, Summary{
		stepInfo:  "Step 1: Valid voter",
		processed: len(validVotes) + len(invalidVotes),
		valid:     len(validVotes),
		invalid:   len(invalidVotes),
		log:       messageLog}
}
