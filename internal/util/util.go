package util

import "strings"

// does the array contain the value?
func Contains(s *[]string, e string) bool {
	for _, a := range *s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func NormalizeVote(countingConfig *CountingConfig, vote string) string {
	vote = strings.TrimSpace(vote)

	if vote == "Write in:" || vote == "Write-in:" || vote == "Write-In" {
		return ""
	}

	vote = strings.ToUpper(vote)

	//TODO add removal of entrys found to be invalid
	//TODO add subsitution for valid but ill-formed entrys
	//both should happen after the trim/upper

	// replace write-in entries with the real candidate
	// vote = strings.ReplaceAll(vote, "ALL CAPS FROM NORMALIZATION", "normal")
	for _, candidate := range countingConfig.CandidatesPresident {
		vote = strings.ReplaceAll(vote, strings.ToUpper(candidate), candidate)
	}
	for _, candidate := range countingConfig.CandidatesSFCChair {
		vote = strings.ReplaceAll(vote, strings.ToUpper(candidate), candidate)
	}
	for _, candidate := range countingConfig.CandidatesSFCAtLarge {
		vote = strings.ReplaceAll(vote, strings.ToUpper(candidate), candidate)
	}
	for _, candidate := range countingConfig.CandidatesGraduateSenate {
		vote = strings.ReplaceAll(vote, strings.ToUpper(candidate), candidate)
	}
	for _, candidate := range countingConfig.CandidatesUndergraduateSenate {
		vote = strings.ReplaceAll(vote, strings.ToUpper(candidate), candidate)
	}

	return vote
}
