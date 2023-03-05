package util

import "testing"

func TestCountIRVVotes(t *testing.T) {
	ballots := []IRVBallot{
		{Choices: []string{"a", "b", "c"}, ID: "a"},
		{Choices: []string{"b", "a", "c"}, ID: "b"},
		{Choices: []string{"c", "b", "a"}, ID: "c"},
		{Choices: []string{"", "d", "a"}, ID: "d"},
		{Choices: []string{"", "", ""}, ID: "e"},
		{Choices: []string{"", "e", ""}, ID: "f"},
		{Choices: []string{"", "", "a"}, ID: "g"},
	}

	results, ballotsCounted := countIRVVotes(&ballots)

	if ballotsCounted != 6 {
		t.Errorf("Expected 6 ballots counted, got %d", ballotsCounted)
	}

	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}

	if results["a"] != 2 {
		t.Errorf("Expected 2 votes for a, got %d", results["a"])
	}

	for _, l := range []string{"b", "c", "d", "e"} {
		if results[l] != 1 {
			t.Errorf("Expected 1 vote for %s, got %d", l, results[l])
		}
	}
}

func TestOverallIRV(t *testing.T) {
	testCases := []struct {
		votes      []Vote
		candidates []string
		winner     string
	}{
		{votes: []Vote{
			{ID: "a", Raw: []string{"1", "2", "3", "", ""}},
			{ID: "b", Raw: []string{"2", "1", "3", "", ""}},
			{ID: "c", Raw: []string{"3", "2", "1", "", ""}},
			{ID: "d", Raw: []string{"1", "", "", "2", "d"}},
			{ID: "e", Raw: []string{"", "", "", "", ""}},
			{ID: "f", Raw: []string{"", "", "", "1", "e"}},
			{ID: "g", Raw: []string{"1", "", "3", "", ""}},
			{ID: "h", Raw: []string{"1", "1", "2", "3", "e"}}, // shouldn't be counted because of duplicate 1's
			{ID: "i", Raw: []string{"1", "3", "2", "1", "e"}}, // shouldn't be counted because of duplicate 1's
		}, candidates: []string{"a", "b", "c"}, winner: "a"},
	}

	for _, tc := range testCases {
		logMessages := RunIRV(tc.votes, tc.candidates, len(tc.candidates), 0)
		if !Contains(&logMessages, "Winner: "+tc.winner+" with 3 votes which is 60.00% of the vote") {
			t.Errorf("Expected winner %s, got %s", tc.winner, logMessages)
		}
	}
}
