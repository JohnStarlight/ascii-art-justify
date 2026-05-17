package internal

import (
	"fmt"
	"strings"
)

// Config stores validated application configuration
// parsed from command-line arguments.
type Config struct {
	// Text is the raw input string to render.
	Text string

	// OutputFile is the optional file
	// where ASCII art should be written.
	OutputFile string

	// BannerPath is the filesystem path
	// of the selected banner file.
	BannerPath string
}

// bannerFiles maps supported banner names
// to their corresponding banner file paths.
var bannerFiles = map[string]string{
	"standard":   "banners/standard.txt",
	"shadow":     "banners/shadow.txt",
	"thinkertoy": "banners/thinkertoy.txt",
}

// ParseArgs validates command-line arguments
// and converts them into a Config struct.
func ParseArgs(args []string) (Config, error) {
	// Valid command formats:
	//
	// go run ./cmd [STRING]
	// go run ./cmd [STRING] [BANNER]
	// go run ./cmd --output=<fileName.txt> [STRING] [BANNER]
	if len(args) != 2 && len(args) != 3 && len(args) != 4 {
		return Config{}, fmt.Errorf(
			"usage:\n" +
				"go run ./cmd [STRING]\n" +
				"go run ./cmd [STRING] [BANNER]\n" +
				"go run ./cmd --output=<fileName.txt> [STRING] [BANNER]",
		)
	}

	var text string
	var outputFile string
	var banner string

	switch len(args) {

	// Example:
	// go run ./cmd "hello"
	case 2:
		text = args[1]
		banner = "standard"

	// Example:
	// go run ./cmd "hello" shadow
	case 3:
		text = args[1]
		banner = args[2]

	// Example:
	// go run ./cmd --output=result.txt "hello" shadow
	case 4:
		// Validate output flag format.
		if !strings.HasPrefix(args[1], "--output=") {
			return Config{}, fmt.Errorf(
				"invalid output flag: expected --output=<fileName.txt>",
			)
		}

		// Remove the flag prefix
		// and keep only the filename.
		outputFile = strings.TrimPrefix(args[1], "--output=")

		// Prevent empty filenames.
		if outputFile == "" {
			return Config{}, fmt.Errorf(
				"output filename cannot be empty",
			)
		}

		text = args[2]
		banner = args[3]
	}

	// Verify that the selected banner exists.
	bannerPath, exists := bannerFiles[banner]
	if !exists {
		return Config{}, fmt.Errorf(
			"unsupported banner %q",
			banner,
		)
	}

	// Validate input characters.
	//
	// Allowed:
	// - printable ASCII characters (32–126)
	// - the escape sequence "\n"
	for i := 0; i < len(text); i++ {
		char := text[i]

		// Handle escape sequences beginning with '\'.
		if char == '\\' {

			// '\' cannot appear as the final character.
			if i+1 >= len(text) {
				return Config{}, fmt.Errorf(
					"invalid escape sequence: trailing backslash",
				)
			}

			// Only "\n" is supported.
			if text[i+1] != 'n' {
				return Config{}, fmt.Errorf(
					"invalid escape sequence \\%c",
					text[i+1],
				)
			}

			// Skip the validated 'n'.
			i++
			continue
		}

		// Reject non-printable ASCII characters.
		if char < 32 || char > 126 {
			return Config{}, fmt.Errorf(
				"invalid character %q (only printable ASCII is supported)",
				char,
			)
		}
	}

	// Return fully validated configuration.
	return Config{
		Text:       text,
		OutputFile: outputFile,
		BannerPath: bannerPath,
	}, nil
}