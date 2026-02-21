package steps

import (
	"fmt"
	"log"
	"maps"
	"os"
	"strings"

	"ethohampton.com/BallotValidator/internal/util"
)

// designed to do all the counting and output a nice little summary
func StepFourtyTwo(countingConfig *util.CountingConfig, votes []util.Vote, outputDirname string) {
	if len(votes) == 0 {
		log.Fatal("No votes to find results for")
	}

	var undergradSenate = make(map[string]int)
	var graduateSenate = make(map[string]int)
	var sfcAtLarge = make(map[string]int)

	for _, vote := range votes {
		///////////////////SENATE/////////////////////////////
		countPopularityVote(countingConfig, &vote, undergradSenate, countingConfig.TallyUndergradeSenateOptionsIndex, countingConfig.TallySenateWritinsCount, countingConfig.TallyUndergraduateSenateWinners)
		countPopularityVote(countingConfig, &vote, graduateSenate, countingConfig.TallyGraduateSenateOptionsIndex, countingConfig.TallySenateWritinsCount, countingConfig.TallyGraduateSenateWinners)

		///////////////////SFC AT LARGE/////////////////////////////
		countPopularityVote(countingConfig, &vote, sfcAtLarge, countingConfig.TallySFCAtLargeOptionsIndex, countingConfig.TallySFCAtLargeWritinsCount, countingConfig.TallySFCAtLargeWinners)
	}
	//log.Println("Counted Popularity Votes")

	//presidental ticket
	presidentResults := util.RunIRV(countingConfig, votes, countingConfig.CandidatesPresident, countingConfig.TallyPresidentOptionsCount, countingConfig.TallyPresidentOptionsIndex)

	//SFC chair
	sfcChairResults := util.RunIRV(countingConfig, votes, countingConfig.CandidatesSFCChair, countingConfig.TallySFCChairOptionsCount, countingConfig.TallySFCChairOptionsIndex)

	_, err := os.Stat(outputDirname)
	if os.IsNotExist(err) && os.Mkdir(outputDirname, 0o755) != nil {
		log.Fatal("Could not create output directory", outputDirname)
	}

	//write to senate file
	writeMultipleVoteResults(undergradSenate, outputDirname+"/undergradSenate.csv")
	writeMultipleVoteResults(graduateSenate, outputDirname+"/graduateSenate.csv")

	//write to SFC At-large file
	writeMultipleVoteResults(sfcAtLarge, outputDirname+"/sfc-at-large.csv")

	//write to president file
	writeIRVResults(presidentResults, outputDirname+"/president.txt")

	//write to SFC chair file
	writeIRVResults(sfcChairResults, outputDirname+"/sfc-chair.txt")
}

func countPopularityVote(countingConfig *util.CountingConfig, vote *util.Vote, position map[string]int, initialPosition, numWriteins, maxVotes int) {
	rawVotes := strings.Split(vote.Raw[initialPosition], ",")
	var votes []string
	// clean up the write in entries
	for _, vote := range rawVotes {
		vote = strings.TrimSpace(vote)
		if vote != "Write in:" && vote != "Write-in:" && vote != "Write-In" {
			votes = append(votes, vote)
		}
	}

	for i := range numWriteins {
		wi := vote.Raw[initialPosition+1+i]
		if wi != "" {
			//TODO seperate commas here before normalizing votes
			votes = append(votes, util.NormalizeVote(countingConfig, wi))
		}
	}

	votes = util.RemoveDuplicateStr(votes)

	// can't pick more than the max for these positions
	if len(votes) > maxVotes {
		log.Println("WARNING: More than the max number of votes for", vote.ONID)
		return
	}

	for _, v := range votes {
		if v != "" {
			position[v]++
		}
	}
}

func writeMultipleVoteResults(results map[string]int, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // nolint:errcheck // don't care about close

	_, err = f.WriteString("Candidate,Votes\n")
	if err != nil {
		log.Fatal(err)
	}

	//print results in order of value
	//resultsCopy results to a new map so we can delete entries
	var resultsCopy = make(map[string]int)
	maps.Copy(resultsCopy, results)

	for len(resultsCopy) > 0 {
		var maxNum = 0
		var maxKey = ""
		for k, v := range resultsCopy {
			// sort by alphabetical order if same number of votes
			if v > maxNum || (v >= maxNum && k < maxKey) {
				maxNum = v
				maxKey = k
			}
		}
		_, err = f.WriteString("\"" + maxKey + "\"" + "," + fmt.Sprint(maxNum) + "\n")
		if err != nil {
			log.Fatal(err)
		}
		delete(resultsCopy, maxKey)
	}

	err = f.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func writeIRVResults(results []string, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // nolint:errcheck // don't care about close

	for _, v := range results {
		_, err = f.WriteString(v + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	err = f.Sync()
	if err != nil {
		log.Fatal(err)
	}
}
