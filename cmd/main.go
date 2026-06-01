package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"ascii-art/internal"
)

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

		defer file.Close()

		writer = file
	}

	// Split on the literal two-character sequence \n as typed on the CLI.
	lines := strings.Split(config.Text, "\\n")

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
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		os.Exit(1)
	}
}
