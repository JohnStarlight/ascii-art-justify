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

		starts := findAll(line, config.Part) // NEW: color substring now comes from config
		partLen := len(config.Part)

		if config.Color != "" && config.Part == "" {
			partLen = len(line)
		}

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

				if config.Color != "" && inColorRange(pos, starts, partLen) {
					sb.WriteString(config.Color + segment + "\033[0m") // NEW: color now comes from config
				} else {
					sb.WriteString(segment)
				}
			}

			rows = append(rows, sb.String()) // NEW: keep the rendered row instead of printing immediately.
		}

		rows = alignRows(rows, config.Align, TerminalWidth()) // NEW: applies left/center/right alignment before printing.

		for _, renderedRow := range rows {
			fmt.Fprintln(writer, renderedRow)
		}
	}

	return nil
}

func findAll(line, part string) []int {
	if part == "" {
		return []int{0}
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

func inColorRange(pos int, starts []int, partLen int) bool {
	for _, start := range starts {
		if pos >= start && pos < start+partLen {
			return true
		}
	}

	return false
}
