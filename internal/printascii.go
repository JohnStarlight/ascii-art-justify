package internal

import (
	"fmt"
	"os"
	"strings"
)

// PrintAscii renders each line of text as ASCII art using the given banner file.
// Each character maps to a 8-row block in the banner file, where blocks are
// separated by a blank line (9 lines per character total).
// Empty lines in the input produce a single blank line in the output,
// except for the first element which is skipped if empty.
func PrintAscii(lines []string, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not open banner file: %w", err)
	}

	// Normalize line endings to \n so banner files with \r\n (Windows)
	// are handled correctly alongside Unix-style files.
	bannerLines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	for i, line := range lines {
		// An empty string means the input contained a \n at this position.
		// We print one blank line as a separator, but skip the very first
		// element to avoid a leading blank line when the input itself is empty.
		if line == "" {
			if i > 0 {
				fmt.Println()
			}
			continue
		}

		// Print all 8 rows of ASCII art for this line of text.
		// Each character c occupies rows [(c-32)*9+1] to [(c-32)*9+8]
		// in the banner file (ASCII printable range starts at 32 = space).
		for row := 1; row <= 8; row++ {
			var sb strings.Builder
			for _, r := range line {
				index := (int(r)-32)*9 + row

				// Safety check to avoid out-of-bounds access
				if /* index < 0 || */ index >= len(bannerLines) {
					return fmt.Errorf("character %q is out of supported range in banner", r)
				}

				sb.WriteString(bannerLines[index])
			}
			fmt.Println(sb.String())
		}
	}

	return nil
}
