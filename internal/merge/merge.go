package merge

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Flags struct {
	InputFiles string
	OutputFile string
	NewColumn  string
}

// compareHeaders checks if two string slices are equal.
func compareHeaders(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Execute(args Flags) error {
	if args.InputFiles == "" {
		return errors.New("no input files specified")
	}

	files := strings.Split(args.InputFiles, ",")
	if len(files) == 0 {
		return errors.New("no input files specified")
	}

	var writer *csv.Writer
	if args.OutputFile != "" {
		outputFile, err := os.Create(args.OutputFile)
		if err != nil {
			return fmt.Errorf("error creating output file: %w", err)
		}
		defer outputFile.Close()
		writer = csv.NewWriter(outputFile)
	} else {
		writer = csv.NewWriter(os.Stdout)
	}
	defer writer.Flush()

	var headers []string
	isFirst := true

	for _, filename := range files {
		filename = strings.TrimSpace(filename)
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("error opening file %s: %w", filename, err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		fileHeaders, err := reader.Read()
		if err != nil {
			return fmt.Errorf("error reading headers from %s: %w", filename, err)
		}

		// Get the value for the new column from the filename
		fileValue := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

		if isFirst {
			// Add new column to headers if specified
			if args.NewColumn != "" {
				headers = append(fileHeaders, args.NewColumn)
			} else {
				headers = fileHeaders
			}
			if err := writer.Write(headers); err != nil {
				return fmt.Errorf("error writing headers: %w", err)
			}
			isFirst = false
		} else if !compareHeaders(fileHeaders, headers[:len(fileHeaders)]) {
			return fmt.Errorf("headers in %s don't match the first file", filename)
		}

		// Copy all records except headers, adding the new column value
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error reading from %s: %w", filename, err)
			}

			// Add the new column value if specified
			if args.NewColumn != "" {
				record = append(record, fileValue)
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing record: %w", err)
			}
		}
	}

	return nil
}
