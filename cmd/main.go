package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"ascii-art/internal"
)

func main() {
	config, err := internal.ParseArgs(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var writer io.Writer
	writer = os.Stdout

	if config.OutputFile != "" {
		file, err := internal.PrepareOutputFile(config.OutputFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		writer = file
		defer file.Close()
	}

	lines := strings.Split(config.Text, "\\n")

	if err := internal.PrintAscii(writer, lines, config.BannerPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
