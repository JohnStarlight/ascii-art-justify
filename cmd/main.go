package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"ascii-art/internal"
)

// fail prints a formatted error message
// and terminates the program with exit status code 1.
func fail(err error) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(1)
}

func main() {
	// Parse and validate command-line arguments.
	// Returns a Config struct containing:
	// - the input text
	// - the selected banner
	// - the optional output filename
	config, err := internal.ParseArgs(os.Args)
	if err != nil {
		fail(err)
	}

	// By default, all rendered ASCII art
	// is written to the terminal.
	var writer io.Writer = os.Stdout

	// If the user provided an output file,
	// create/truncate the file and redirect
	// all output writes to it instead.
	if config.OutputFile != "" {
		file, err := internal.PrepareOutputFile(config.OutputFile)
		if err != nil {
			fail(err)
		}

		// Replace terminal output with file output.
		writer = file

		// Ensure the file is properly closed
		// when main finishes execution.
		defer file.Close()
	}

	// Convert literal "\n" sequences into
	// separate logical text lines.
	//
	// Example:
	// "Hello\nWorld"
	// becomes:
	// []string{"Hello", "World"}
	lines := strings.Split(config.Text, "\\n")

	// Render the ASCII art using the selected banner.
	if err := internal.PrintAscii(writer, lines, config.BannerPath); err != nil {
		fail(err)
	}
}