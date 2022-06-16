package main

import (
	"bufio"
	"os"
	"strings"
	"unicode"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type book struct {
	fulltext       []string
	bodyBoundaries []int
	body           []string
	colleges       []college
}

func (book *book) loadFile(filename string) {
	// reads a windows 1252 encoded file
	// adds the file's fulltext to the book structure
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decodingReader := transform.NewReader(file, charmap.Windows1252.NewDecoder())
	lines := []string{}

	scanner := bufio.NewScanner(decodingReader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	book.fulltext = lines
}

func (book *book) findBodyBoundaries() {
	// we're expecting the crucial part of the book to be demarcated
	// by starting/ending lines with text "SENIOR COLLEGES"
	lineNums := []int{}
	for n, v := range book.fulltext {
		if strings.Trim(v, " ") == "SENIOR COLLEGES" {
			lineNums = append(lineNums, n)
		}
	}
	book.bodyBoundaries = lineNums
}

func (book *book) getBody() {
	// gets the text lines between the boundary line numbers
	start := book.bodyBoundaries[0]
	end := book.bodyBoundaries[1]
	body := book.fulltext[start:end]
	book.body = body
}

func (book *book) addColleges() {
	// creates a college struct for each parsable college in the ocr's body
	// a parsable college is identified by all uppercase letters
	// fills out their 'Name' and 'startLineNum'
	// leaves other fields as null values for later function to fill in
	for lineNum, lineText := range book.body {
		trimmed := strings.TrimFunc(lineText, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsPunct(r)
		})
		if isAllUpper(trimmed) {
			c := college{
				Name:         trimmed,
				startLineNum: lineNum,
				SingleLiners: make(map[string]string),
				MultiLiners:  make(map[string]map[string]string),
				UnusedText:   []string{},
			}
			book.colleges = append(book.colleges, c)
		}
	}
}

func (book *book) addCollegeEnds() {
	// fills in college.endLineNum, using the next item in the colleges list
	for n, college := range book.colleges {
		if n == 0 {
			continue
		}
		book.colleges[n-1].endLineNum = college.startLineNum - 1
	}
	book.colleges[len(book.colleges)-1].endLineNum = len(book.body)
}
