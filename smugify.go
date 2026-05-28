package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"github.com/The0xGeek/smugify/utils"
)

func main() {
	// cli flags ==========
	trustFile := flag.String("trust", "", "Path to the SVG template file (used for context)")
	attachFile := flag.String("attach", "", "Path to the file to be attached")
	outputName := flag.String("out", "", "Custom filename for the downloaded attachment")

	flag.Parse()

	// validation ==========
	if *attachFile == "" {
		exitWithError("missing required flag: -attach")
	}

	if *outputName == "" {
		*outputName = *attachFile
	}

	// read files ==========
	data, err := os.ReadFile(*attachFile)
	if err != nil {
		exitWithError(fmt.Sprintf("%v", err))
	}

	// basic size guard ==========
	if len(data) == 0 {
		exitWithError(fmt.Sprintf("open %v: attach file is empty", *attachFile))
	}

	// base64 encode ==========
	base64String := base64.StdEncoding.EncodeToString(data)

	// generate svg ==========
	var outSVG string

	if *trustFile == "" {
		outSVG, err = utils.GenerateSVG(base64String, *outputName)
	} else {
		outSVG, err = utils.GenerateTrustedSVG(*trustFile, base64String, *outputName)
	}

	if err != nil {
		exitWithError(fmt.Sprintf("%v", err))
	}

	// write output ==========
	outputFile := "output.svg"

	err = os.WriteFile(outputFile, []byte(outSVG), 0644)
	if err != nil {
		exitWithError(fmt.Sprintf("Failed to write output file: %v", err))
	}

	// success ==========
	fmt.Printf("Done! Output written to %s\n", outputFile)
}

// centeralize error handling ==========
func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}