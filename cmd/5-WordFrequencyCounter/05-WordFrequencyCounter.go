package main

import (
	"regexp"
	"strings"
)

// CountWordFrequency takes a string containing multiple words and returns
// a map where each key is a word and the value is the number of times that
// word appears in the string. The comparison is case-insensitive.
//
// Words are defined as sequences of letters and digits.
// All words are converted to lowercase before counting.
// All punctuation, spaces, and other non-alphanumeric characters are ignored.
//
// For example:
// Input: "The quick brown fox jumps over the lazy dog."
// Output: map[string]int{"the": 2, "quick": 1, "brown": 1, "fox": 1, "jumps": 1, "over": 1, "lazy": 1, "dog": 1}
func CountWordFrequency(text string) map[string]int {
	ret := make(map[string]int)

	words := strings.Fields(strings.ToLower(text))
	// fmt.Println("words = ", words)

	nonAlnumRegex := regexp.MustCompile(`[^a-z0-9-]`)
	for _, part := range words {
		var signalWord string
		var doubleWord []string
		signalWord = nonAlnumRegex.ReplaceAllString(part, "")
		// fmt.Println("signalWord = ", signalWord)
		if strings.Contains(signalWord, "-") {
			doubleWord = append(doubleWord, strings.Split(signalWord, "-")...)
			// fmt.Println("signalWord = ", doubleWord)
			for _, v := range doubleWord {
				ret[v]++
			}
			continue
		}
		ret[signalWord]++

	}
	return ret
}

func main() {
	text1 := "Go, go, go! Let's learn Go programming."
	text2 := "  Spaces,   tabs,\t\tand\nnew-lines are ignored!  "
	CountWordFrequency(text1)
	CountWordFrequency(text2)
}
