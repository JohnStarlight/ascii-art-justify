package main

import (
	"fmt"
	"os"
	"strings"
	"io"

	"ascii-art/internal"
)

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Println("Error: invalid usage")
		fmt.Println("Usage: go run ./cmd \"your-text-here\"")
		os.Exit(1)
	}

	var text string

	var outputFile string

	var banner string

	var writer io.Writer

	bannerFiles := map[string]string {
		"standard" : "banners/standard.txt",
		"shadow" : "banners/shadow.txt",
		"thinkertoy" : "banners/thinkertoy.txt",
	}

	if len(os.Args) == 2 {
		text = os.Args[1]
		banner = "standard"
	}
	if len(os.Args) == 3 {
		text = os.Args[1]
		banner = os.Args[2]
	}

	if len(os.Args) == 4 {
		if !strings.HasPrefix(os.Args[1], "--output="){
			fmt.Println("Error: invalid usage")
			fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
			os.Exit(1)
		}
		outputFile = strings.TrimPrefix(os.Args[1], "--output=")
		text = os.Args[2]
		banner = os.Args[3]
	}
	value, exists := bannerFiles[banner]
	 if !exists {
		fmt.Printf("Error: unsupported banner %q\n", banner)
		os.Exit(1)
	 }
	 
	 writer = os.Stdout

	if outputFile != "" {
		file, err := internal.PrepareOutputFile(outputFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		writer = file
		defer file.Close()
	}

	// Only printable ASCII characters (32–126) are supported.
	// Reject anything outside that range (e.g. accented letters, emoji).
	for _, r := range text {
		if r < 32 || r > 126 {
			fmt.Printf("Error: invalid character %q (only printable ASCII is supported)\n", r)
			os.Exit(1)
		}
	}

	// Split on the literal two-character sequence "\n" (backslash + n),
	// which is how multi-line input is passed from the command line.
	// e.g. "Hello\nThere" becomes ["Hello", "There"].
	lines := strings.Split(text, "\\n")

	if err := internal.PrintAscii(writer, lines, value); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
