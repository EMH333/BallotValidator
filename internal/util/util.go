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

func CleanVote(vote string) string {
	vote = strings.TrimSpace(vote)

	if vote == "Write in:" || vote == "Write-in:" || vote == "Write-In" {
		return ""
	}

	vote = strings.ToUpper(vote)

	//TODO add removal of entrys found to be invalid
	//TODO add subsitution for valid but ill-formed entrys
	//both should happen after the trim/upper

	// fix spelling error for real person
	vote = strings.ReplaceAll(vote, "CAROYLN PEARCE", "CAROLYN PEARCE")

	// actual candidates for sfc at large, so dedupe those votes
	vote = strings.ReplaceAll(vote, "SOPHIA NOWERS", "Sophia Nowers")
	vote = strings.ReplaceAll(vote, "ABUKAR MOHAMMED", "Abukar Mohamed")
	vote = strings.ReplaceAll(vote, "COLE PETERS", "Cole Peters")

	return vote
}
