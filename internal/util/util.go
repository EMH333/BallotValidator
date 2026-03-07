package util

import (
	"iter"
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

	if strings.EqualFold(vote, "WRITE IN:") || strings.EqualFold(vote, "WRITE-IN:") ||
		strings.EqualFold(vote, "WRITE-IN") || strings.EqualFold(vote, "WRITE IN") {
		return ""
	}

	// normalize pres/vp candidates into same form as ballot
	vote = strings.Replace(vote, " & ", " and ", 1)

	// replace write-in entries with the real candidate
	for _, candidate := range candidates {
		if strings.EqualFold(vote, candidate) {
			return candidate
		}
	}

	return strings.ToUpper(vote)
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

func StringSliceToMap(sl []string) iter.Seq2[string, struct{}] {
	return func(yield func(string, struct{}) bool) {
		for _, k := range sl {
			if !yield(k, struct{}{}) {
				return
			}
		}
	}
}
