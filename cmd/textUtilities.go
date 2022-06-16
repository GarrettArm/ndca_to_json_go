package main

import (
	"strings"
	"unicode"
)

func isAllUpper(s string) bool {
	// identifies lines that are all uppercase
	// assumed to be college names & the start of a college block
	if len(s) == 0 {
		return false
	}
	// fail a string with any lowercase letters
	if strings.ToUpper(s) != s {
		return false
	}
	// fail a string that has a number in it
	for _, v := range s {
		if unicode.IsDigit(v) {
			return false
		}
	}
	// fail a string with any character that's neither a space, a letter, or a puctuation mark
	for _, v := range s {
		if !(unicode.IsLetter(v) || unicode.IsSpace(v) || unicode.IsPunct(v)) {
			return false
		}
	}
	// fail a string that doesn't have at least one letter
	hasLetter := false
	for _, v := range s {
		if unicode.IsLetter(v) {
			hasLetter = true
		}
	}
	if !hasLetter {
		return false
	}

	return true
}

func isAllNumeric(s string) bool {
	// identifies lines that are all numeric
	// assumed to be page numbers in OCR text
	for _, v := range s {
		if !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}
