package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"ascii-art/internal"
)

// fail prints the error and exits the program with status code 1.
func fail(err error) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(1)
}

func main() {
	config, err := internal.ParseArgs(os.Args)
	if err != nil {
		fail(err)
	}

	// By default, render output to the terminal.
	var writer io.Writer = os.Stdout

	// If an output file was requested, create it and redirect writes there.
	if config.OutputFile != "" {
		file, err := internal.PrepareOutputFile(config.OutputFile)
		if err != nil {
			fail(err)
		}

		writer = file
		defer file.Close()
	}

	// Split literal "\n" sequences into separate logical lines.
	lines := strings.Split(config.Text, "\\n")

	if err := internal.PrintAscii(writer, lines, config.BannerPath); err != nil {
		fail(err)
	}
}