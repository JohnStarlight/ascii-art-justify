package internal

import (
	"os"
	"fmt"
)

func PrepareOutputFile(filename string) (*os.File, error){

	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	return file, nil
}
