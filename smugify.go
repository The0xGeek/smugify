package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"github.com/The0xGeek/smugify/utils"
)

func main() {
	// GET USER INPUTS FROM TERMINAL ==========
	trustFile := flag.String("trust", "", "Path to the SVG template file (used for context)")
	attachFile := flag.String("attach", "", "Path to the file to be attached")
	outputName := flag.String("out", "", "Custom filename for the downloaded attachment")

	flag.Parse()

	// CHECK USER INPUTS ==========
	if *attachFile == "" {
		fmt.Fprintf(os.Stderr, "Error: The '-attach' flag is required but was not provided.\n")
		os.Exit(1)
	}

	if *outputName == "" {
		*outputName = *attachFile
	}

	// OPEN THE FILE ==========
	data, err := os.ReadFile(*attachFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// CONVERT ATTACHMENT TO BASE64 ==========
	base64String := base64.StdEncoding.EncodeToString(data)

	// GENERATE FINAL SVG FILE ==========
	var outSVG string
	if *trustFile == "" {
		outSVG = utils.GenerateSVG(base64String, *outputName)
	} else {
		outSVG = utils.GenerateTrustedSVG(*trustFile, base64String, *outputName)
	}

	// WRITE OUTPUT FILE ==========
	err = os.WriteFile("output.svg", []byte(outSVG), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// SUCCESS MESSAGE ==========
	fmt.Println("Done!")
}
