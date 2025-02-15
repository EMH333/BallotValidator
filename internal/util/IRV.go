package util

import (
	"fmt"
	"log"
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

type IRVBallot struct {
	Choices []string
	ID      string
}

func RunIRV(countingConfig *CountingConfig, votes []Vote, includedCandidates []string, numCandidates, offset int) []string {
	if len(votes) == 0 {
		return []string{"No votes cast"}
	}

	var logMessages []string

	//first process into ballots
	ballots, createMessages := createIRVBallots(countingConfig, &votes, includedCandidates, numCandidates, offset)

	logMessages = append(logMessages, createMessages...)
	logMessages = append(logMessages, fmt.Sprint("Number of ballots: ", len(ballots)), "")

	var roundNumber int = 1
	//now run the IRV
	for {
		logMessages = append(logMessages, "Round "+strconv.Itoa(roundNumber))
		//count up the votes for each candidate
		candidateVotes, ballotsCountedThisRound := countIRVVotes(&ballots)
		var majority int = (ballotsCountedThisRound / 2) + 1 // as per statute, the majority is based on the number of votes cast in the round, not overall

		logMessages = append(logMessages, "Number of ballots remaining this round: "+fmt.Sprint(ballotsCountedThisRound))
		logMessages = append(logMessages, "----------------------------------------------------")

		//copy candidateVotes to a new map so we can delete entries as we print them
		var candidateVotesCopy map[string]int = make(map[string]int)
		for candidate, votes := range candidateVotes {
			candidateVotesCopy[candidate] = votes
		}

		//print candidates in order of votes
		for len(candidateVotesCopy) > 0 {
			var max int = 0
			var maxKey string = ""
			for candidate, votes := range candidateVotesCopy {
				// sort by alphabetical order if same number of votes
				if votes > max || (votes >= max && candidate < maxKey) {
					max = votes
					maxKey = candidate
				}
			}
			logMessages = append(logMessages, maxKey+" has "+strconv.Itoa(max)+" votes")
			delete(candidateVotesCopy, maxKey)
		}

		//check if there is a winner yet
		winner := ""
		for candidate, votes := range candidateVotes {
			//if someone has over the majority of the vote then they are the winner
			if votes >= majority || len(candidateVotes) <= 1 {
				// a check here, since only one candidate should ever have a majority of the votes
				if winner != "" {
					log.Fatal("Multiple candidates have a majority of the votes in a single round. This should never happen.")
				}

				winner = candidate
			}
		}

		if winner != "" {
			logMessages = append(logMessages, "", "Winner: "+winner+" with "+strconv.Itoa(candidateVotes[winner])+" votes"+" which is "+strconv.FormatFloat(float64(candidateVotes[winner]*100)/float64(ballotsCountedThisRound), 'f', 2, 64)+"% of the vote")
			break
		}

		//remove the candidate with the least votes
		//if there is a tie, the candidate with the lowest alphabetical order is removed //TODO confirm this is the correct way to handle ties
		var lowestCandidate string
		var lowestVotes int
		var secondLowestVotes int
		for candidate, votes := range candidateVotes {
			if lowestVotes == 0 || votes < lowestVotes || (votes <= lowestVotes && candidate < lowestCandidate) {
				lowestCandidate = candidate
				lowestVotes = votes
			}
		}
		//find the second lowest number of votes (AFTER finding the absolute lowest)
		for _, votes := range candidateVotes {
			if votes > lowestVotes && (secondLowestVotes == 0 || votes < secondLowestVotes) {
				secondLowestVotes = votes
			}
		}

		// loop through and see if there are multiple canidates w/ the lowest number of votes
		var numWithLowestVotes int
		var lowestCandidates []string
		for c, v := range candidateVotes {
			if v == lowestVotes {
				numWithLowestVotes++
				lowestCandidates = append(lowestCandidates, c)
			}
		}

		logMessages = append(logMessages, "", "Lowest number of votes: "+strconv.Itoa(lowestVotes))
		//logMessages = append(logMessages, "Second lowest number of votes: "+strconv.Itoa(secondLowestVotes))
		//logMessages = append(logMessages, "Number of candidates with lowest number of votes: "+strconv.Itoa(numWithLowestVotes))
		//if we can remove all the lowest candidates without affecting the other results, then do it
		if lowestVotes == 1 && (numWithLowestVotes*lowestVotes) < secondLowestVotes {
			for _, c := range lowestCandidates {
				logMessages = append(logMessages, "Removing "+c+" with "+fmt.Sprint(lowestVotes)+" from the election")
				removeFromBallots(&ballots, c)
			}
		} else {
			//othwerwise, remove the lowest candidate (using algorithm from above)
			logMessages = append(logMessages, "Removing "+lowestCandidate+" with "+fmt.Sprint(lowestVotes)+" from the election")
			removeFromBallots(&ballots, lowestCandidate)
		}

		roundNumber++
		logMessages = append(logMessages, "", "", "")
	}

	return logMessages
}

func removeFromBallots(ballots *[]IRVBallot, candidate string) {
	//remove the candidate from all ballots
	for i := range *ballots {
		for j := range (*ballots)[i].Choices {
			if (*ballots)[i].Choices[j] == candidate {
				(*ballots)[i].Choices[j] = ""
			}
		}
	}
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

func createIRVBallots(countingConfig *CountingConfig, votes *[]Vote, includedCandidates []string, numCandidates, offset int) ([]IRVBallot, []string) {
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
					logMessages = append(logMessages, "Invalid rank from "+vote.Raw[i]+" for "+vote.ID)
					validBallot = false
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
		writeinRank := vote.Raw[offset+numCandidates]
		if writeinRank != "" {
			rank, err := strconv.ParseInt(writeinRank, 10, 64)
			if err != nil {
				logMessages = append(logMessages, "Invalid rank from "+vote.Raw[offset+numCandidates]+" for "+vote.ID)
				validBallot = false
			} else {
				//check to make sure we aren't overriding a value
				if ballot.Choices[rank-1] != "" {
					logMessages = append(logMessages, "Error: "+vote.ID+" tried to override "+ballot.Choices[rank-1]+" with "+vote.Raw[offset+numCandidates+1])
					validBallot = false
				}

				writeInName := NormalizeVote(countingConfig, vote.Raw[offset+numCandidates+1])
				ballot.Choices[rank-1] = writeInName //set the rank choice to the candidate
			}
		}

		if validBallot {
			ballots = append(ballots, ballot)
		} else {
			logMessages = append(logMessages, "Invalid ballot: "+vote.ID)
		}
	}

	return ballots, logMessages
}
