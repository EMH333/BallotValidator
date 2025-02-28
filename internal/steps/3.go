package steps

import (
	"log"

	"ethohampton.com/BallotValidator/internal/util"
)

// returns valid, invaild, and summary
func StepThree(countingConfig *util.CountingConfig, votes []util.Vote, validVotersGraduate, validVotersUndergraduate *[]string) ([]util.Vote, []util.Vote, util.Summary) {
	var initialSize int = len(votes)

	var logMessages []string

	var validVotes []util.Vote
	var invalidVotes []util.Vote

	for _, v := range votes {
		beginningColumns := len(v.Raw)

		if util.Contains(validVotersGraduate, v.ONID) && v.Raw[countingConfig.StepThreeChoiceIndex] != "Graduate Student" {
			invalidVotes = append(invalidVotes, v) //not actually invalid, just copied directly over, valid will actually fix it
			//clear the all rows voting for reps
			start := v.Raw[:countingConfig.StepThreeStartIndex]
			end := v.Raw[countingConfig.StepThreeEndIndex:]
			v.Raw = append(start, make([]string, countingConfig.StepThreeEndIndex-countingConfig.StepThreeStartIndex)...) //nolint:gocritic
			v.Raw = append(v.Raw, end...)
			logMessages = append(logMessages, "Incorrect representatives vote from "+v.ONID+" (supposed to be graduate) with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05")+" was "+v.Raw[countingConfig.StepThreeChoiceIndex])
		} else if util.Contains(validVotersUndergraduate, v.ONID) && v.Raw[countingConfig.StepThreeChoiceIndex] != "Undergraduate Student" {
			invalidVotes = append(invalidVotes, v) //not actually invalid, just copied directly over, valid will actually fix it
			//clear the all rows voting for reps
			start := v.Raw[:countingConfig.StepThreeStartIndex]
			end := v.Raw[countingConfig.StepThreeEndIndex:]
			v.Raw = append(start, make([]string, countingConfig.StepThreeEndIndex-countingConfig.StepThreeStartIndex)...) //nolint:gocritic
			v.Raw = append(v.Raw, end...)
			logMessages = append(logMessages, "Incorrect representatives vote from "+v.ONID+" (supposed to be undergraduate) with response ID "+v.ID+" at "+v.Timestamp.Format("2006-Jan-02 15:04:05")+" was "+v.Raw[countingConfig.StepThreeChoiceIndex])
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
