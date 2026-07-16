package test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ascii-art/internal"
)

func TestEmptyLine(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(
		&buf,
		internal.Config{ // NEW: PrintAscii now receives Config instead of separate color/part/banner arguments.
			BannerPath: "../banners/standard.txt",
		},
		[]string{""},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

func TestSingleWord(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(
		&buf,
		internal.Config{ // NEW: banner path is now passed through Config.
			BannerPath: "../banners/standard.txt",
		},
		[]string{"Hi"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// TrimRight strips the trailing newline that Fprintln appends to each row,
	// otherwise Split produces a spurious empty element and len == 9.
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	if len(lines) != 8 {
		t.Errorf("expected 8 lines of ASCII art, got %d", len(lines))
	}
}

func TestNewlineSeparator(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(
		&buf,
		internal.Config{ // NEW: keeps renderer inputs grouped in Config for extensions like --align.
			BannerPath: "../banners/standard.txt",
		},
		[]string{"Hi", "", "There"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// "Hi"     -> 8 lines
	// ""       -> 1 empty line
	// "There"  -> 8 lines
	//
	// Total: 17 lines
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	if len(lines) != 17 {
		t.Errorf("expected 17 lines (8 + 1 + 8), got %d", len(lines))
	}
}

func TestSpecificCharacter(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(
		&buf,
		internal.Config{ // NEW: old empty color/part arguments are no longer needed.
			BannerPath: "../banners/standard.txt",
		},
		[]string{"A"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	lines := strings.Split(output, "\n")

	// lines[4] is the 5th row (0-indexed); 'A' has its crossbar there.
	if !strings.Contains(lines[4], "_") {
		t.Errorf(
			"expected 5th line of 'A' to contain '_', got %q",
			lines[4],
		)
	}
}

func TestParseArgsOutputFlag(t *testing.T) {
	cfg, err := internal.ParseArgs([]string{"prog", "--output=out.txt", "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFile != "out.txt" {
		t.Errorf("OutputFile = %q, want %q", cfg.OutputFile, "out.txt")
	}
}

func TestParseArgsOutputEmptyValue(t *testing.T) {
	_, err := internal.ParseArgs([]string{"prog", "--output=", "Hello"})
	if err == nil {
		t.Fatal("expected error for --output= with empty value, got nil")
	}
}

func TestPrintAsciiToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.txt")

	file, err := internal.PrepareOutputFile(path)
	if err != nil {
		t.Fatalf("PrepareOutputFile: %v", err)
	}
	defer file.Close()

	err = internal.PrintAscii(file, internal.Config{BannerPath: "../banners/standard.txt"}, []string{"Hi"})
	if err != nil {
		t.Fatalf("PrintAscii: %v", err)
	}
	file.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) != 8 {
		t.Errorf("expected 8 lines in output file, got %d", len(lines))
	}
}

func TestPrintAsciiOutOfRangeChar(t *testing.T) {
	var buf bytes.Buffer
	err := internal.PrintAscii(
		&buf,
		internal.Config{BannerPath: "../banners/standard.txt"},
		[]string{"\x80"},
	)
	if err == nil {
		t.Fatal("expected error for out-of-range character, got nil")
	}
}

func TestParseArgsInvalidEscape(t *testing.T) {
	_, err := internal.ParseArgs([]string{"prog", "Hello\\xWorld"})
	if err == nil {
		t.Fatal("expected error for invalid escape sequence, got nil")
	}
}

func TestParseArgsTrailingBackslash(t *testing.T) {
	_, err := internal.ParseArgs([]string{"prog", "Hello\\"})
	if err == nil {
		t.Fatal("expected error for trailing backslash, got nil")
	}
}

func TestPrintAsciiCorruptBanner(t *testing.T) {
	path := filepath.Join(t.TempDir(), "corrupt.txt")
	os.WriteFile(path, []byte("not a valid banner\n"), 0o600)

	var buf bytes.Buffer
	err := internal.PrintAscii(
		&buf,
		internal.Config{BannerPath: path},
		[]string{"Hi"},
	)
	if err == nil {
		t.Fatal("expected error for corrupt banner file, got nil")
	}
}
