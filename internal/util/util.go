package util

import (
	"log"
	"maps"
	"slices"
	"strings"
)

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

	crossCheck := vote

	// replace write-in entries with the real candidate
	// vote = strings.ReplaceAll(vote, "ALL CAPS FROM NORMALIZATION", "normal")
	for _, candidate := range candidates {
		vote = strings.ReplaceAll(vote, strings.ToUpper(candidate), candidate)

		// early return if normalized
		if vote != crossCheck {
			return vote
		}
	}

	return vote
}

func RemoveIneligibleWriteins(resultsMap map[string]int, candidates []string, writeInThreshold int) {
	maps.DeleteFunc(resultsMap, func(candidate string, votes int) bool {
		// all good if registered candidate
		if slices.Contains(candidates, candidate) {
			return false
		}

		// if write-in and meets threshold, then all good
		if votes >= writeInThreshold {
			log.Printf("Write-in candidate %s met threshold with %d votes\n", candidate, votes)
			return false
		}
		return true
	})
}
