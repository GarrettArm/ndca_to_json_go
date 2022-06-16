# ndca_to_json_go

A Go port of a python OCR parser.  Accepts a specialized OCR text format, then outputs a json file.

#### run the script without compiling:

install go on your computer, then:

`go run ./cmd/*.go`

#### build an executable:

For your computer:

`go build -o ${filename.extension} ./cmd/*.go`

Windows:

`GOOS=windows GOARCH=amd64 go build -o ndca-to-json.exe ./cmd/*.go`

Linux:

`GOOS=linux GOARCH=amd64 go build -o ndca-to-json.exe ./cmd/*.go`

M1 apple:

`GOOS=darwin GOARCH=arm64 go build -o ndca-to-json.exe ./cmd/*.go`
