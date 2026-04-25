package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"ascii-art/internal"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("The correct usage of this program is:")
		fmt.Println("  1. Navigate to the project root directory")
		fmt.Println("  2. Run: go run ./cmd \"your-text-here\"")
		return
	}

	text := os.Args[1]

	// Only printable ASCII characters (32–126) are supported.
	// Reject anything outside that range (e.g. accented letters, emoji).
	for _, r := range text {
		if r < 32 || r > 126 {
			fmt.Printf("Invalid character detected: %q (only printable ASCII is supported)\n", r)
			os.Exit(1)
		}
	}

	fmt.Println("In which style would you like that?")
	fmt.Println("1 = Standard")
	fmt.Println("2 = Shadow")
	fmt.Println("3 = Thinkertoy")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	input = strings.TrimSpace(input)

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > 3 {
		fmt.Println("Invalid choice. Please enter \"1\" for Standard, \"2\" for Shadow or \"3\" for Thinkertoy.")
		os.Exit(1)
	}

	banners := []string{"banners/standard.txt", "banners/shadow.txt", "banners/thinkertoy.txt"}
	filename := banners[choice-1]

	// Split on the literal two-character sequence "\n" (backslash + n),
	// which is how multi-line input is passed from the command line.
	// e.g. "Hello\nThere" becomes ["Hello", "There"].
	lines := strings.Split(text, "\\n")
	internal.PrintAscii(lines, filename)
}
