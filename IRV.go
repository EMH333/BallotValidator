package main

import (
	"fmt"
	"strconv"
)

/*
The instant run-off voting functions as follows, for a single race:
- we go through each ballot and put their chioces in an array that lines up with the number of canidates per person
	For example, in a race with A, B, and C as canidates with 1 write in allowed (D and E in this example), a potential ballots would look like [A,C,D,] and [,E,B,A]
- We go through an count up the first choice (index 0 of each ballot) votes for each candidate. nil or empty strings are ignored and removed from the offending ballots
- If (for single-seat races), someone has a majority of the votes (including those who didn't vote in the race?) then they win
- Otherwise, we remove the canidate with the least number of votes (looping through every single entry of every single ballot) and reset to step 2
*/

//TODO deal with ties better

type IRVBallot struct {
	Choices []string
	ID      string
}

func runIRV(votes []Vote, includedCandidates []string, numCandidates, offset int) []string {
	var majority int = (len(votes) / 2) + 1
	var logMessages []string

	//first process into ballots
	ballots, createMessages := createIRVBallots(&votes, includedCandidates, numCandidates, offset)

	logMessages = append(logMessages, createMessages...)
	logMessages = append(logMessages, fmt.Sprint("Number of ballots: ", len(ballots)))

	var roundNumber int = 1
	//now run the IRV
	for {
		logMessages = append(logMessages, "Round "+strconv.Itoa(roundNumber))
		//count up the votes for each candidate
		candidateVotes, ballotsCountedThisRound := countIRVVotes(&ballots)
		logMessages = append(logMessages, "Number of ballots remaining this round: "+fmt.Sprint(ballotsCountedThisRound))
		logMessages = append(logMessages, "----------------------------------------------------")
		winner := ""
		for candidate, votes := range candidateVotes {
			logMessages = append(logMessages, candidate+" has "+strconv.Itoa(votes)+" votes")
			//if someone has over the majority of the vote then they are the winner
			if votes >= majority || len(candidateVotes) <= 1 {
				winner = candidate
			}
		}

		if winner != "" {
			logMessages = append(logMessages, "Winner: "+winner)
			break
		}

		//remove the candidate with the least votes
		var lowestCandidate string
		var lowestVotes int
		for candidate, votes := range candidateVotes {
			if lowestVotes == 0 || votes < lowestVotes {
				lowestCandidate = candidate
				lowestVotes = votes
			}
		}
		logMessages = append(logMessages, "Removing "+lowestCandidate+" from the election")
		//remove the candidate from all ballots
		for i := range ballots {
			for j := range ballots[i].Choices {
				if ballots[i].Choices[j] == lowestCandidate {
					ballots[i].Choices[j] = ""
				}
			}
		}

		roundNumber++
		logMessages = append(logMessages, "", "", "")
	}

	return logMessages
}

func countIRVVotes(ballots *[]IRVBallot) (map[string]int, int) {
	candidateVotes := make(map[string]int)
	ballotsCountedThisRound := 0
	for _, ballot := range *ballots {
		//remove any empty strings from each ballot
		for {
			if len(ballot.Choices) > 0 && ballot.Choices[0] == "" {
				ballot.Choices = ballot.Choices[1:]
			} else {
				break
			}
		}

		// if valid ballot, count the vote
		if len(ballot.Choices) > 0 {
			candidateVotes[ballot.Choices[0]]++
			ballotsCountedThisRound++
		}
	}
	return candidateVotes, ballotsCountedThisRound
}

func createIRVBallots(votes *[]Vote, includedCandidates []string, numCandidates, offset int) ([]IRVBallot, []string) {
	var ballots []IRVBallot
	var logMessages []string
	for _, vote := range *votes {
		var ballot IRVBallot
		ballot.ID = vote.ID
		ballot.Choices = make([]string, numCandidates+1) //include a write-in slot

		var validBallot = true

		// handle preregistered candidates
		for i := offset; i < offset+numCandidates; i++ {
			if vote.Raw[i] != "" {
				//if rank doesn't parse, then they left it blank
				rank, err := strconv.ParseInt(vote.Raw[i], 10, 64)
				if err != nil {
					//logMessages = append(logMessages, "Invalid rank from "+vote.Raw[i]+" for "+vote.ID)
					continue
				}

				//check to make sure we aren't overriding a value
				if ballot.Choices[rank-1] != "" {
					logMessages = append(logMessages, "Error: "+vote.ID+" tried to override "+ballot.Choices[rank-1]+" with "+includedCandidates[i-offset])
					validBallot = false
					break
				}

				ballot.Choices[rank-1] = includedCandidates[i-offset] //set the rank choice to the candidate
			}
		}

		// handle write-in
		rank, err := strconv.ParseInt(vote.Raw[offset+numCandidates], 10, 64)
		//if rank doesn't parse, then they left it blank
		if err != nil {
			//logMessages = append(logMessages, "Invalid rank from "+vote.Raw[offset+numCandidates]+" for "+vote.ID)
		} else {
			//check to make sure we aren't overriding a value
			if ballot.Choices[rank-1] != "" {
				logMessages = append(logMessages, "Error: "+vote.ID+" tried to override "+ballot.Choices[rank-1]+" with "+vote.Raw[offset+numCandidates+1])
				validBallot = false
			}
			writeInName := cleanVote(vote.Raw[offset+numCandidates+1])
			ballot.Choices[rank-1] = writeInName //set the rank choice to the candidate
		}

		if validBallot {
			ballots = append(ballots, ballot)
		} else {
			logMessages = append(logMessages, "Invalid ballot: "+vote.ID)
		}
	}

	return ballots, logMessages
}
