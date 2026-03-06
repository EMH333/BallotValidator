package util

import "strings"

func RemoveDuplicateOrEmptyStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if item == "" {
			continue
		}
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func NormalizeVote(countingConfig *CountingConfig, candidates []string, vote string) string {
	vote = strings.TrimSpace(vote)
	vote = strings.ToUpper(vote)

	if vote == "WRITE IN:" || vote == "WRITE-IN:" || vote == "WRITE-IN" || vote == "WRITE IN" {
		return ""
	}

	// normalize pres/vp candidates into same form as ballot
	vote = strings.ReplaceAll(vote, " & ", " and ")

	//TODO add removal of entrys found to be invalid
	//TODO add subsitution for valid but ill-formed entrys
	//both should happen after the trim/upper

	crossCheck := vote

	// replace write-in entries with the real candidate
	// vote = strings.ReplaceAll(vote, "ALL CAPS FROM NORMALIZATION", "normal")
	for _, candidate := range candidates {
		vote = strings.ReplaceAll(vote, strings.ToUpper(candidate), candidate)
	}

	if crossCheck != vote {
		return vote
	}

	return vote
}
