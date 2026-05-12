package test

import (
	"bytes"
	"strings"
	"testing"

	"ascii-art/internal"
)

// Ελέγχει ότι κενή είσοδος δεν τυπώνει τίποτα
func TestEmptyLine(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(&buf, []string{""}, "../banners/standard.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

// Ελέγχει ότι η έξοδος για ένα string έχει ακριβώς 8 γραμμές ASCII art
func TestSingleWord(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(&buf, []string{"Hi"}, "../banners/standard.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) != 8 {
		t.Errorf("expected 8 lines of ASCII art, got %d", len(lines))
	}
}

// Ελέγχει ότι το \n χωρίζει σωστά σε δύο blocks ASCII art με μία κενή γραμμή
func TestNewlineSeparator(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(&buf, []string{"Hi", "", "There"}, "../banners/standard.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// "Hi" = 8 γραμμές, "" = 1 κενή γραμμή, "There" = 8 γραμμές → σύνολο 17
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) != 17 {
		t.Errorf("expected 17 lines (8 + 1 + 8), got %d", len(lines))
	}
}

func TestSpecificCharacter(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(&buf, []string{"A"}, "../banners/standard.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	lines := strings.Split(output, "\n")
	if !strings.Contains(lines[4], "_") {
		t.Errorf("expected 5th line of 'A' to contain '_', got %q", lines[4])
	}
}
