package util

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func writeFile(fileName string, data []string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // nolint:errcheck // don't care about close
	for _, record := range data {
		_, err = f.WriteString(record + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestLoadStringArrayFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "BallotValidator")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // nolint:errcheck // don't care about not removing files

	writeFile(dir+"/test.txt", []string{"a", "b", "c"})
	writeFile(dir+"/test2.txt", []string{"a", "b", "c", "d, e"})

	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "Basic", args: args{fileName: dir + "/test.txt"}, want: []string{"a", "b", "c"}},
		{name: "With Comma", args: args{fileName: dir + "/test2.txt"}, want: []string{"a", "b", "c", "d, e"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadStringArrayFile(tt.args.fileName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadStringArrayFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadValidVoters(t *testing.T) {
	dir, err := os.MkdirTemp("", "BallotValidator")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // nolint:errcheck // don't care about removing

	writeFile(dir+"/test.txt", []string{
		"First Name,Last Name,Email,ONID,G_UG_STATUS",
		"first1,last1,email1@example.com,onid1,UG",
		"first2,last2,email2@example.com,onid2,UG",
		"first3,last3,email3@example.com,onid3,G"})

	var countingConfig = &CountingConfig{
		ValidVotersFile:        dir + "/test.txt",
		ValidVotersEmailIndex:  2,
		ValidVotersStatusIndex: 4,
	}

	type args struct {
		config    *CountingConfig
		gradorund string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "Basic Graduate", args: args{config: countingConfig, gradorund: "G"}, want: []string{"email3@example.com"}},
		{name: "Basic Undergraduate", args: args{config: countingConfig, gradorund: "UG"}, want: []string{"email1@example.com", "email2@example.com"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadValidVoters(tt.args.config, tt.args.gradorund); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadValidVoters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadAlreadyVoted(t *testing.T) {
	dir, err := os.MkdirTemp("", "BallotValidator")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // nolint:errcheck // don't care about removal

	writeFile(dir+"/alreadyVoted-0-0.csv", []string{
		"a@example.com",
		"b@example.com",
	})

	writeFile(dir+"/alreadyVoted-0-1.csv", []string{
		"a@example.com",
		"b@example.com",
		"c@example.com",
	})

	writeFile(dir+"/alreadyVoted-2-2.csv", []string{
		"d@example.com",
	})

	var countingConfig = &CountingConfig{
		AlreadyVotedDir: dir,
	}

	type args struct {
		config  *CountingConfig
		upToDay int64
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "One file", args: args{config: countingConfig, upToDay: 1}, want: []string{"a@example.com", "b@example.com"}},
		{name: "Ignore All", args: args{config: countingConfig, upToDay: 0}, want: []string{}},
		{name: "Overlaping Files", args: args{config: countingConfig, upToDay: 2}, want: []string{"a@example.com", "b@example.com", "c@example.com"}},
		{name: "All", args: args{config: countingConfig, upToDay: 10}, want: []string{"a@example.com", "b@example.com", "c@example.com", "d@example.com"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadAlreadyVoted(tt.args.config, tt.args.upToDay); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadAlreadyVoted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadVotesCSV(t *testing.T) {
	type args struct {
		config   *CountingConfig
		fileName string
		startDay int
		endDay   int
	}

	tests := []struct {
		name string
		args args
		want []Vote
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadVotesCSV(tt.args.config, tt.args.fileName, tt.args.startDay, tt.args.endDay); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadVotesCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCandidateColumnIndexValidation(t *testing.T) {
	testingConfig := &CountingConfig{
		TallyPresidentOptionsIndex: 0,
		TallySFCChairOptionsIndex:  2,
		CandidatesPresident:        []string{"John", "Alice"},
		CandidatesSFCChair:         []string{"Bob", "Eve"},
	}

	// all correct
	err := candidateColumnIndexValidation(testingConfig, []string{"John", "Alice", "Bob", "Eve"})
	if err != nil {
		t.Errorf("Did not expect an error: %e", err)
	}

	// president wrong
	err = candidateColumnIndexValidation(testingConfig, []string{"Alice", "John", "Bob", "Eve"})
	if err == nil {
		t.Errorf("Expected an error and got none")
	}

	err = candidateColumnIndexValidation(testingConfig, []string{"John", "Alice", "Eve", "Bob"})
	if err == nil {
		t.Errorf("Expected an error and got none")
	}
}
