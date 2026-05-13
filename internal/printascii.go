package internal

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	asciiStart       = 32
	charHeight       = 8
	linesPerChar     = 9
	expectedNewlines = 855
)

// PrintAscii renders input text as ASCII art using the selected banner file.
func PrintAscii(writer io.Writer, lines []string, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf(
			"could not open banner file: %w",
			err,
		)
	}

	// Validate banner structure before rendering.
	if strings.Count(string(data), "\n") != expectedNewlines {
		return fmt.Errorf(
			"banner file is corrupt or invalid",
		)
	}

	// Normalize Windows line endings (\r\n) into Unix format (\n).
	normalized := strings.ReplaceAll(string(data), "\r\n", "\n")
	bannerLines := strings.Split(normalized, "\n")

	for i, line := range lines {
		// Empty logical lines represent explicit "\n" separators.
		if line == "" {
			if i > 0 {
				fmt.Fprintln(writer)
			}
			continue
		}

		// Render all 8 visual rows for the current line.
		for row := 1; row <= charHeight; row++ {
			var sb strings.Builder

			for _, r := range line {
				index := (int(r)-asciiStart)*linesPerChar + row

				// Protect against invalid banner indexing.
				if index >= len(bannerLines) {
					return fmt.Errorf(
						"character %q is out of supported range in banner",
						r,
					)
				}

				sb.WriteString(bannerLines[index])
			}

			fmt.Fprintln(writer, sb.String())
		}
	}

	return nil
}