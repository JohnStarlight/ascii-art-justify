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
	// each character occupies inside a banner file.
	//
	// 8 visual rows + 1 separator line.
	linesPerChar = 9

	// expectedNewlines is the total number
	// of newline characters expected inside
	// a valid banner file.
	expectedNewlines = 855
)

// PrintAscii renders text as ASCII art
// using the selected banner file.
//
// The rendered output is written to writer,
// which may be:
// - the terminal
// - a file
// - a buffer during testing
func PrintAscii(writer io.Writer, lines []string, filename string) error {
	// Read the entire banner file into memory.
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf(
			"could not open banner file: %w",
			err,
		)
	}

	// Validate the banner structure before rendering.
	//
	// A valid banner file must contain
	// exactly the expected number of newline characters.
	if strings.Count(string(data), "\n") != expectedNewlines {
		return fmt.Errorf(
			"banner file is corrupt or invalid",
		)
	}

	// Normalize Windows line endings (\r\n)
	// into Unix format (\n).
	//
	// This ensures compatibility across operating systems.
	normalized := strings.ReplaceAll(string(data), "\r\n", "\n")

	// Split the banner file into individual lines.
	bannerLines := strings.Split(normalized, "\n")

	// Process each logical input line separately.
	for i, line := range lines {

		// Empty logical lines represent explicit "\n".
		//
		// Example:
		// "Hello\n\nWorld"
		//
		// The middle empty string creates
		// one empty output line.
		if line == "" {
			if i > 0 {
				fmt.Fprintln(writer)
			}
			continue
		}

		// Render all 8 visual rows
		// for the current text line.
		for row := 1; row <= charHeight; row++ {

			// strings.Builder efficiently builds
			// the full row before printing it.
			var sb strings.Builder

			// Process every rune (character)
			// of the current input line.
			for _, r := range line {

				// Calculate the corresponding line index
				// inside the banner file.
				//
				// Formula:
				// (ASCII value - 32) * 9 + current row
				index := (int(r)-asciiStart)*linesPerChar + row

				// Prevent out-of-range banner access.
				if index >= len(bannerLines) {
					return fmt.Errorf(
						"character %q is out of supported range in banner",
						r,
					)
				}

				// Append the correct ASCII-art segment
				// for the current row of the character.
				sb.WriteString(bannerLines[index])
			}

			// Print the completed visual row.
			fmt.Fprintln(writer, sb.String())
		}
	}

	return nil
}