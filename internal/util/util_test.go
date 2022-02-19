package util

import "testing"

// never hurt to do some sanity checks on the most basic stuff
func TestContains(t *testing.T) {
	var testCases = []struct {
		slice  []string
		value  string
		expect bool
	}{
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{"a", "b", "c"}, "d", false},
	}

	for _, testCase := range testCases {
		if Contains(&testCase.slice, testCase.value) != testCase.expect {
			t.Errorf("Expected %v to be %v", testCase.slice, testCase.expect)
		}
	}
}

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
