/*
Task : Palindrome Check
Write a Go function that takes a string as input and checks whether it is a palindrome or not. A palindrome is a word, phrase, number, or other sequence of characters that reads the same forward and backward (ignoring spaces, punctuation, and capitalization).
[Optional]: Write test for your function
*/
package main

import (
	"strings"
	"unicode"
)

func CheckPalindrome(str string) bool {
	lowerStr := strings.ToLower(str)
	left, right := 0, len(str)-1
	for left < right {

		chL := rune(lowerStr[left])
		for !unicode.IsLetter(chL) {
			left += 1
			chL = rune(lowerStr[left])
		}

		chR := rune(lowerStr[right])
		for !unicode.IsLetter(chR) {
			right -= 1
			chR = rune(lowerStr[right])
		}
		// fmt.Println(string(str[left]), string(str[right]))
		if str[left] != str[right] {
			return false
		}

		left += 1
		right -= 1
	}

	return true
}
