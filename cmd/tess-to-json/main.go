package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func readWindowFile(filename string) []string {
	// reads a windows 1252 encoded file
	// outputs the full text, split into lines
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
	return lines
}

func findSeniorCollegesLines(textLines []string) []int {
	l := []int{}
	for n, v := range textLines {
		if strings.Trim(v, " ") == "SENIOR COLLEGES" {
			l = append(l, n)
		}
	}
	return l
}

func getBody(textLines []string, seniorCollegesLines []int) []string {
	start := seniorCollegesLines[0]
	end := seniorCollegesLines[1]
	body := textLines[start:end]
	return body
}

type collegeBlock struct {
	CollegeName  string                       `json:"collegeName"`
	SingleLiners map[string]string            `json:"SingleLiners"`
	MultiLiners  map[string]map[string]string `json:"MultiLiners"`
	startLineNum int
	endLineNum   int
	Text         []string `json:"Text"`
	UnusedText   []string `json:"UnusedText"`
}

func getColleges(body []string) []collegeBlock {
	// creates a collegeBlock struct for each parsable college in the ocr's body
	// fills out their 'collegeName' and 'startLineNum'
	// leaves the 'endLineNum' and 'text' as null values for later function to fill in
	allColleges := []collegeBlock{}
	for lineNum, lineText := range body {
		trimmed := strings.TrimFunc(lineText, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsPunct(r)
		})
		if isAllUpper(trimmed) {
			c := collegeBlock{
				CollegeName:  trimmed,
				startLineNum: lineNum,
				SingleLiners: make(map[string]string),
				MultiLiners:  make(map[string]map[string]string),
				UnusedText:   []string{},
			}
			allColleges = append(allColleges, c)
		}
	}
	return allColleges
}

func addEnds(colleges []collegeBlock, body []string) {
	// fills in collegeBlock.endLineNum, using next item in the colleges list
	for n, college := range colleges {
		if n == 0 {
			continue
		}
		colleges[n-1].endLineNum = college.startLineNum - 1
	}
	colleges[len(colleges)-1].endLineNum = len(body)
}

func (c *collegeBlock) addText(body []string) {
	// fills in the collegeBlock.text field
	// reads the OCR body lines, then filters garbage lines
	var accumulate []string
	for _, line := range body[c.startLineNum:c.endLineNum] {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		if strings.Contains(trimmed, "NEW PAGE") {
			continue
		}
		if strings.Contains(trimmed, "www.collegiatedirectories.com") {
			continue
		}
		if isAllNumeric(trimmed) {
			continue
		}
		accumulate = append(accumulate, trimmed)
	}
	c.Text = accumulate
}

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

func (c *collegeBlock) addSingleLiners() {
	// converts some lines in collegeBlock.text into key:value
	// SINGLE_LINERS is a matchlist of expected keys
	// if you split the text line on "-" and the first portion is in SINGLE_LINERS
	// then the collegeBlock.SingleLiners[SINGLE_LINER] = second portion
	SINGLE_LINERS := []string{
		"Affiliation",
		"Conference",
		"Enrollment",
		"Colors",
		"Nickname",
		"Pres.",
		"Stadium",
		"Arena",
		"AD",
		"Acad. Adv.",
		"Acad. Affairs",
		"FB Secy.",
		"Secy.",
		"Ath. Communications",
		"Fac. Rep.",
		"PE Dir.",
		"Intra. Dir.",
		"Tkt. Mgr.",
		"SWA",
		"Asst. Aquatics Dir.",
		"Ath. Secy.",
		"Mgr. FB Ops",
		"Aquatics Dir.",
	}
	singleLinersSet := sliceToSet(SINGLE_LINERS)

	for _, textLine := range c.Text {
		// if textLine is splitable & pre is in SingleLinersSet
		pre, post, found := strings.Cut(textLine, "-")
		if found && singleLinersSet[pre] {
			c.SingleLiners[pre] = post
		} else {
			c.UnusedText = append(c.UnusedText, textLine)
		}
	}
}

