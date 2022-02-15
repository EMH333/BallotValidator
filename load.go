package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var EPOCH, epochErr = time.Parse("2006-Jan-02 03:04:05", "2022-Feb-14 00:00:01")

// values to use when importing from csv
const IMPORT_TIMESTAMP = 1 //using end date so it is consistent across submission times
const IMPORT_ONID = 75
const IMPORT_COMPLETE = 6

func loadVotesCSV(fileName string, startDay, endDay int) []Vote {
	// make sure our epoch is valid
	if epochErr != nil {
		log.Fatal(epochErr)
	}

	var validStartTime = EPOCH.Add(time.Duration(startDay) * 24 * time.Hour)
	var validEndTime = EPOCH.Add(time.Duration(endDay+1) * 24 * time.Hour) // add one day to end day

	var votes []Vote
	//return []string{"TODO"}
	//load csv file
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// read csv values using csv.Reader
	//with modifications to handle the specifics of the valid votes list
	csvReader := csv.NewReader(f)
	csvReader.Comma = ','
	csvReader.TrimLeadingSpace = true

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		//skip the first few rows
		if rec[IMPORT_TIMESTAMP] == "EndDate" || rec[IMPORT_TIMESTAMP] == "End Date" || strings.Contains(rec[IMPORT_TIMESTAMP], "ImportId") {
			continue
		}

		timestamp, err := time.Parse("1/2/2006 15:04", rec[IMPORT_TIMESTAMP]) //2/14/2022 9:10
		if err != nil {
			log.Fatal(err)
		}

		//make sure it is only reading the correct day
		if timestamp.Before(validStartTime) || timestamp.After(validEndTime) {
			continue
		}

		ONID := rec[IMPORT_ONID]
		//sanity check to make sure the ONID looks like an email
		if !strings.Contains(ONID, "@oregonstate.edu") {
			log.Fatalf("ONID is not an email address: %s\n", ONID)
		}

		//make sure it is a complete row
		if rec[IMPORT_COMPLETE] != "TRUE" {
			log.Fatalf("Vote is not complete: %+v\n", rec)
		}

		//append rec to votes
		votes = append(votes, Vote{Raw: rec, Timestamp: timestamp, ONID: ONID})
	}

	return votes
}
