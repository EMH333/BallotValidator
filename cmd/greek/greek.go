package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"ethohampton.com/BallotCleaner/internal/steps"
	"ethohampton.com/BallotCleaner/internal/util"
)

//determine percentage of each fraternity and sorority that has voted
// load votes
// load fraternity and sorority data
// for each member, determine if they have voted
// if they have voted then add them to the appropriate fraternity or sorority
// regardless, update total members for that fraternity or sorority

const START = 0
const END = 100 //no harm in going overboard here

type Member struct {
	ONID         string
	Organization string
}

type Organization struct {
	Total int
	Voted int
}

func main() {
	var dataFile string
	// in the form of `program <file_to_process>`
	if len(os.Args) == 2 {
		dataFile = os.Args[1]
	} else {
		log.Fatal("Need to specify a file to process")
	}

	_, err := os.Stat("output")
	if os.IsNotExist(err) && os.Mkdir("output", 0755) != nil {
		log.Fatal("Could not create output directory")
	}

	// Load the votes
	log.Println("Loading votes...")
	votes := util.LoadVotesCSV("data/ballots/"+dataFile, START, END, util.IMPORT_ONID)
	log.Printf("%d votes loaded\n", len(votes))

	// reuse step two to get the ONID emails for all eligible students that have already voted
	_, _, alreadyVoted, _ := steps.StepTwo(votes, &[]string{})
	log.Printf("There are %d people who have already voted\n", len(alreadyVoted))

	// load greek members
	log.Println("Loading members...")
	members := loadMembers("data/Greek-Chapter-Rosters.csv")

	var orgs map[string]Organization = make(map[string]Organization)

	for _, member := range members {
		//create org if first time seeing it
		if _, ok := orgs[member.Organization]; !ok {
			orgs[member.Organization] = Organization{0, 0}
		}
		org := orgs[member.Organization]

		//see if member already has voted
		if util.Contains(&alreadyVoted, member.ONID) {
			org.Voted++
		}

		//always increment total
		org.Total++
		orgs[member.Organization] = org
	}

	log.Printf("%d organizations\n", len(orgs))

	//print the map
	var participationOutput []string
	for org, orgData := range orgs {
		participationOutput = append(participationOutput, fmt.Sprintf("%s,%d,%d, =%d/%d", org, orgData.Voted, orgData.Total, orgData.Voted, orgData.Total))
	}

	sort.Strings(participationOutput)

	//write to file
	util.StoreStringArrayFile(participationOutput, "greek-info.csv")

}

func loadMembers(file string) []Member {
	reg, err := regexp.Compile(`\b(Alpha|Beta|Gamma|Delta|Epsilon|Zeta|Eta|Theta|Iota|Kappa|Lambda|Mu|Nu|Xi|Omicron|Pi|Rho|Sigma|Tau|Upsilon|Phi|Chi|Psi|Omega|Fiji|Farmhouse|Acacia|AEPi|DTD|KDChi|SDO|Sig)`)
	if err != nil {
		log.Fatal(err)
	}

	var members []Member

	//load csv file
	f, err := os.Open(file)
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

		org := strings.TrimSpace(rec[1])

		//remove all characters after Sorority or Fraternity
		org = trimStringFromWord(org, "Sorority")
		org = trimStringFromWord(org, "Fraternity")
		org = trimStringFromWord(org, "Philanthropy")

		//get all the stuff without greek alphabet
		orgNoName := strings.TrimSpace(reg.ReplaceAllString(org, ""))

		//replace all the non-greek stuff at the end with nothing and keep the greek stuff
		finalOrg := strings.TrimSpace(strings.Replace(org, orgNoName, "", -1))

		//This indicates that I parsed something wrong
		if finalOrg == "" {
			log.Println("No organization found for:", org)
			continue
		}

		if finalOrg == "Pi" {
			log.Println("Pi found for:", org)
		}

		members = append(members, Member{ONID: rec[2], Organization: finalOrg})
	}

	return members
}

func trimStringFromWord(s, word string) string {
	if idx := strings.Index(s, word); idx != -1 {
		return s[:idx]
	}
	return s
}
