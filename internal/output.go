package internal

import (
	"fmt"
	"os"
)

// PrepareOutputFile creates or truncates
// the specified output file.
//
// If the file already exists,
// its contents are erased.
func PrepareOutputFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create output file: %w",
			err,
		)
	}

	return file, nil
}