package main

import (
	"fmt"
	"strings"
)

type college struct {
	Name         string                       `json:"Name"`
	SingleLiners map[string]string            `json:"SingleLiners"`
	MultiLiners  map[string]map[string]string `json:"MultiLiners"`
	startLineNum int
	endLineNum   int
	Text         []string `json:"Text"`
	UnusedText   []string `json:"UnusedText"`
}

func (c *college) addText(body []string) {
	// fills in the college.text field
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

func (c *college) addSingleLiners() {
	// converts some lines in college.text into key:value
	// SINGLE_LINERS is a matchlist of expected keys
	// if you split the text line on "-" and the first portion is in SINGLE_LINERS
	// then the college.SingleLiners[SINGLE_LINER] = second portion
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

func (c *college) addMultiLiners() {
	// converts some lines in college.text into key:role:value
	// MULTI_LINERS is a matchlist of expected keys
	// ROLES is a matchlist of expected roles
	// if you split the text line on "-" and the first portion is in MULTI_LINERS
	// then a college[MULTI_LINER] = {"Lead": value, }
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
		// sometimes clumps span across book pages,
		// so if the label is repeated then consider it a continuation of the clump
		if counter+1 >= len(sectionStarts) {
			clump[label] = c.Text[startLine:]
		} else {
			clump[label] = c.Text[startLine:sectionStarts[counter+1]]
		}
	}
	c.MultiLiners = c.parseSecondLayer(clump)
}

func (c *college) parseSecondLayer(clump map[string][]string) map[string]map[string]string {
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
				c.removeFromUnprocessed(line)
				continue
			}
			if found && rolesSet[pre] {
				sport[pre] = post
				c.removeFromUnprocessed(line)
				continue
			}
			if !found || !rolesSet[pre] {
				// There are unlabeled lines in the source text that probably indicate Asst.
				// so, we're making an "Asst." key if one doesn't exist, else concating to the existing one
				_, ok := sport["Asst."]
				if !ok {
					sport["Asst."] = line
					c.removeFromUnprocessed(line)
					continue
				} else {
					sport["Asst."] = fmt.Sprintf("%s %s", sport["Asst."], line)
					c.removeFromUnprocessed(line)
					continue
				}
			}
		}
		allSports[label] = sport
	}
	return allSports
}

func (c *college) removeFromUnprocessed(usedLine string) {
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
