package internal

import (
	"fmt"
	"strings"
)

type Config struct {
	Text       string
	OutputFile string
	BannerPath string
}

var bannerFiles = map[string]string{
	"standard":   "banners/standard.txt",
	"shadow":     "banners/shadow.txt",
	"thinkertoy": "banners/thinkertoy.txt",
}

// ParseArgs receives args as a parameter rather than reading os.Args directly
// so tests can inject arbitrary argument lists without spawning a subprocess.
func ParseArgs(args []string) (Config, error) {
	// args[0] is always the binary name, so the minimum meaningful length is 2.
	if len(args) < 2 || len(args) > 4 {
		return Config{}, fmt.Errorf(
			"usage:\n" +
				"go run ./cmd [STRING]\n" +
				"go run ./cmd [STRING] [BANNER]\n" +
				"go run ./cmd --output=<fileName.txt> [STRING] [BANNER]",
		)
	}

	var text string
	var outputFile string
	banner := "standard"

	switch len(args) {
	case 2:
		text = args[1]

	case 3:
		text = args[1]
		banner = args[2]

	case 4:
		args[1] = strings.ToLower(args[1]) // accept --Output=, --OUTPUT=, etc.
		if !strings.HasPrefix(args[1], "--output=") {
			return Config{}, fmt.Errorf(
				"invalid output flag: expected --output=<fileName.txt>",
			)
		}

		outputFile = strings.TrimPrefix(args[1], "--output=")

		if outputFile == "" {
			return Config{}, fmt.Errorf(
				"output filename cannot be empty",
			)
		}

		text = args[2]
		banner = args[3]
	}

	banner = strings.ToLower(banner)
	bannerPath, exists := bannerFiles[banner]
	if !exists {
		return Config{}, fmt.Errorf(
			"unsupported banner %q",
			banner,
		)
	}

	if err := validateText(text); err != nil {
		return Config{}, err
	}

	return Config{
		Text:       text,
		OutputFile: outputFile,
		BannerPath: bannerPath,
	}, nil
}

func validateText(text string) error {
	// Byte iteration (not range) so we can manipulate i to skip the second
	// character of a \n escape sequence without rune-width complications.
	for i := 0; i < len(text); i++ {
		char := text[i]

		if char == '\\' {
			if i+1 >= len(text) {
				return fmt.Errorf(
					"invalid escape sequence: trailing backslash",
				)
			}

			// Only "\n" is supported.
			if text[i+1] != 'n' {
				return fmt.Errorf(
					"invalid escape sequence \\%c",
					text[i+1],
				)
			}

			i++ // advance past the 'n' to consume the full \n escape
			continue
		}

		if char < 32 || char > 126 {
			return fmt.Errorf(
				"invalid character %q (only printable ASCII is supported)",
				char,
			)
		}
	}
	return nil
}
