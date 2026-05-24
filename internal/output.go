package internal

import (
	"fmt"
	"os"
)

// Creates or truncates output file and returns a handle to it.
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
