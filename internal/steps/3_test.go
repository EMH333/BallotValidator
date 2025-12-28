package steps

import (
	"testing"
	"time"

	"ethohampton.com/BallotValidator/internal/util"
)

// TODO create valid test for masking vote based on student status now that we can control valid offsets
func TestStepThreeNoConfig(t *testing.T) {
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

func TestStepThreeUndergradWrong(t *testing.T) {
	var votes []util.Vote = []util.Vote{
		{Raw: []string{"Undergraduate Student", "456", "789"}, Timestamp: time.Now(), ONID: "abc", ID: ""},
		{Raw: []string{"Graduate Student", "789", "012"}, Timestamp: time.Now(), ONID: "efg", ID: ""},
		{Raw: []string{"Undergraduate Student", "012", "123"}, Timestamp: time.Now(), ONID: "hij", ID: ""},
		{Raw: []string{"Undergraduate Student", "123", "456"}, Timestamp: time.Now(), ONID: "lmn", ID: ""},
		{Raw: []string{"Undergraduate Student", "123", "456"}, Timestamp: time.Now(), ONID: "opq", ID: ""},
	}

	var validGrad = []string{}
	var validUndergrad = []string{"abc", "efg", "hij", "lmn", "opq"}

	valid, invalid, _ := StepThree(&util.CountingConfig{StepThreeChoiceIndex: 0}, votes, &validGrad, &validUndergrad)

	if len(valid) != len(votes) {
		t.Errorf("Valid voters was %d, expected %d", len(valid), len(votes))
	}

	if len(invalid) != 1 {
		t.Errorf("Invalid voters was %d, expected 1", len(invalid))
	}
}

func TestStepThreeGradWrong(t *testing.T) {
	var votes []util.Vote = []util.Vote{
		{Raw: []string{"Graduate Student", "456", "789"}, Timestamp: time.Now(), ONID: "abc", ID: ""},
		{Raw: []string{"Graduate Student", "789", "012"}, Timestamp: time.Now(), ONID: "efg", ID: ""},
		{Raw: []string{"Undergraduate Student", "012", "123"}, Timestamp: time.Now(), ONID: "hij", ID: ""},
		{Raw: []string{"Graduate Student", "123", "456"}, Timestamp: time.Now(), ONID: "lmn", ID: ""},
		{Raw: []string{"Graduate Student", "123", "456"}, Timestamp: time.Now(), ONID: "opq", ID: ""},
	}

	var validGrad = []string{"abc", "efg", "hij", "lmn", "opq"}
	var validUndergrad = []string{}

	valid, invalid, _ := StepThree(&util.CountingConfig{StepThreeChoiceIndex: 0}, votes, &validGrad, &validUndergrad)

	if len(valid) != len(votes) {
		t.Errorf("Valid voters was %d, expected %d", len(valid), len(votes))
	}

	if len(invalid) != 1 {
		t.Errorf("Invalid voters was %d, expected 1", len(invalid))
	}
}
