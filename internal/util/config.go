package util

import (
	"encoding/json"
	"log"
	"os"
)

type CountingConfig struct {
	ElectionEpoch     string // used to verify start/end times make sense
	ElectionStartTime string
	ElectionEndTime   string
	BallotTimeFormat  string
	ElectionNumDays   int // used to verify when program should print results instead of just stats

	// file to load valid voters from
	ValidVotersFile        string
	ValidVotersEmailIndex  int
	ValidVotersStatusIndex int
	// for non-results runs, where to find list of voters who have already voted
	AlreadyVotedDir string

	// csv columns to use when importing from csv
	ImportTimestamp int //using end date so it is consistent across submission times
	ImportType      int
	ImportONID      int
	ImportComplete  int
	ImportID        int

	// config for grad vs undergrad ballot selection
	// start includes first time to remove
	// end is first item after the ones to remove
	StepThreeChoiceIndex int
	StepThreeStartIndex  int
	StepThreeEndIndex    int

	// csv columns and offsets for each position
	TallyUndergradeSenateOptionsIndex int
	TallyGraduateSenateOptionsIndex   int
	TallySenateWritinsCount           int // number of write-in columns

	TallySFCAtLargeOptionsIndex int
	TallySFCAtLargeWritinsCount int

	TallyPresidentOptionsIndex int
	TallyPresidentOptionsCount int // number of presidential candidates

	TallySFCChairOptionsIndex int
	TallySFCChairOptionsCount int // number of SFC Chair candidates

	// number of winners
	TallyUndergraduateSenateWinners int
	TallyGraduateSenateWinners      int
	TallySFCAtLargeWinners          int

	// Candidates
	CandidatesPresident           []string
	CandidatesSFCChair            []string
	CandidatesSFCAtLarge          []string
	CandidatesUndergraduateSenate []string
	CandidatesGraduateSenate      []string
}

func LoadCountingConfig(location string) CountingConfig {
	fileContent, err := os.ReadFile(location)

	if err != nil {
		log.Fatal("Error when opening config file: ", err)
	}

	var payload CountingConfig
	err = json.Unmarshal(fileContent, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return payload
}
