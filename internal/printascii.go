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

// NEW: receives the whole Config so extensions like --align can be used
// without adding more and more parameters to PrintAscii.
func PrintAscii(
	writer io.Writer,
	config Config,
	lines []string,
) error {
	data, err := os.ReadFile(config.BannerPath) // NEW: banner path now comes from config
	if err != nil {
		return fmt.Errorf(
			"could not open banner file: %w",
			err,
		)
	}

	normalized := strings.ReplaceAll(string(data), "\r\n", "\n")

	if strings.Count(normalized, "\n") != expectedNewlines {
		return fmt.Errorf("banner file is corrupt or invalid")
	}

	bannerLines := strings.Split(normalized, "\n")

	for i, line := range lines {
		if line == "" {
			if i > 0 {
				fmt.Fprintln(writer)
			}
			continue
		}

		// NEW: colorAt maps every byte position of the line to the color
		// that applies there, so multiple --color flags can coexist.
		colorAt := colorIndexes(line, config)

		rows := make([]string, 0, charHeight) // NEW: stores the 8 rendered ASCII rows before printing, so we can align them.

		for row := 1; row <= charHeight; row++ {
			var sb strings.Builder

			for pos, r := range line {
				index := (int(r)-asciiStart)*linesPerChar + row

				if index >= len(bannerLines) {
					return fmt.Errorf(
						"character %q is out of supported range in banner",
						r,
					)
				}

				segment := bannerLines[index]

				if idx := colorAt[pos]; idx >= 0 {
					sb.WriteString(config.Colors[idx])
					sb.WriteString(segment)
					sb.WriteString("\033[0m")
				} else {
					sb.WriteString(segment)
				}
			}

			rows = append(rows, sb.String()) // NEW: keep the rendered row instead of printing immediately.
		}

		terminalWidth := TerminalWidth() // NEW: read terminal width once before applying alignment.

		if config.Align == "justify" {
			rows = justifyRows(line, bannerLines, terminalWidth, config) // NEW: justify needs the original text line to spread spaces between words.
		} else {
			rows = alignRows(rows, config.Align, terminalWidth) // NEW: left/center/right use the already-rendered rows.
		}

		for _, renderedRow := range rows {
			fmt.Fprintln(writer, renderedRow)
		}
	}

	return nil
}

// NEW: colorIndexes returns, for each byte position of line, the index into
// config.Colors that applies there, or -1 when the position is uncolored.
// The first --color flag whose substring covers a position wins, so overlaps
// between substrings resolve in flag order.
func colorIndexes(line string, config Config) []int {
	indexes := make([]int, len(line))
	for i := range indexes {
		indexes[i] = -1
	}

	for colorIdx, part := range config.Parts {
		if colorIdx >= len(config.Colors) {
			break
		}

		// An omitted substring colors the whole line.
		if part == "" {
			for i := range indexes {
				if indexes[i] == -1 {
					indexes[i] = colorIdx
				}
			}
			continue
		}

		for _, start := range findAll(line, part) {
			for i := start; i < start+len(part) && i < len(indexes); i++ {
				if indexes[i] == -1 {
					indexes[i] = colorIdx
				}
			}
		}
	}

	return indexes
}

func findAll(line, part string) []int {
	if part == "" {
		return nil
	}

	var starts []int
	offset := 0

	for {
		idx := strings.Index(line[offset:], part)
		if idx == -1 {
			break
		}

		starts = append(starts, offset+idx)
		offset += idx + len(part)
	}

	return starts
}
