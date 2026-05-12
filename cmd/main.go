package main

import (
	"fmt"
	"os"
	"strings"
	"io"

	"ascii-art/internal"
)

func main() {
	
	// Split on the literal two-character sequence "\n" (backslash + n),
	// which is how multi-line input is passed from the command line.
	// e.g. "Hello\nThere" becomes ["Hello", "There"].
	lines := strings.Split(text, "\\n")

	if err := internal.PrintAscii(writer, lines, value); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
