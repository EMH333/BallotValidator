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
		{
			name: "Don't Delete",
			votes: []util.Vote{
				{
					ID:   "1",
					ONID: "1",
					Raw:  []string{"2023-02-20 17:04:15", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "a,b,c,d,e,f"},
				},
			},
			valid:   1,
			invalid: 0,
			cureCSV: []string{
				"1,no",
			},
		},
		{
			name: "Do Delete",
			votes: []util.Vote{
				{
					ID:   "1",
					ONID: "1",
					Raw:  []string{"2023-02-20 17:04:15", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "a,b,c,d,e,f"},
				},
				{
					ID:   "2",
					ONID: "1",
					Raw:  []string{"2023-02-21 10:04:15", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "a,b,c,d,e,f,g"},
				},
			},
			valid:   1,
			invalid: 1,
			cureCSV: []string{
				"1,yes",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, _ := StepCure(tt.votes)
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
