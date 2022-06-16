package main

import (
	"bufio"
	"os"
	"strings"

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
