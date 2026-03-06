package util

import (
	"reflect"
	"testing"
)

func TestRemoveDuplicateStr(t *testing.T) {
	var testCases = []struct {
		slice  []string
		expect []string
	}{
		{[]string{"a", "b", "c", "a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
	}

	for _, testCase := range testCases {
		results := RemoveDuplicateOrEmptyStr(testCase.slice)
		if len(results) != len(testCase.expect) {
			t.Errorf("Expected length %d to be %d", len(results), len(testCase.expect))
		}
		for i := range results {
			if results[i] != testCase.expect[i] {
				t.Errorf("Expected %v to be %v", results, testCase.expect)
				return
			}
		}
	}
}

func TestNormalizeVote(t *testing.T) {
	var testCases = []struct {
		vote   string
		expect string
	}{
		{"a", "A"},
		{"a ", "A"},
		{" A", "A"},
		{"b", "b"},
		{"Write in:", ""},
		{"Write-in:", ""},
		{"Write-In", ""},
	}

	for _, testCase := range testCases {
		if NormalizeVote(&CountingConfig{}, []string{"b"}, testCase.vote) != testCase.expect {
			t.Errorf("Expected %v to be %v", testCase.vote, testCase.expect)
		}
	}
}

func TestRemoveIneligibleWriteins(t *testing.T) {
	tests := []struct {
		name             string
		resultsMap       map[string]int
		candidates       []string
		writeInThreshold int
		expected         map[string]int
	}{
		{
			name: "No candidates to remove",
			resultsMap: map[string]int{
				"a": 5,
				"b": 10,
				"c": 15,
			},
			candidates:       []string{"a", "b", "c"},
			writeInThreshold: 20,
			expected: map[string]int{
				"a": 5,
				"b": 10,
				"c": 15,
			},
		},
		{
			name: "Remove all writeins",
			resultsMap: map[string]int{
				"a": 5,
				"b": 10,
				"c": 15,
				"d": 3,
				"e": 4,
			},
			candidates:       []string{"a", "b", "c"},
			writeInThreshold: 20,
			expected: map[string]int{
				"a": 5,
				"b": 10,
				"c": 15,
			},
		},
		{
			name: "Eligible write-in",
			resultsMap: map[string]int{
				"a": 5,
				"b": 10,
				"c": 15,
				"d": 30,
				"e": 4,
			},
			candidates:       []string{"a", "b", "c"},
			writeInThreshold: 20,
			expected: map[string]int{
				"a": 5,
				"b": 10,
				"c": 15,
				"d": 30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RemoveIneligibleWriteins(tt.resultsMap, tt.candidates, tt.writeInThreshold)
			if !reflect.DeepEqual(tt.resultsMap, tt.expected) {
				t.Errorf("Expected %v to be %v", tt.resultsMap, tt.expected)
			}
		})
	}
}
