package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"ascii-art/internal"
)

// run contains the main application logic.
func run() error {
	config, err := internal.ParseArgs(os.Args)
	if err != nil {
		return err
	}

	var writer io.Writer = os.Stdout

	if config.OutputFile != "" {
		file, err := internal.PrepareOutputFile(config.OutputFile)
		if err != nil {
			return err
		}

		// Ensure the file is closed before run exits.
		defer file.Close()

		writer = file
	}

	// Split the input text into separate logical lines.
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

// Any error is returned to main for centralized handling.
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		os.Exit(1)
	}
}
