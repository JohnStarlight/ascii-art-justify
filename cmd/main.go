package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"ascii-art/internal"
)

// run contains the main application logic.
//
// It:
// - parses and validates CLI arguments
// - prepares the output destination
// - splits logical lines
// - renders ASCII art
//
// Any error is returned to main for centralized handling.
func run() error {
	// Parse and validate command-line arguments.
	config, err := internal.ParseArgs(os.Args)
	if err != nil {
		return err
	}

	// By default, write ASCII art to the terminal.
	var writer io.Writer = os.Stdout

	// If an output file was requested,
	// create/truncate the file and redirect output there.
	if config.OutputFile != "" {
		file, err := internal.PrepareOutputFile(config.OutputFile)
		if err != nil {
			return err
		}

		// Ensure the file is closed before run exits.
		defer file.Close()

		writer = file
	}

	// Convert literal "\n" sequences
	// into separate logical lines.
	lines := strings.Split(config.Text, "\\n")

	// Render ASCII art using the selected banner.
	if err := internal.PrintAscii(
		writer,
		lines,
		config.BannerPath,
	); err != nil {
		return err
	}

	return nil
}

func main() {
	// main is responsible only for:
	// - starting the application
	// - handling the final error
	// - returning the correct exit status
	if err := run(); err != nil {
		// Print errors to stderr instead of stdout.
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		// Exit with non-zero status code
		// to indicate program failure.
		os.Exit(1)
	}
}