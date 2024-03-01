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

	// replace write-in entries with the real candidate
	vote = strings.ReplaceAll(vote, "ADISON ROWE", "Adison Rowe")
	vote = strings.ReplaceAll(vote, "ALEXA YESSENIA GOMEZ SILVA", "Alexa Yessenia Gomez Silva")
	vote = strings.ReplaceAll(vote, "CAMRYN LAU", "Camryn Lau")
	vote = strings.ReplaceAll(vote, "CARTER HAWES", "Carter Hawes")
	vote = strings.ReplaceAll(vote, "CHRISTINA LEWANDOWSKI", "Christina Lewandowski")
	vote = strings.ReplaceAll(vote, "COLE PETERS", "Cole Peters")
	vote = strings.ReplaceAll(vote, "DESTINY JONES", "Destiny Jones")
	vote = strings.ReplaceAll(vote, "DONOVAN GIOVANNI MORALES-COONRAD", "Donovan Giovanni Morales-Coonrad")
	vote = strings.ReplaceAll(vote, "DYLAN PERFECT", "Dylan Perfect")
	vote = strings.ReplaceAll(vote, "ELIZABETH ECKMAN", "Elizabeth Eckman")
	vote = strings.ReplaceAll(vote, "EMERSON PEARSON", "Emerson Pearson")
	vote = strings.ReplaceAll(vote, "EMILY LILAK", "Emily Lilak")
	vote = strings.ReplaceAll(vote, "HENRIETTA RUTAREMWA", "Henrietta Rutaremwa")
	vote = strings.ReplaceAll(vote, "HERAN GAO", "HeRan Gao")
	vote = strings.ReplaceAll(vote, "JACOB FIELD", "Jacob Field")
	vote = strings.ReplaceAll(vote, "JACOB RICHARD RUTHARDT", "Jacob Richard Ruthardt")
	vote = strings.ReplaceAll(vote, "JAKE HUSH", "Jake Hush")
	vote = strings.ReplaceAll(vote, "JULIAN LOESCH", "Julian Loesch")
	vote = strings.ReplaceAll(vote, "KATYAYANI (KATYA) KARLAPATI", "Katyayani (Katya) Karlapati")
	vote = strings.ReplaceAll(vote, "KIERAN HOSTETLER-MCLAUGHLIN", "Kieran Hostetler-McLaughlin")
	vote = strings.ReplaceAll(vote, "KYLE LOCKE", "Kyle Locke")
	vote = strings.ReplaceAll(vote, "LAUREN CAMOU", "Lauren Camou")
	vote = strings.ReplaceAll(vote, "LEAH ANN WRIGHT", "Leah Ann Wright")
	vote = strings.ReplaceAll(vote, "LILLIAN JUDITH GOODYEAR", "Lillian Judith Goodyear")
	vote = strings.ReplaceAll(vote, "MADISON WUSSTIG", "Madison Wusstig")
	vote = strings.ReplaceAll(vote, "MARISA IKEHARA", "Marisa Ikehara")
	vote = strings.ReplaceAll(vote, "MERCEDEZ ALLEN", "Mercedez Allen")
	vote = strings.ReplaceAll(vote, "MORGAN WOODRICH", "Morgan Woodrich")
	vote = strings.ReplaceAll(vote, "NATHAN RASCHKES", "Nathan Raschkes")
	vote = strings.ReplaceAll(vote, "REBECCA J LANG", "Rebecca J Lang")
	vote = strings.ReplaceAll(vote, "ROMAN LEWIS", "Roman Lewis")
	vote = strings.ReplaceAll(vote, "SAEGIS ABBOTT", "Saegis Abbott")
	vote = strings.ReplaceAll(vote, "SAM JEWETT", "Sam Jewett")
	vote = strings.ReplaceAll(vote, "SHAWN AUNDRAE DURR", "Shawn Aundrae Durr")
	vote = strings.ReplaceAll(vote, "SIANNA STONE", "Sianna Stone")
	vote = strings.ReplaceAll(vote, "SOPHIA NOWERS", "Sophia Nowers")
	vote = strings.ReplaceAll(vote, "SPENCER THOMAS KOWASH", "Spencer Thomas Kowash")
	vote = strings.ReplaceAll(vote, "TEAGHAN KNOX", "Teaghan Knox")

	return vote
}
