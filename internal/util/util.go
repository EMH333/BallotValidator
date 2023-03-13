package util

import "strings"

// does the array contain the value?
func Contains(s *[]string, e string) bool {
	for _, a := range *s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func CleanVote(vote string) string {
	vote = strings.TrimSpace(vote)

	if vote == "Write in:" || vote == "Write-in:" || vote == "Write-In" {
		return ""
	}

	vote = strings.ToUpper(vote)

	//TODO add removal of entrys found to be invalid
	//TODO add subsitution for valid but ill-formed entrys
	//both should happen after the trim/upper

	// fix spelling error for real person
	vote = strings.ReplaceAll(vote, "CAROYLN PEARCE", "CAROLYN PEARCE")

	// actual candidates for sfc at large, so dedupe those votes
	vote = strings.ReplaceAll(vote, "SOPHIA NOWERS", "Sophia Nowers")
	vote = strings.ReplaceAll(vote, "ABUKAR MOHAMMED", "Abukar Mohamed")
	vote = strings.ReplaceAll(vote, "COLE PETERS", "Cole Peters")

	//for senate, same person
	vote = strings.ReplaceAll(vote, "DOLA POPOOLA", "Ayodola Kayode-Popoola")
	//for senate
	vote = strings.ReplaceAll(vote, "ABHISHEK ENAGUTHI", "Abhishek Enaguthi")
	vote = strings.ReplaceAll(vote, "ALEXANDER KERNER", "Alexander Kerner")
	vote = strings.ReplaceAll(vote, "ANTHONY BUTLER-TORREZ", "Anthony Butler-Torrez")
	vote = strings.ReplaceAll(vote, "AUDREY PORTER", "Audrey Porter")
	vote = strings.ReplaceAll(vote, "AYODOLA KAYODE-POPOOLA", "Ayodola Kayode-Popoola")
	vote = strings.ReplaceAll(vote, "CAMRYN LAU", "Camryn Lau")
	vote = strings.ReplaceAll(vote, "CAROLYN PEARCE", "Carolyn Pearce")
	vote = strings.ReplaceAll(vote, "CONNOR ROBERTS", "Connor Roberts")
	vote = strings.ReplaceAll(vote, "ELIZABETH ECKMAN", "Elizabeth Eckman")
	vote = strings.ReplaceAll(vote, "ERICA NYARKO", "Erica Nyarko")
	vote = strings.ReplaceAll(vote, "EVAN RUDISILE", "Evan Rudisile")
	vote = strings.ReplaceAll(vote, "GABRIEL THOMISON", "Gabriel Thomison")
	vote = strings.ReplaceAll(vote, "HAILEY BROWN", "Hailey Brown")
	vote = strings.ReplaceAll(vote, "JAMIE HAMLIN", "Jamie Hamlin")
	vote = strings.ReplaceAll(vote, "KATYAYANI KARLAPATI", "Katyayani Karlapati")
	vote = strings.ReplaceAll(vote, "MAIA BROWN", "Maia Brown")
	vote = strings.ReplaceAll(vote, "MARCUS PAUL ANTOLLI", "Marcus Paul Antolli")
	vote = strings.ReplaceAll(vote, "NATHAN SCHMIDT", "Nathan Schmidt")
	vote = strings.ReplaceAll(vote, "OLIVIA CARTWRIGHT", "Olivia Cartwright")
	vote = strings.ReplaceAll(vote, "RHYAN STEPHENSON", "Rhyan Stephenson")
	vote = strings.ReplaceAll(vote, "RICHARD DAVID DEININGER", "Richard David Deininger")
	vote = strings.ReplaceAll(vote, "RICHARD PIAZZA", "Richard Piazza")
	vote = strings.ReplaceAll(vote, "RILEY WALSH", "Riley Walsh")
	vote = strings.ReplaceAll(vote, "SARAH THEALL", "Sarah Theall")
	vote = strings.ReplaceAll(vote, "SHAWN DURR", "Shawn Durr")
	vote = strings.ReplaceAll(vote, "TEAGHAN KNOX", "Teaghan Knox")

	return vote
}
