package util

import "testing"

func TestRemoveDuplicateStr(t *testing.T) {
	var testCases = []struct {
		slice  []string
		expect []string
	}{
		{[]string{"a", "b", "c", "a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
	}

	for _, testCase := range testCases {
		results := RemoveDuplicateStr(testCase.slice)
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
		{"Write in:", ""},
		{"Write-in:", ""},
		{"Write-In", ""},
	}

	for _, testCase := range testCases {
		if NormalizeVote(&CountingConfig{}, testCase.vote) != testCase.expect {
			t.Errorf("Expected %v to be %v", testCase.vote, testCase.expect)
		}
	}
}
