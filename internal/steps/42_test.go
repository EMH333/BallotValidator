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
		countPopularityVote(&v, &results, 0, 2, 8)
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

func TestMaxVotesPopularity(t *testing.T) {
	var votes []util.Vote = []util.Vote{
		{Raw: []string{"a,b,c", "", ""}, ONID: "a"},
		{Raw: []string{"a,b,c,d", "", ""}, ONID: "b"}, //too many votes
		{Raw: []string{"a, Write in:", "e", "f"}, ONID: "c"},
		{Raw: []string{"a, Write-in:", "", "f"}, ONID: "d"},
		{Raw: []string{"a,b, Write in:", "e", "f"}, ONID: "e"},
	}
	var results map[string]int = make(map[string]int)

	for _, v := range votes {
		countPopularityVote(&v, &results, 0, 2, 3)
	}

	if results["a"] != 3 {
		t.Errorf("Expected 3 votes for a, got %d", results["a"])
	}

	if results["d"] != 0 {
		t.Errorf("Expected 0 votes for d, got %d", results["d"])
	}

}
