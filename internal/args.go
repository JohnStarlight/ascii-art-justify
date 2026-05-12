package internal

import (
	"fmt"
	"strings"
)

type Config struct {
	Text string
	OutputFile string
	BannerPath string
}

func ParseArgs(args []string) (Config, error) {
if len(args) != 2 && len(args) != 3 && len(args) != 4 {
	return Config{}, fmt.Errorf(
    "invalid usage: expected 1, 2 or 3 arguments",
)
}

var text string

var outputFile string

var banner string

bannerFiles := map[string]string{
	"standard":    "banners/standard.txt",
	"shadow":      "banners/shadow.txt",
	"thinkertoy":  "banners/thinkertoy.txt",
}

if len(args) == 2 {
	text = args[1]
	banner = "standard"
}

if len(args) == 3 {
	text = args[1]
	banner = args[2]
}

if len(args) == 4 {
	if !strings.HasPrefix(args[1], "--output=") {
	return Config{}, fmt.Errorf(
		"invalid output flag: expected --output=<fileName.txt>",
	)
}

	outputFile = strings.TrimPrefix(args[1], "--output=")
	text = args[2]
	banner = args[3]
}

value, exists := bannerFiles[banner]
if !exists {
	return Config{}, fmt.Errorf(
		"unsupported banner %q",
		banner,
	)
}
return Config{
	Text: text,
	OutputFile: outputFile,
	BannerPath: value,
}, nil

}