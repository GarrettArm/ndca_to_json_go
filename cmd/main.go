package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func main() {
	// builds a book item whose colleges items are printed to a json file
	book := book{}
	book.loadFile("source-data/ndca_2007_08_tesseract_full_vol_read.txt")
	book.findBodyBoundaries()
	book.getBody()
	book.addColleges()
	book.addCollegeEnds()

	for n := range book.colleges {
		book.colleges[n].addText(book.body)
		book.colleges[n].addSingleLiners()
		book.colleges[n].addMultiLiners()
	}

	toJson, err := json.Marshal(book.colleges)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("output.json", toJson, os.ModePerm)
}
