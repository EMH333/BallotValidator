package steps

import (
	"strings"
	"testing"

	"ethohampton.com/BallotCleaner/internal/util"
)

func TestCountPopularityVote(t *testing.T) {
	var votes []util.Vote = []util.Vote{
		{Raw: []string{"a,b,c,d,Write in:", "", ""}},
		{Raw: []string{"a,b,c,d", "", ""}},
		{Raw: []string{"Write in:", "e", "f"}},
		{Raw: []string{"Write-in:", "", "f"}},
		{Raw: []string{"", "e", ""}},
		{Raw: []string{"a,b,c,d,Write-in:", "e", "f"}},
		{Raw: []string{"a", "Something super crazy 'logan' !@#$$%^&*()_+=-;:'|][{},./<>?", ""}},
		{Raw: []string{"b,c", "", ""}},
		{Raw: []string{"d,Write in:,Write in:", "f", "e"}},
	}

	var results map[string]int = make(map[string]int)

	for _, v := range votes {
		countPopularityVote(&v, &results, 0, 2)
	}

	for k, v := range results {
		if strings.HasPrefix(k, "SOMETHING") {
			if v != 1 {
				t.Errorf("Expected 1 vote for %s, got %d", k, v)
			}
			continue
		}
		if v != 4 {
			t.Errorf("Expected 4 votes for %s, got %d", k, v)
		}
	}
}
