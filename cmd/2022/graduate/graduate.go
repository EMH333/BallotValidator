package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"sort"
	"strings"

	"ethohampton.com/BallotValidator/internal/steps"
	"ethohampton.com/BallotValidator/internal/util"
)

//determine percentage of each graduate major that has voted
// load votes
// load graduate student data
// for each member, determine if they have voted
// if they have voted then add them to the appropriate major
// regardless, update total members for that major

//ran this to find deps with multiple majors
/* document.getElementsByClassName('sitemap')[0].childNodes[0].childNodes.forEach((i)=>{
  i.childNodes[1].childNodes.forEach((dep)=>{
    let majors = dep.childNodes[1]
    if(majors == undefined) {return;}
    let numGrad = 0;
    majors.childNodes.forEach((majorItems) => {
      let majorName = majorItems.firstChild.innerText;
      if(majorName !== undefined && (majorName.includes("MS") || majorName.includes("PhD"))){
        numGrad++;
      }
    });
    if(numGrad > 1){
      dep.style.backgroundColor = "green";
    }
  })
})
*/

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
	if os.IsNotExist(err) && os.Mkdir("output", 0o755) != nil {
		log.Fatal("Could not create output directory")
	}

	var countingConfig util.CountingConfig

	// Load the votes
	log.Println("Loading votes...")
	votes := util.LoadVotesCSV(&countingConfig, "data/ballots/"+dataFile, START, END)
	log.Printf("%d votes loaded\n", len(votes))

	// reuse step two to get the ONID emails for all eligible students that have already voted
	_, _, alreadyVoted, _ := steps.StepTwo(votes, &[]string{})
	log.Printf("There are %d people who have already voted\n", len(alreadyVoted))

	// load graduate students
	log.Println("Loading members...")
	students := loadStudents("data/graduate-students.csv")

	var orgs = make(map[string]Organization)

	for _, member := range students {
		//create org if first time seeing it
		if _, ok := orgs[member.Organization]; !ok {
			orgs[member.Organization] = Organization{0, 0}
		}
		org := orgs[member.Organization]

		//see if member already has voted
		if slices.Contains(alreadyVoted, member.ONID) {
			org.Voted++
		}

		//always increment total
		org.Total++
		orgs[member.Organization] = org
	}

	log.Printf("%d degrees\n", len(orgs))

	//print the map
	var output []string
	for org, orgData := range orgs {
		output = append(output, fmt.Sprintf("%s,%d,%d, =%d/%d", org, orgData.Voted, orgData.Total, orgData.Voted, orgData.Total))
	}

	sort.Strings(output)

	//write to file
	util.StoreStringArrayFile(output, "graduate-info.csv", true)

}

func loadStudents(file string) []Member {
	var members []Member

	//load csv file
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // nolint:errcheck // don't care about close

	// read csv values using csv.Reader
	//with modifications to handle the specifics of the valid votes list
	csvReader := csv.NewReader(f)
	csvReader.Comma = ','
	csvReader.TrimLeadingSpace = true

	_, err = csvReader.Read() //skip header
	if err != nil {
		log.Fatal(err)
	}

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		spitName := strings.Split(rec[4], "-")
		var org string
		if len(spitName) == 2 {
			org = strings.TrimSpace(spitName[1])
		} else {
			org = strings.TrimSpace(spitName[0])
		}

		org = strings.ReplaceAll(org, ",", "") //remove commas

		if org == "Kinesiology" || org == "Nutrition" {
			org = "Bio & Pop Health Sciences"
		}

		if org == "Applied Anthropology" ||
			org == "College Student Services Admin" ||
			org == "Women Gender and Sexuality" {
			org = "Lang Culture & Society"
		}

		if org == "Applied Ethics" ||
			org == "History" ||
			org == "History of Science" {
			org = "History Philosophy & Religion"
		}

		if org == "Nuclear Engineering" || org == "Radiation Health Physics" {
			org = "Nuclear Science & Engineering"
		}

		if org == "Industrial Engineering" ||
			org == "Materials Science" ||
			org == "Mechanical Engineering" ||
			org == "Robotics" {
			org = "MIME"
		}

		if org == "Computer Science" ||
			org == "Electrical and Computer Engr" ||
			org == "Artificial Intelligence" {
			org = "EECS"
		}

		if org == "Bioengineering" ||
			org == "Chemical Engineering" ||
			org == "Environmental Engineering" {
			org = "Chem & Bio & Enviromental Engr"
		}

		if org == "Fisheries Science" ||
			org == "Wildlife Science" {
			org = "Fisheries & Wildlife & Convervation Sciences"
		}

		if org == "Crop Science" ||
			org == "Soil Science" {
			org = "Crop and Soil Science"
		}

		if org == "Animal Science" ||
			org == "Rangeland Ecology & Management" {
			org = "Animal & Rangeland Science"
		}

		if org == "Creative Writing" {
			org = "English"
		}

		//This indicates that I parsed something wrong
		if org == "" {
			log.Println("No organization found for:", org)
			continue
		}

		members = append(members, Member{ONID: rec[2], Organization: org})
	}

	return members
}
