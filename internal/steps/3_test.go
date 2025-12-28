package steps

import (
	"testing"
	"time"

	"ethohampton.com/BallotValidator/internal/util"
)

// TODO create valid test now that we can control valid offsets
func TestStepThree(t *testing.T) {
	var votes []util.Vote = []util.Vote{
		{Raw: []string{"123", "456", "789"}, Timestamp: time.Now(), ONID: "abc", ID: ""},
		{Raw: []string{"456", "789", "012"}, Timestamp: time.Now(), ONID: "efg", ID: ""},
		{Raw: []string{"789", "012", "123"}, Timestamp: time.Now(), ONID: "hij", ID: ""},
		{Raw: []string{"012", "123", "456"}, Timestamp: time.Now(), ONID: "lmn", ID: ""},
		{Raw: []string{"012", "123", "456"}, Timestamp: time.Now(), ONID: "opq", ID: ""},
	}

	var validGrad = []string{}
	var validUndergrad = []string{}

	valid, invalid, _ := StepThree(&util.CountingConfig{}, votes, &validGrad, &validUndergrad)

	if len(valid)+len(invalid) != len(votes) {
		t.Error("Total vote counts don't match")
	}

	if len(invalid) != 0 {
		t.Error("Received invalid voters unexpectedly")
	}
}
