package internal

import (
	"fmt"
	"strings"
)

type Config struct {
	Text       string
	Color      string
	Part       string
	OutputFile string
	Align      string // NEW: stores the selected alignment from --align=<type>
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
	if len(args) < 2 {
		return Config{}, usageError()
	}

	var colorFlag, outputFlag, alignFlag string // NEW: alignFlag temporarily stores --align value
	var positionals []string
	outputProvided := false

	for _, arg := range args[1:] {
		lower := strings.ToLower(arg)

		switch {
		case strings.HasPrefix(lower, "--color="):
			colorFlag = arg[len("--color="):]
			if colorFlag == "" {
				return Config{}, colorUsageError()
			}

		case strings.HasPrefix(lower, "--align="): // NEW: accepts exactly --align=<type>
			alignFlag = strings.ToLower(arg[len("--align="):])
			if alignFlag == "" {
				return Config{}, alignUsageError()
			}

		case strings.HasPrefix(lower, "--align"): // NEW: rejects wrong formats like --align center or --align
			return Config{}, alignUsageError()

		case strings.HasPrefix(lower, "--color"):
			return Config{}, colorUsageError()

		case strings.HasPrefix(lower, "--output="):
			outputProvided = true
			outputFlag = arg[len("--output="):]

		default:
			positionals = append(positionals, arg)
		}
	}

	if outputProvided && outputFlag == "" {
		return Config{}, fmt.Errorf("output filename cannot be empty")
	}

	var text, part string
	banner := "standard"

	if colorFlag != "" {
		switch len(positionals) {
		case 1:
			text = positionals[0]
		case 2:
			if _, ok := bannerFiles[strings.ToLower(positionals[1])]; ok {
				text = positionals[0]
				banner = strings.ToLower(positionals[1])
			} else {
				part = positionals[0]
				text = positionals[1]
			}
		case 3:
			part = positionals[0]
			text = positionals[1]
			banner = strings.ToLower(positionals[2])
		default:
			return Config{}, usageError()
		}
	} else {
		switch len(positionals) {
		case 1:
			text = positionals[0]
		case 2:
			text = positionals[0]
			banner = strings.ToLower(positionals[1])
		default:
			return Config{}, usageError()
		}
	}

	bannerPath, exists := bannerFiles[banner]
	if !exists {
		return Config{}, fmt.Errorf("unsupported banner %q", banner)
	}

	// NEW: validates --align value after all arguments have been parsed.
	// Only the four alignments required by the subject are accepted.
	if alignFlag != "" {
		switch alignFlag {
		case "left", "right", "center", "justify":
		// valid alignment, continue normally
		default:
			return Config{}, alignUsageError()
		}
	}

	if err := validateText(text); err != nil {
		return Config{}, err
	}

	var color string
	if colorFlag != "" {
		var err error
		color, err = Palette(colorFlag)
		if err != nil {
			return Config{}, err
		}
	}

	return Config{
		Text:       text,
		Color:      color,
		Part:       part,
		OutputFile: outputFlag,
		Align:      alignFlag, // NEW: saves the alignment inside Config so main/renderer can use it later
		BannerPath: bannerPath,
	}, nil
}

func colorUsageError() error {
	return fmt.Errorf(
		"Usage: go run . [OPTION] [STRING]\n\n" +
			"EX: go run . --color=<color> <substring to be colored> \"something\"",
	)
}

// NEW: specific usage message required by the ascii-art-justify subject.
func alignUsageError() error {
	return fmt.Errorf(
		"Usage: go run . [OPTION] [STRING] [BANNER]\n\n" +
			"EX: go run . --align=right something standard",
	)
}

func usageError() error {
	return fmt.Errorf(
		"usage:\n" +
			"  go run ./cmd [STRING]\n" +
			"  go run ./cmd [STRING] [BANNER]\n" +
			"  go run ./cmd --output=<file.txt> [STRING] [BANNER]\n" +
			"  go run ./cmd --color=<color> [STRING]\n" +
			"  go run ./cmd --color=<color> [SUBSTRING] [STRING]\n" +
			"  go run ./cmd --color=<color> [SUBSTRING] [STRING] [BANNER]\n" +
			"  go run ./cmd --align=<type> [STRING] [BANNER]",
	)
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
