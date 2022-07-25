# ndca_to_json_go

A Go port of a python OCR parser.  Accepts a specialized OCR text format, then outputs a json file.

#### run the compiled program:

   - Download the latest release from https://github.com/GarrettArm/ndca_to_json_go/releases
   - Place the executable file into some folder, then create a new folder "source-data" within that folder.
   - Place the file "ndca_2007_08_tesseract_full_vol_read.txt" into that "source-data" folder.
   - Run the "ncda-to-json" program you downloaded in step 1.
   - See the newly created file "output.json" in the folder.  (Firefox is a good json reader)

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