func (c *collegeBlock) addMultiLiners() {
	// converts some lines in collegeBlock.text into key:role:value
	// MULTI_LINERS is a matchlist of expected keys
	// ROLES is a matchlist of expected roles
	// if you split the text line on "-" and the first portion is in MULTI_LINERS
	// then a collegeBlock[MULTI_LINER] = {"Lead": value, }
	// and the next lines are added to the map
	MULTI_LINERS := []string{
		"Football",
		"Basketball",
		"Baseball",
		"Cross Country",
		"Diving",
		"Golf",
		"Soccer",
		"Tennis",
		"Track",
		"Archery",
		"Aquatics",
		"Badminton",
		"Bowling",
		"Broomball",
		"Cheer",
		"Curling",
		"Cycling",
		"Equestrian",
		"Fencing",
		"Golf",
		"Gymnastics",
		"Handball",
		"Hockey",
		"Indoor soccer",
		"Lacrosse",
		"Rodeo",
		"Rugby",
		"Sailing",
		"Ski",
		"Squash",
		"Swim",
		"Strength",
		"Badminton",
		"Cross Country",
		"Crew",
	}
	multiLinersSet := sliceToSet(MULTI_LINERS)

	var sectionStarts []int
	for lineNum, line := range c.Text {
		pre, _, found := strings.Cut(line, "-")
		if found && multiLinersSet[pre] {
			sectionStarts = append(sectionStarts, lineNum)
		}
	}
	clump := make(map[string][]string)
	for counter, startLine := range sectionStarts {
		label, _, _ := strings.Cut(c.Text[startLine], "-")
		if counter+1 >= len(sectionStarts) {
			clump[label] = c.Text[startLine:]
		} else {
			clump[label] = c.Text[startLine:sectionStarts[counter+1]]
		}
	}
	c.MultiLiners = parseSecondLayer(clump, c)
}

func parseSecondLayer(clump map[string][]string, c *collegeBlock) map[string]map[string]string {
	ROLES := []string{
		"Asst.",
		"Assoc.",
		"Video Coord.",
		"Dir. Bask Ops.",
		"Bask. Secy.",
		"Dir. FB Ops.",
	}
	rolesSet := sliceToSet(ROLES)

	var allSports = make(map[string]map[string]string)

	for label, lines := range clump {
		sport := make(map[string]string)
		for lineNum, line := range lines {
			pre, post, found := strings.Cut(line, "-")
			if lineNum == 0 {
				sport["Lead"] = post
				removeFromUnprocessed(c, line)
				continue
			}
			if found && rolesSet[pre] {
				sport[pre] = post
				removeFromUnprocessed(c, line)
				continue
			}
			if !found || !rolesSet[pre] {
				// There are unlabeled lines in the source text that probably indicate Asst.
				// so, we're making an "Asst." key if one doesn't exist, else concating to the existing one
				_, ok := sport["Asst."]
				if !ok {
					sport["Asst."] = line
					removeFromUnprocessed(c, line)
					continue
				} else {
					sport["Asst."] = fmt.Sprintf("%s %s", sport["Asst."], line)
					removeFromUnprocessed(c, line)
					continue
				}
			}
		}
		allSports[label] = sport
	}
	return allSports
}

func removeFromUnprocessed(c *collegeBlock, usedLine string) {
	// remove one elem from the c.UsedText slice
	// return out of function immediately, to avoid removing two identical elements
	for count, line := range c.UnusedText {
		if count == len(c.UnusedText) {
			c.UnusedText = c.UnusedText[:count]
			return
		}
		if line == usedLine {
			c.UnusedText = append(c.UnusedText[:count], c.UnusedText[count+1:]...)
			return
		}
	}
}

func sliceToSet(slice []string) map[string]bool {
	set := make(map[string]bool)
	for _, v := range slice {
		set[v] = true
	}
	return set
}

func main() {
	textLines := readWindowFile("source-data/ndca_2007_08_tesseract_full_vol_read.txt")
	seniorCollegesLines := findSeniorCollegesLines(textLines)
	body := getBody(textLines, seniorCollegesLines)
	colleges := getColleges(body)
	addEnds(colleges, body)
	for n := range colleges {
		colleges[n].addText(body)
		colleges[n].addSingleLiners()
		colleges[n].addMultiLiners()
	}
	toJson, err := json.Marshal(colleges)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("output.json", toJson, os.ModePerm)
}
