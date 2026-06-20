package internal

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// NEW: TerminalWidth returns the current terminal width in columns.
// It is needed for --align because center/right/justify depend on
// how wide the terminal is.
func TerminalWidth() int {
	// NEW: first try `stty size`, using the real terminal input.
	// It usually gives the most accurate current terminal size.
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin

	output, err := cmd.Output()
	if err == nil {
		parts := strings.Fields(string(output))

		if len(parts) == 2 {
			width, convErr := strconv.Atoi(parts[1])
			if convErr == nil && width > 0 {
				return width
			}
		}
	}

	// NEW: fallback to `tput cols`.
	output, err = exec.Command("tput", "cols").Output()
	if err == nil {
		widthText := strings.TrimSpace(string(output))

		width, convErr := strconv.Atoi(widthText)
		if convErr == nil && width > 0 {
			return width
		}
	}

	// NEW: fallback to COLUMNS.
	if cols := os.Getenv("COLUMNS"); cols != "" {
		width, err := strconv.Atoi(cols)
		if err == nil && width > 0 {
			return width
		}
	}

	return 80
}
