package steps

import (
	"testing"
	"time"

	"ethohampton.com/BallotValidator/internal/util"
)

func TestStepOne(t *testing.T) {
	var validVotersGraduate = []string{"123", "456"}
	var validVotersUndergraduate = []string{"789", "012"}
	var validVotersUndefined = []string{"345", "678"}

	var votes = []util.Vote{
		{Raw: []string{"123", "456", "789"}, Timestamp: time.Now(), ONID: "123", ID: ""},
		{Raw: []string{"012", "123", "456"}, Timestamp: time.Now(), ONID: "056", ID: ""}, //invalid
		{Raw: []string{"456", "789", "012"}, Timestamp: time.Now(), ONID: "456", ID: ""},
		{Raw: []string{"012", "123", "456"}, Timestamp: time.Now(), ONID: "000", ID: ""}, //invalid
		{Raw: []string{"789", "012", "123"}, Timestamp: time.Now(), ONID: "789", ID: ""},
		{Raw: []string{"012", "123", "456"}, Timestamp: time.Now(), ONID: "012", ID: ""},
		{Raw: []string{"012", "123", "456"}, Timestamp: time.Now(), ONID: "345", ID: ""},
	}

	valid, invalid, _ := StepOne(votes, &validVotersGraduate, &validVotersUndergraduate, &validVotersUndefined)

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
