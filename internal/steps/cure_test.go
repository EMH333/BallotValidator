package steps

import (
	"reflect"
	"testing"

	"ethohampton.com/BallotCleaner/internal/util"
)

func TestStepCure(t *testing.T) {
	tests := []struct {
		name    string
		votes   []util.Vote
		valid   int
		invalid int
		cureCSV []string
	}{
		//TODO add more tests here
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, _ := StepCure(tt.votes, "data/senateVotes.csv")
			if len(got) != tt.valid {
				t.Errorf("StepCure() valid = %v, want %v", got, tt.valid)
			}
			if len(got1) != tt.invalid {
				t.Errorf("StepCure() invalid = %v, want %v", got1, tt.invalid)
			}
			if !reflect.DeepEqual(got2, tt.cureCSV) {
				t.Errorf("StepCure() cureCSV = %v, want %v", got2, tt.cureCSV)
			}
		})
	}
}
