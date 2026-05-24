package internal

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// asciiStart is the ASCII value
	// of the first printable character (' ').
	asciiStart = 32

	// charHeight is the number of visual rows
	// used to render each ASCII-art character.
	charHeight = 8

	// linesPerChar represents how many lines
	// each character occupies in the banner file:
	// 8 visual rows + 1 separator line.
	linesPerChar = 9

	// expectedNewlines is the total number
	// of newline characters expected
	// inside a valid banner file.
	expectedNewlines = 855
)

// PrintAscii renders text as ASCII art
// using the selected banner file.
//
// Output is written to writer,
// which may be:
// - the terminal
// - a file
// - a test buffer
func PrintAscii(
	writer io.Writer,
	lines []string,
	filename string,
) error {
	// Read the entire banner file into memory.
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf(
			"could not open banner file: %w",
			err,
		)
	}

	// Normalize Windows line endings (\r\n)
	// into Unix format (\n).
	normalized := strings.ReplaceAll(
		string(data),
		"\r\n",
		"\n",
	)

	// Validate banner structure before rendering.
	if strings.Count(normalized, "\n") != expectedNewlines {
		return fmt.Errorf(
			"banner file is corrupt or invalid",
		)
	}

	// Split the banner into individual lines.
	bannerLines := strings.Split(normalized, "\n")

	// Process each logical input line separately.
	for i, line := range lines {

		// Empty logical lines represent explicit "\n".
		if line == "" {
			if i > 0 {
				fmt.Fprintln(writer)
			}
			continue
		}

		// Render all 8 visual rows
		// for the current logical line.
		for row := 1; row <= charHeight; row++ {
			var sb strings.Builder

			// Process every character in the line.
			for _, r := range line {

				// Calculate the corresponding line index
				// inside the banner file.
				index := (int(r)-asciiStart)*linesPerChar + row

				// Prevent invalid banner indexing.
				if index >= len(bannerLines) {
					return fmt.Errorf(
						"character %q is out of supported range in banner",
						r,
					)
				}

				// Append the correct ASCII-art segment
				// for the current visual row.
				sb.WriteString(bannerLines[index])
			}

			// Print the completed visual row.
			fmt.Fprintln(writer, sb.String())
		}
	}

	return nil
}
