package steps

import (
	"log"

	"ethohampton.com/BallotValidator/internal/util"
)

//start includes first time to remove
//end is first item after the ones to remove
const THREE_START = 31
const THREE_END = THREE_START + 27

const THREE_CHOICE = 30

//returns valid, invaild, and summary
func StepThree(votes []util.Vote, validVotersGraduate, validVotersUndergraduate *[]string) ([]util.Vote, []util.Vote, util.Summary) {
	var initialSize int = len(votes)

	var logMessages []string

	var validVotes []util.Vote
	var invalidVotes []util.Vote

	for _, v := range votes {
		beginningColumns := len(v.Raw)

		if util.Contains(validVotersGraduate, v.ONID) && v.Raw[THREE_CHOICE] != "Graduate Student" {
			invalidVotes = append(invalidVotes, v) //not actually invalid, just copied directly over, valid will actually fix it
			//clear the all rows voting for reps
			start := v.Raw[:THREE_START]
			end := v.Raw[THREE_END:]
			v.Raw = append(start, make([]string, THREE_END-THREE_START)...)
			v.Raw = append(v.Raw, end...)
			logMessages = append(logMessages, "Incorrect representatives vote from "+v.ONID+" (supposed to be graduate) with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05"))
		} else if util.Contains(validVotersUndergraduate, v.ONID) && v.Raw[THREE_CHOICE] != "Undergraduate Student" {
			invalidVotes = append(invalidVotes, v) //not actually invalid, just copied directly over, valid will actually fix it
			//clear the all rows voting for reps
			start := v.Raw[:THREE_START]
			end := v.Raw[THREE_END:]
			v.Raw = append(start, make([]string, THREE_END-THREE_START)...)
			v.Raw = append(v.Raw, end...)
			logMessages = append(logMessages, "Incorrect representatives vote from "+v.ONID+" (supposed to be undergraduate) with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05"))
		}

		endColumns := len(v.Raw)

		//make sure the modifications didn't change # of columns which would screw up data
		if beginningColumns != endColumns {
			log.Println("Error: columns changed size during step 3")
		}

		//add to valid votes regardless (corrections have been made above if needed)
		validVotes = append(validVotes, v)
	}

	// in this case only comparing valid, since no votes should be removed here
	if len(validVotes) != initialSize {
		log.Fatal("Step 3 vote count doesn't match")
	}

	return validVotes, invalidVotes, util.Summary{
		StepInfo:  "Step 3: Grad/undergrad",
		Processed: len(validVotes),
		Valid:     len(validVotes),
		Invalid:   len(invalidVotes),
		Log:       logMessages}
}
