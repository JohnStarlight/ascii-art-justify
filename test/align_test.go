package test

import (
	"bytes"
	"strings"
	"testing"

	"ascii-art/internal"
)

// testVisibleLen counts characters that occupy terminal columns,
// skipping ANSI escape sequences, mirroring what a terminal displays.
func testVisibleLen(s string) int {
	count := 0
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			continue
		}

		if inEscape {
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}

		count++
	}

	return count
}

func renderAligned(t *testing.T, config internal.Config, text string) []string {
	t.Helper()

	config.BannerPath = "../banners/standard.txt"

	var buf bytes.Buffer
	if err := internal.PrintAscii(&buf, config, []string{text}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
}

func TestParseArgsAlignValues(t *testing.T) {
	for _, align := range []string{"left", "right", "center", "justify"} {
		t.Run(align, func(t *testing.T) {
			cfg, err := internal.ParseArgs([]string{"prog", "--align=" + align, "Hello"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.Align != align {
				t.Errorf("Align = %q, want %q", cfg.Align, align)
			}
		})
	}
}

func TestParseArgsAlignInvalidValue(t *testing.T) {
	_, err := internal.ParseArgs([]string{"prog", "--align=weird", "Hello"})
	if err == nil {
		t.Fatal("expected error for invalid alignment, got nil")
	}
}

func TestParseArgsAlignMissingEquals(t *testing.T) {
	_, err := internal.ParseArgs([]string{"prog", "--align", "right", "Hello"})
	if err == nil {
		t.Fatal("expected error for --align without '=', got nil")
	}
}

func TestParseArgsAlignEmptyValue(t *testing.T) {
	_, err := internal.ParseArgs([]string{"prog", "--align=", "Hello"})
	if err == nil {
		t.Fatal("expected error for --align= with empty value, got nil")
	}
}

// NEW: --align and --color are independent flags and must combine cleanly.
func TestParseArgsAlignWithColor(t *testing.T) {
	cfg, err := internal.ParseArgs([]string{"prog", "--align=justify", "--color=red", "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Align != "justify" {
		t.Errorf("Align = %q, want %q", cfg.Align, "justify")
	}
	if len(cfg.Colors) != 1 {
		t.Errorf("Colors = %v, want one entry", cfg.Colors)
	}
}

// Left alignment (and no alignment at all) must leave the render untouched.
func TestAlignLeftMatchesDefault(t *testing.T) {
	plain := renderAligned(t, internal.Config{}, "Hi")
	left := renderAligned(t, internal.Config{Align: "left"}, "Hi")

	if strings.Join(plain, "\n") != strings.Join(left, "\n") {
		t.Errorf("expected --align=left output to match unaligned output")
	}
}

func TestAlignRightPadsToTerminalWidth(t *testing.T) {
	width := internal.TerminalWidth()

	plain := renderAligned(t, internal.Config{}, "Hi")
	right := renderAligned(t, internal.Config{Align: "right"}, "Hi")

	if width <= len(plain[0]) {
		t.Skipf("terminal too narrow (%d cols) to exercise right alignment", width)
	}

	for i, line := range right {
		want := strings.Repeat(" ", width-len(plain[i])) + plain[i]
		if line != want {
			t.Errorf("line %d: right-aligned output mismatch\ngot  %q\nwant %q", i, line, want)
		}
	}
}

func TestAlignCenterPadsHalfway(t *testing.T) {
	width := internal.TerminalWidth()

	plain := renderAligned(t, internal.Config{}, "Hi")
	center := renderAligned(t, internal.Config{Align: "center"}, "Hi")

	if width <= len(plain[0]) {
		t.Skipf("terminal too narrow (%d cols) to exercise center alignment", width)
	}

	for i, line := range center {
		want := strings.Repeat(" ", (width-len(plain[i]))/2) + plain[i]
		if line != want {
			t.Errorf("line %d: center-aligned output mismatch\ngot  %q\nwant %q", i, line, want)
		}
	}
}

// Justify must stretch every row to exactly the terminal width by
// spreading spaces between the words.
func TestJustifySpansTerminalWidth(t *testing.T) {
	width := internal.TerminalWidth()

	if width < 40 {
		t.Skipf("terminal too narrow (%d cols) to exercise justify", width)
	}

	justified := renderAligned(t, internal.Config{Align: "justify"}, "Hi Yo")

	if len(justified) != 8 {
		t.Fatalf("expected 8 rendered rows, got %d", len(justified))
	}

	for i, line := range justified {
		if got := testVisibleLen(line); got != width {
			t.Errorf("line %d: expected visible width %d, got %d", i, width, got)
		}
	}
}

// A single word cannot be justified, so the output must match plain rendering.
func TestJustifySingleWordUnchanged(t *testing.T) {
	plain := renderAligned(t, internal.Config{}, "Hi")
	justified := renderAligned(t, internal.Config{Align: "justify"}, "Hi")

	if strings.Join(plain, "\n") != strings.Join(justified, "\n") {
		t.Errorf("expected single-word justify output to match plain output")
	}
}

// NEW: justify must preserve substring coloring while stretching rows.
func TestJustifyKeepsColors(t *testing.T) {
	width := internal.TerminalWidth()

	if width < 40 {
		t.Skipf("terminal too narrow (%d cols) to exercise justify", width)
	}

	red, err := internal.Palette("red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	justified := renderAligned(
		t,
		internal.Config{
			Align:  "justify",
			Colors: []string{red},
			Parts:  []string{"Hi"},
		},
		"Hi Yo",
	)

	for i, line := range justified {
		if !strings.Contains(line, red) {
			t.Errorf("line %d: expected red color code in justified output, got %q", i, line)
		}
		if got := testVisibleLen(line); got != width {
			t.Errorf("line %d: expected visible width %d, got %d", i, width, got)
		}
	}
}

func TestJustifyThreeWords(t *testing.T) {
	width := internal.TerminalWidth()

	if width < 60 {
		t.Skipf("terminal too narrow (%d cols) to exercise justify with three words", width)
	}

	justified := renderAligned(t, internal.Config{Align: "justify"}, "Hi Yo Za")

	if len(justified) != 8 {
		t.Fatalf("expected 8 rendered rows, got %d", len(justified))
	}

	for i, line := range justified {
		if got := testVisibleLen(line); got != width {
			t.Errorf("line %d: expected visible width %d, got %d", i, width, got)
		}
	}
}

func TestMultiLineRightAlign(t *testing.T) {
	width := internal.TerminalWidth()

	if width < 20 {
		t.Skipf("terminal too narrow (%d cols)", width)
	}

	var buf bytes.Buffer
	err := internal.PrintAscii(
		&buf,
		internal.Config{Align: "right", BannerPath: "../banners/standard.txt"},
		[]string{"Hi", "", "Yo"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 17 {
		t.Fatalf("expected 17 lines (8 + blank + 8), got %d", len(lines))
	}

	for i, line := range lines {
		if i == 8 {
			continue
		}
		if !strings.HasPrefix(line, " ") {
			t.Errorf("line %d: expected right-aligned line to start with spaces, got %q", i, line)
		}
	}
}

func TestAlignRightWithColor(t *testing.T) {
	width := internal.TerminalWidth()

	if width < 20 {
		t.Skipf("terminal too narrow (%d cols)", width)
	}

	red, err := internal.Palette("red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	right := renderAligned(t, internal.Config{
		Align:  "right",
		Colors: []string{red},
		Parts:  []string{""},
	}, "Hi")

	for i, line := range right {
		if !strings.Contains(line, red) {
			t.Errorf("line %d: expected ANSI color codes to be present", i)
		}
		if got := testVisibleLen(line); got != width {
			t.Errorf("line %d: expected visible width %d, got %d", i, width, got)
		}
	}
}
