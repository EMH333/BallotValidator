package main

import "strings"

// does the array contain the value?
func contains(s *[]string, e string) bool {
	for _, a := range *s {
		if a == e {
			return true
		}
	}
	return false
}

func removeDuplicateStr(strSlice []string) []string {
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

func cleanVote(vote string) string {
	if vote == "Write in:" || vote == "Write-in:" {
		return ""
	}

	vote = strings.TrimSpace(vote)
	vote = strings.ToUpper(vote)

	//TODO add removal of entrys found to be invalid
	//TODO add subsitution for valid but ill-formed entrys
	//both should happen after the trim/upper

	return vote
}
