package main

import (
	"flag"
	"log"
	"regexp"
)

var rComment = regexp.MustCompile(`//\s*@inject_field:\s+(\S+)\s+(\S+)$`)

// customField is the information of the custom customField to inject to struct.
type customField struct {
	fieldName string
	fieldType string
}

// textArea is information of the struct to inject custom fields to.
type textArea struct {
	name  string
	start int
	end   int
	// where to insert custom fields, should be after end of last field.
	// We have to record it because between the closing brace and the end
	// of the struct can contains comment, white spaces...
	insertPos int
	fields    []*customField
}

func main() {
	var inputFile string

	flag.StringVar(&inputFile, "input", "", "path to input file")

	flag.Parse()

	if len(inputFile) == 0 {
		log.Fatal("input file is mandatory")
	}

	areas, err := parseFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	if err = writeFile(inputFile, areas); err != nil {
		log.Fatal(err)
	}
}
