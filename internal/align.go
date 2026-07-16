package internal

import "strings"

type wordBlock struct {
	text  string
	start int
}

// NEW: alignRows receives already-rendered ASCII rows and adds spaces
// in front of each row depending on the selected alignment.
func alignRows(rows []string, align string, terminalWidth int) []string {
	if align == "" || align == "left" {
		return rows
	}

	aligned := make([]string, 0, len(rows))

	for _, row := range rows {
		padding := terminalWidth - visibleLen(row)

		if padding <= 0 {
			aligned = append(aligned, row)
			continue
		}

		switch align {
		case "right":
			aligned = append(aligned, strings.Repeat(" ", padding)+row)

		case "center":
			leftPadding := padding / 2
			aligned = append(aligned, strings.Repeat(" ", leftPadding)+row)

		default:
			aligned = append(aligned, row)
		}
	}

	return aligned
}

// NEW: justifyRows spreads the available terminal space
// between words instead of adding padding to the left side.
func justifyRows(
	line string,
	bannerLines []string,
	terminalWidth int,
	config Config,
) []string {
	words := splitWordsWithPositions(line)

	// NEW: colorAt is computed once from the full line so each word keeps
	// the exact coloring it had at its original position.
	colorAt := colorIndexes(line, config)

	// A single word cannot be justified.
	if len(words) < 2 {
		return renderWordRows(line, 0, bannerLines, config, colorAt)
	}

	wordRows := make([][]string, 0, len(words))
	wordWidths := make([]int, 0, len(words))
	totalWordWidth := 0

	for _, word := range words {
		renderedWordRows := renderWordRows(
			word.text,
			word.start,
			bannerLines,
			config,
			colorAt,
		)

		plainWordRows := renderPlainWordRows(word.text, bannerLines)
		width := maxRowWidth(plainWordRows)

		wordRows = append(wordRows, renderedWordRows)
		wordWidths = append(wordWidths, width)

		totalWordWidth += width
	}

	gaps := len(words) - 1
	spacesToDistribute := terminalWidth - totalWordWidth

	if spacesToDistribute <= 0 {
		return renderWordRows(line, 0, bannerLines, config, colorAt)
	}

	spacePerGap := spacesToDistribute / gaps
	extraSpaces := spacesToDistribute % gaps

	justified := make([]string, 0, charHeight)

	for row := 0; row < charHeight; row++ {
		var sb strings.Builder

		for wordIndex, rows := range wordRows {
			sb.WriteString(padRight(rows[row], wordWidths[wordIndex]))

			if wordIndex < gaps {
				spaces := spacePerGap

				if wordIndex < extraSpaces {
					spaces++
				}

				sb.WriteString(strings.Repeat(" ", spaces))
			}
		}

		justified = append(justified, sb.String())
	}

	return justified
}

// NEW: splitWordsWithPositions keeps each word together with its original
// start position inside the input line. This is needed for substring coloring.
func splitWordsWithPositions(line string) []wordBlock {
	var words []wordBlock

	i := 0
	for i < len(line) {
		for i < len(line) && line[i] == ' ' {
			i++
		}

		if i >= len(line) {
			break
		}

		start := i

		for i < len(line) && line[i] != ' ' {
			i++
		}

		words = append(words, wordBlock{
			text:  line[start:i],
			start: start,
		})
	}

	return words
}

// NEW: renderWordRows renders text into its 8 ASCII-art rows.
// It is color-aware, so justify can work together with --color:
// colorAt was computed from the full line, and startPosition translates
// each character back to its absolute position inside that line.
func renderWordRows(
	text string,
	startPosition int,
	bannerLines []string,
	config Config,
	colorAt []int,
) []string {
	rows := make([]string, 0, charHeight)

	for row := 1; row <= charHeight; row++ {
		var sb strings.Builder

		for pos, r := range text {
			index := (int(r)-asciiStart)*linesPerChar + row
			segment := bannerLines[index]

			absolutePos := startPosition + pos

			if idx := colorAt[absolutePos]; idx >= 0 {
				sb.WriteString(config.Colors[idx])
				sb.WriteString(segment)
				sb.WriteString("\033[0m")
			} else {
				sb.WriteString(segment)
			}
		}

		rows = append(rows, sb.String())
	}

	return rows
}

// NEW: renderPlainWordRows renders without color.
// Justify uses this only to measure the visual width correctly.
func renderPlainWordRows(
	word string,
	bannerLines []string,
) []string {
	rows := make([]string, 0, charHeight)

	for row := 1; row <= charHeight; row++ {
		var sb strings.Builder

		for _, r := range word {
			index := (int(r)-asciiStart)*linesPerChar + row
			segment := bannerLines[index]

			sb.WriteString(segment)
		}

		rows = append(rows, sb.String())
	}

	return rows
}

// NEW: maxRowWidth returns the widest row inside an ASCII-art block.
func maxRowWidth(rows []string) int {
	max := 0

	for _, row := range rows {
		if len(row) > max {
			max = len(row)
		}
	}

	return max
}

// NEW: padRight adds spaces to the end of a row until it reaches width.
func padRight(row string, width int) string {
	padding := width - visibleLen(row)

	if padding <= 0 {
		return row
	}

	return row + strings.Repeat(" ", padding)
}

// NEW: visibleLen counts visible characters only.
// ANSI color codes are ignored because they do not take terminal space.
func visibleLen(s string) int {
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
