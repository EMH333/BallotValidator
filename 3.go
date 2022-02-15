package main

import "log"

//start includes first time to remove
//end is first item after the ones to remove
const THREE_START = 31 //TODO
const THREE_END = THREE_START + 27

const THREE_CHOICE = 30 //TODO what column are we looking at to make the choice?

//returns valid, invaild, and summary
func stepThree(votes []Vote, validVotersGraduate, validVotersUndergraduate *[]string) ([]Vote, []Vote, Summary) {
	var initialSize int = len(votes)

	var logMessages []string

	var validVotes []Vote
	var invalidVotes []Vote

	for _, v := range votes {
		if contains(validVotersGraduate, v.ONID) && v.Raw[THREE_CHOICE] != "Graduate Student" {
			invalidVotes = append(invalidVotes, v) //not actually invalid, just copied directly over, valid will actually fix it
			start := v.Raw[:THREE_START]
			end := v.Raw[THREE_END:]
			v.Raw = append(start, make([]string, THREE_END-THREE_START)...)
			v.Raw = append(v.Raw, end...)
			logMessages = append(logMessages, "Incorrect representatives vote from "+v.ONID+" (supposed to be graduate) with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05"))
		} else if contains(validVotersUndergraduate, v.ONID) && v.Raw[THREE_CHOICE] != "Undergraduate Student" {
			invalidVotes = append(invalidVotes, v) //not actually invalid, just copied directly over, valid will actually fix it
			start := v.Raw[:THREE_START]
			end := v.Raw[THREE_END:]
			v.Raw = append(start, make([]string, THREE_END-THREE_START)...)
			v.Raw = append(v.Raw, end...)
			logMessages = append(logMessages, "Incorrect representatives vote from "+v.ONID+" (supposed to be undergraduate) with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05"))
		}

		//add to valid votes regardless (corrections have been made above if needed)
		validVotes = append(validVotes, v)
	}

	// in this case only comparing valid, since no votes should be removed here
	if len(validVotes) != initialSize {
		log.Fatal("Step 3 vote count doesn't match")
	}

	log.Println("Step 3: Invalid votes:", len(invalidVotes))
	log.Println("Step 3: Valid votes:", len(validVotes))

	return validVotes, invalidVotes, Summary{processed: len(validVotes) + len(invalidVotes), valid: len(validVotes), invalid: len(invalidVotes), log: logMessages}
}
