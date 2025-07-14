/*

Task:  Word Frequency Count
Write a Go function that takes a string as input and returns a dictionary containing the frequency of each word in the string. Treat words in a case-insensitive manner and ignore punctuation marks.
[Optional]: Write test for your function

*/

package main

import (
	"strings"
	"unicode"
)

func WordCounter(words string) map[string]int {
	lowerWords := strings.ToLower(words)
	alphabet := make(map[string]int)
	for _, ch := range lowerWords {
		if unicode.IsLetter(ch) {
			alphabet[string(ch)] += 1
		}
	}
	return alphabet
}
