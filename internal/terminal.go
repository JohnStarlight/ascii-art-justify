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
	// NEW: first try the COLUMNS environment variable.
	// Some terminals/shells already store the terminal width there.
	if cols := os.Getenv("COLUMNS"); cols != "" {
		width, err := strconv.Atoi(cols)
		if err == nil && width > 0 {
			return width
		}
	}

	// NEW: fallback to `tput cols`, which asks the terminal
	// how many columns it currently has.
	output, err := exec.Command("tput", "cols").Output()
	if err == nil {
		widthText := strings.TrimSpace(string(output))

		width, convErr := strconv.Atoi(widthText)
		if convErr == nil && width > 0 {
			return width
		}
	}

	// NEW: safe fallback width if the program cannot detect the terminal size.
	return 80
}
