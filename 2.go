package main

import "log"

//returns valid, invaild, onids that voted and summary
func stepTwo(votes []Vote, alreadyVoted *[]string) ([]Vote, []Vote, []string, Summary) {
	var initialSize int = len(votes)

	var logMessages []string

	var validVotes []Vote
	var invalidVotes []Vote
	var votedToday []string

	for _, v := range votes {
		if contains(alreadyVoted, v.ONID) || contains(&votedToday, v.ONID) {
			invalidVotes = append(invalidVotes, v)
			logMessages = append(logMessages, "Invalid vote from "+v.ONID+" with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05"))
		} else {
			validVotes = append(validVotes, v)
			votedToday = append(votedToday, v.ONID)
		}
	}

	if len(validVotes)+len(invalidVotes) != initialSize {
		log.Fatal("Step 2 vote counts don't match")
	}

	return validVotes, invalidVotes, votedToday, Summary{
		stepInfo:  "Step 2: Dedupe",
		processed: len(validVotes) + len(invalidVotes),
		valid:     len(validVotes),
		invalid:   len(invalidVotes),
		log:       logMessages}
}
