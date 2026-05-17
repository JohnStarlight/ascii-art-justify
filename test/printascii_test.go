package test

import (
	"bytes"
	"strings"
	"testing"

	"ascii-art/internal"
)

// TestEmptyLine verifies that an empty input
// produces no ASCII-art output.
func TestEmptyLine(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(
		&buf,
		[]string{""},
		"../banners/standard.txt",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

// TestSingleWord verifies that rendering
// a normal string produces exactly 8 lines
// of ASCII art.
func TestSingleWord(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(
		&buf,
		[]string{"Hi"},
		"../banners/standard.txt",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	if len(lines) != 8 {
		t.Errorf("expected 8 lines of ASCII art, got %d", len(lines))
	}
}

// TestNewlineSeparator verifies that "\n"
// correctly separates ASCII-art blocks
// with one empty line between them.
func TestNewlineSeparator(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(
		&buf,
		[]string{"Hi", "", "There"},
		"../banners/standard.txt",
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

// TestSpecificCharacter verifies
// that rendering the character 'A'
// produces expected ASCII-art content.
func TestSpecificCharacter(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(
		&buf,
		[]string{"A"},
		"../banners/standard.txt",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	lines := strings.Split(output, "\n")

	// Verify that the 5th visual row
	// contains an underscore,
	// which is expected for the letter 'A'.
	if !strings.Contains(lines[4], "_") {
		t.Errorf(
			"expected 5th line of 'A' to contain '_', got %q",
			lines[4],
		)
	}
}