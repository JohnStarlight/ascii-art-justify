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
