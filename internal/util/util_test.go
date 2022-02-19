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
