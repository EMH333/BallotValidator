package main

// does the array contain the value?
func contains(s *[]string, e string) bool {
	for _, a := range *s {
		if a == e {
			return true
		}
	}
	return false
}
