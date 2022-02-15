package main

import (
	"testing"
	"time"
)

func TestStepOne(t *testing.T) {
	var validVotersGraduate []string = []string{"123", "456"}
	var validVotersUndergraduate []string = []string{"789", "012"}
	var validVotersUndefined []string = []string{"345", "678"}
	var votes []Vote = []Vote{
		{[]string{"123", "456", "789"}, time.Now(), "123", ""},
		{[]string{"012", "123", "456"}, time.Now(), "056", ""}, //invalid
		{[]string{"456", "789", "012"}, time.Now(), "456", ""},
		{[]string{"012", "123", "456"}, time.Now(), "000", ""}, //invalid
		{[]string{"789", "012", "123"}, time.Now(), "789", ""},
		{[]string{"012", "123", "456"}, time.Now(), "012", ""},
		{[]string{"012", "123", "456"}, time.Now(), "345", ""},
	}

	valid, invalid, _ := stepOne(votes, &validVotersGraduate, &validVotersUndergraduate, &validVotersUndefined)

	if len(valid)+len(invalid) != len(votes) {
		t.Error("Total vote counts don't match")
	}

	if len(valid) != 5 {
		t.Error("Valid vote counts don't match")
	}

	if len(invalid) != 2 {
		t.Error("Invalid vote counts don't match")
	}

	for _, v := range valid {
		switch v.ONID {
		case "123":
			continue
		case "456":
			continue
		case "789":
			continue
		case "012":
			continue
		case "345":
			continue
		default:
			t.Errorf("Invalid id in valid votes: %s\n", v.ONID)
		}
	}

	for _, v := range invalid {
		switch v.ONID {
		case "000":
			continue
		case "056":
			continue
		default:
			t.Errorf("Invalid id in invalid votes: %s\n", v.ONID)
		}
	}
}
