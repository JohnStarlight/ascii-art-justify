package internal

import (
	"fmt"
	"strings"
)

// Config stores validated application configuration parsed from CLI arguments.
type Config struct {
	Text       string
	OutputFile string
	BannerPath string
}

// bannerFiles maps banner names to their corresponding banner file paths.
var bannerFiles = map[string]string{
	"standard":   "banners/standard.txt",
	"shadow":     "banners/shadow.txt",
	"thinkertoy": "banners/thinkertoy.txt",
}

// ParseArgs validates and parses command-line arguments into a Config.
func ParseArgs(args []string) (Config, error) {
	if len(args) != 2 && len(args) != 3 && len(args) != 4 {
		return Config{}, fmt.Errorf(
			"usage:\n"+
				"go run ./cmd [STRING]\n"+
				"go run ./cmd [STRING] [BANNER]\n"+
				"go run ./cmd --output=<fileName.txt> [STRING] [BANNER]",
		)
	}

	var text string
	var outputFile string
	var banner string

	switch len(args) {
	case 2:
		text = args[1]
		banner = "standard"

	case 3:
		text = args[1]
		banner = args[2]

	case 4:
		if !strings.HasPrefix(args[1], "--output=") {
			return Config{}, fmt.Errorf(
				"invalid output flag: expected --output=<fileName.txt>",
			)
		}

		outputFile = strings.TrimPrefix(args[1], "--output=")

		if outputFile == "" {
			return Config{}, fmt.Errorf("output filename cannot be empty")
		}

		text = args[2]
		banner = args[3]
	}

	bannerPath, exists := bannerFiles[banner]
	if !exists {
		return Config{}, fmt.Errorf(
			"unsupported banner %q",
			banner,
		)
	}

	// Only printable ASCII characters (32–126) are supported.
	for _, r := range text {
		if r == '\\' {
			continue
		}

		if r < 32 || r > 126 {
			return Config{}, fmt.Errorf(
				"invalid character %q (only printable ASCII is supported)",
				r,
			)
		}
	}

	return Config{
		Text:       text,
		OutputFile: outputFile,
		BannerPath: bannerPath,
	}, nil
}