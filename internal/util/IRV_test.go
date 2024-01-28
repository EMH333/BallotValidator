package util

import (
	"reflect"
	"strings"
	"testing"
)

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
			{ID: "i", Raw: []string{"1", "3", "2", "2", "e"}}, // shouldn't be counted because of duplicate 2's
		}, candidates: []string{"a", "b", "c"}, winner: "Winner: a with 3 votes which is 60.00% of the vote"},
		{
			// test that a tie is broken by the number of votes for the next choice
			votes: []Vote{
				{ID: "a", Raw: []string{"1", "2", "3", "", ""}},
				{ID: "b", Raw: []string{"2", "1", "3", "", ""}},
				{ID: "c", Raw: []string{"3", "2", "1", "", ""}},
				{ID: "d", Raw: []string{"2", "", "", "1", "d"}},
				{ID: "e", Raw: []string{"", "2", "", "1", "e"}},
			},
			candidates: []string{"a", "b", "c"},
			// TODO this will change if we change the way we break last place ties
			winner: "Winner: b with 3 votes which is 60.00% of the vote",
		},
	}

	for _, tc := range testCases {
		logMessages := RunIRV(tc.votes, tc.candidates, len(tc.candidates), 0)
		if !Contains(&logMessages, tc.winner) {
			t.Errorf("Expected %s, got:\n%s\n\n", tc.winner, strings.Join(logMessages, "\n"))
		}
	}
}

func TestCreateIRVBallots(t *testing.T) {
	votes := []Vote{
		{ID: "a", Raw: []string{"1", "2", "3", "", ""}},
		{ID: "b", Raw: []string{"2", "1", "3", "", ""}},
		{ID: "c", Raw: []string{"3", "2", "1", "", ""}},
		{ID: "d", Raw: []string{"1", "", "", "2", "d"}},
		{ID: "e", Raw: []string{"", "", "", "", ""}},
		{ID: "f", Raw: []string{"", "", "", "1", "e"}},
		{ID: "g", Raw: []string{"1", "", "3", "", ""}},
		{ID: "h", Raw: []string{"1", "1", "2", "3", "e"}},
		{ID: "i", Raw: []string{"1", "3", "2", "2", "e"}},
	}

	includedCandidates := []string{"a", "b", "c"}
	numCandidates := len(includedCandidates)
	offset := 0

	expectedBallots := []IRVBallot{
		{ID: "a", Choices: []string{"a", "b", "c", ""}},
		{ID: "b", Choices: []string{"b", "a", "c", ""}},
		{ID: "c", Choices: []string{"c", "b", "a", ""}},
		{ID: "d", Choices: []string{"a", "D", "", ""}}, // capital D for write-in
		{ID: "e", Choices: []string{"", "", "", ""}},
		{ID: "f", Choices: []string{"E", "", "", ""}}, // capital E for write-in
		{ID: "g", Choices: []string{"a", "", "c", ""}},
	}

	expectedLogMessages := []string{
		"Error: h tried to override a with b",
		"Invalid ballot: h",
		"Error: i tried to override c with e",
		"Invalid ballot: i",
	}

	ballots, logMessages := createIRVBallots(&votes, includedCandidates, numCandidates, offset)

	if len(ballots) != len(expectedBallots) {
		t.Errorf("Expected %d ballots, got %d", len(expectedBallots), len(ballots))
	}

	for i, expectedBallot := range expectedBallots {
		if !reflect.DeepEqual(ballots[i], expectedBallot) {
			t.Errorf("Expected ballot %+v, got %+v", expectedBallot, ballots[i])
		}
	}

	if len(logMessages) != len(expectedLogMessages) {
		t.Errorf("Expected %d log messages, got %d", len(expectedLogMessages), len(logMessages))
	}

	for i, expectedLogMessage := range expectedLogMessages {
		if logMessages[i] != expectedLogMessage {
			t.Errorf("Expected log message '%s', got '%s'", expectedLogMessage, logMessages[i])
		}
	}
}

func TestRemoveFromBallots(t *testing.T) {
	ballots := []IRVBallot{
		{Choices: []string{"a", "b", "c"}, ID: "a"},
		{Choices: []string{"b", "a", "c"}, ID: "b"},
		{Choices: []string{"c", "b", "a"}, ID: "c"},
		{Choices: []string{"", "d", "a"}, ID: "d"},
		{Choices: []string{"", "", ""}, ID: "e"},
		{Choices: []string{"", "e", ""}, ID: "f"},
		{Choices: []string{"", "", "a"}, ID: "g"},
	}

	expectedBallots := []IRVBallot{
		{Choices: []string{"", "b", "c"}, ID: "a"},
		{Choices: []string{"b", "", "c"}, ID: "b"},
		{Choices: []string{"c", "b", ""}, ID: "c"},
		{Choices: []string{"", "d", ""}, ID: "d"},
		{Choices: []string{"", "", ""}, ID: "e"},
		{Choices: []string{"", "e", ""}, ID: "f"},
		{Choices: []string{"", "", ""}, ID: "g"},
	}

	removeFromBallots(&ballots, "a")
	removeFromBallots(&ballots, "g")

	if len(ballots) != len(expectedBallots) {
		t.Errorf("Expected %d ballots, got %d", len(expectedBallots), len(ballots))
	}

	for i, expectedBallot := range expectedBallots {
		if !reflect.DeepEqual(ballots[i], expectedBallot) {
			t.Errorf("Expected ballot %+v, got %+v", expectedBallot, ballots[i])
		}
	}
}
