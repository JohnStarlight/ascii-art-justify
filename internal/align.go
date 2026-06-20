package internal

import "strings"

// NEW: alignRows receives already-rendered ASCII rows and adds spaces
// in front of each row depending on the selected alignment.
func alignRows(rows []string, align string, terminalWidth int) []string {
	if align == "" || align == "left" {
		return rows
	}

	aligned := make([]string, 0, len(rows))

	for _, row := range rows {
		padding := terminalWidth - len(row)

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
	_ = config // NEW: placeholder so we can later support justify together with color options.

	words := strings.Fields(line)

	// A single word cannot be justified.
	if len(words) < 2 {
		return renderWordRows(line, bannerLines)
	}

	wordRows := make([][]string, 0, len(words)) // NEW: stores the 8 rendered rows for each word.
	wordWidths := make([]int, 0, len(words))    // NEW: stores the real ASCII width of each word block.
	totalWordWidth := 0                         // NEW: total width of all word blocks without spaces.

	for _, word := range words {
		renderedWordRows := renderWordRows(word, bannerLines)
		width := maxRowWidth(renderedWordRows) // NEW: use the widest row, not only row 0.

		wordRows = append(wordRows, renderedWordRows)
		wordWidths = append(wordWidths, width)

		totalWordWidth += width
	}

	gaps := len(words) - 1
	spacesToDistribute := terminalWidth - totalWordWidth

	if spacesToDistribute <= 0 {
		return renderWordRows(line, bannerLines)
	}

	spacePerGap := spacesToDistribute / gaps
	extraSpaces := spacesToDistribute % gaps

	justified := make([]string, 0, charHeight) // NEW: final 8 rows after spreading spaces between words.

	for row := 0; row < charHeight; row++ {
		var sb strings.Builder

		for wordIndex, rows := range wordRows {
			sb.WriteString(padRight(rows[row], wordWidths[wordIndex])) // NEW: keep each word block aligned to its real width.

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

// NEW: renderWordRows renders one word into its 8 ASCII-art rows.
// Justify needs this because it must place spaces between whole word blocks.
func renderWordRows(
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
	padding := width - len(row)

	if padding <= 0 {
		return row
	}

	return row + strings.Repeat(" ", padding)
}
