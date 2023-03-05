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
	defer f.Close()
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
	defer os.RemoveAll(dir)

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
	defer os.RemoveAll(dir)

	writeFile(dir+"/test.txt", []string{
		"First Name,Last Name,Email,ONID",
		"first1,last1,email1@example.com,onid1",
		"first2,last2,email2@example.com,onid2",
		"first3,last3,email3@example.com,onid3"})

	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "Basic", args: args{fileName: dir + "/test.txt"}, want: []string{"email1@example.com", "email2@example.com", "email3@example.com"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadValidVoters(tt.args.fileName); !reflect.DeepEqual(got, tt.want) {
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
	defer os.RemoveAll(dir)

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

	type args struct {
		folderName string
		upToDay    int64
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "One file", args: args{folderName: dir, upToDay: 1}, want: []string{"a@example.com", "b@example.com"}},
		{name: "Ignore All", args: args{folderName: dir, upToDay: 0}, want: []string{}},
		{name: "Overlaping Files", args: args{folderName: dir, upToDay: 2}, want: []string{"a@example.com", "b@example.com", "c@example.com"}},
		{name: "All", args: args{folderName: dir, upToDay: 10}, want: []string{"a@example.com", "b@example.com", "c@example.com", "d@example.com"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadAlreadyVoted(tt.args.folderName, tt.args.upToDay); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadAlreadyVoted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadVotesCSV(t *testing.T) {
	type args struct {
		fileName string
		startDay int64
		endDay   int64
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
			if got := LoadVotesCSV(tt.args.fileName, tt.args.startDay, tt.args.endDay); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadVotesCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}
