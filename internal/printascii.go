package internal

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	asciiStart = 32 // ASCII value of ' ', the first printable character.

	charHeight = 8

	// 8 visual rows + 1 blank separator line per character in the banner file.
	linesPerChar = 9

	// 95 printable ASCII characters (32–126) × 9 lines each.
	expectedNewlines = 855
)

func PrintAscii(
	writer io.Writer,
	lines []string,
	filename string,
) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf(
			"could not open banner file: %w",
			err,
		)
	}

	// Normalize Windows line endings so banner indexing is consistent.
	normalized := strings.ReplaceAll(
		string(data),
		"\r\n",
		"\n",
	)

	if strings.Count(normalized, "\n") != expectedNewlines {
		return fmt.Errorf(
			"banner file is corrupt or invalid",
		)
	}

	bannerLines := strings.Split(normalized, "\n")

	for i, line := range lines {

		// Empty logical lines represent explicit "\n" in the input.
		if line == "" {
			if i > 0 { // suppress a blank line before the first rendered block
				fmt.Fprintln(writer)
			}
			continue
		}

		// row starts at 1 to skip the blank separator at index 0 of each character block.
		for row := 1; row <= charHeight; row++ {
			var sb strings.Builder

			for _, r := range line {
				// Each character occupies linesPerChar lines in the banner file.
				index := (int(r)-asciiStart)*linesPerChar + row

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
