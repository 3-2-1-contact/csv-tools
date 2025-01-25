package split

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type FilterConfig struct {
	Column string
	Values map[string]bool
}

// Execute is the main entry point
func Execute(flags Flags) error {
	// Parse format string if provided
	var err error
	if flags.FormatStr != "" {
		flags.Formats, err = ParseFormatString(flags.FormatStr)
		if err != nil {
			return fmt.Errorf("invalid format specification: %w", err)
		}
	}
	if err := validateFlags(flags); err != nil {
		return err
	}

	reader, closer, err := getReader(flags.InputFile)
	if err != nil {
		return err
	}
	if closer != nil {
		defer closer.Close()
	}

	filterConfig, err := parseFilterArgument(flags.Filter)
	if err != nil {
		return err
	}

	headers, indices, err := processHeaders(reader, filterConfig.Column, flags.SplitCol)
	if err != nil {
		return err
	}
	for i := range headers {
		headers[i] = strings.ToLower(headers[i])
	}
	var formats []ColumnFormat
	if flags.FormatStr != "" {
		formats, err = ParseFormatString(flags.FormatStr)
		if err != nil {
			return fmt.Errorf("invalid format specification: %w", err)
		}
	}
	groupedRecords, err := groupRecords(reader, indices, filterConfig.Values)
	if err != nil {
		return err
	}

	//filterValueStr := createFilterValueString(flags.Filter)
	if err := writeOutputFiles(headers, flags.Filter, groupedRecords, formats); err != nil {
		return fmt.Errorf("error writing output files: %w", err)
	}
	return nil
}

func validateFlags(args Flags) error {
	if args.Filter == "" || args.SplitCol == "" {
		return fmt.Errorf("filter and split column must be specified")
	}
	return nil
}

func getReader(inputFile string) (*csv.Reader, io.Closer, error) {
	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			return nil, nil, fmt.Errorf("error opening input file: %w", err)
		}
		return csv.NewReader(file), file, nil
	}
	return csv.NewReader(os.Stdin), nil, nil
}

func parseFilterArgument(filter string) (FilterConfig, error) {
	filterParts := strings.Split(filter, "=")
	if len(filterParts) != 2 {
		return FilterConfig{}, fmt.Errorf("filter must be in format column=value1,value2,...")
	}

	filterValues := strings.Split(filterParts[1], ",")
	filterValuesMap := make(map[string]bool)
	for _, value := range filterValues {
		filterValuesMap[strings.ToLower(strings.TrimSpace(value))] = true
	}

	return FilterConfig{
		Column: filterParts[0],
		Values: filterValuesMap,
	}, nil
}

type ColumnIndices struct {
	FilterIdx int
	SplitIdx  int
}

func processHeaders(reader *csv.Reader, filterColumn, splitColumn string) ([]string, ColumnIndices, error) {
	headers, err := reader.Read()
	if err != nil {
		return nil, ColumnIndices{}, fmt.Errorf("error reading headers: %w", err)
	}

	indices := ColumnIndices{FilterIdx: -1, SplitIdx: -1}
	for i, header := range headers {
		if strings.EqualFold(header, filterColumn) {
			indices.FilterIdx = i
		}
		if strings.EqualFold(header, splitColumn) {
			indices.SplitIdx = i
		}
	}

	if indices.FilterIdx == -1 {
		return nil, ColumnIndices{}, fmt.Errorf("filter column '%s' not found in headers", filterColumn)
	}
	if indices.SplitIdx == -1 {
		return nil, ColumnIndices{}, fmt.Errorf("split column '%s' not found in headers", splitColumn)
	}

	return headers, indices, nil
}

func groupRecords(reader *csv.Reader, indices ColumnIndices, filterValues map[string]bool) (map[string][][]string, error) {
	groupedRecords := make(map[string][][]string)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				break
			}
			return nil, fmt.Errorf("error reading record: %w", err)
		}

		recordFilterValue := strings.ToLower(record[indices.FilterIdx])
		if filterValues[recordFilterValue] {
			splitValue := strings.ToLower(record[indices.SplitIdx])
			groupedRecords[splitValue] = append(groupedRecords[splitValue], record)
		}
	}
	return groupedRecords, nil
}

func createFilterValueString(filter string) string {
	parts := strings.Split(filter, "=")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func writeOutputFiles(headers []string, filter string, recordGroups map[string][][]string, formats []ColumnFormat) error {
	filter = strings.Replace(filter, " ", "-", -1)
	filter = strings.Replace(filter, "=", "_", -1)
	filter = strings.Replace(filter, ",", "_", -1)
	for value, records := range recordGroups {

		outputFilename := fmt.Sprintf("%s-%s.csv", value, strings.ToLower(filter))
		if err := writeCSVFile(outputFilename, headers, records, formats); err != nil {
			return err
		}
	}
	return nil
}

// func writeCSVFile(filename string, headers []string, records [][]string) error {
func writeCSVFile(filename string, headers []string, records [][]string, formats []ColumnFormat) error {
	outputFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filename, err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("error writing headers to %s: %w", filename, err)
	}

	// Format specific columns in the records
	formattedRecords := make([][]string, len(records))
	for i, record := range records {
		formattedRecord := make([]string, len(record))
		for j, value := range record {
			switch j {
			case 0: // First two columns
				// Try to parse as float and format with required precision
				if f, err := strconv.ParseFloat(value, 64); err == nil {
					formattedRecord[j] = fmt.Sprintf("%.6f", f)
				} else {
					formattedRecord[j] = value // Keep original if not a valid float
				}
			case 1:
				if f, err := strconv.ParseFloat(value, 64); err == nil {
					formattedRecord[j] = fmt.Sprintf("%.5f", f)
				} else {
					formattedRecord[j] = value // Keep original if not a valid float
				}
			default:
				formattedRecord[j] = value // Keep other columns as is
			}
		}
		formattedRecords[i] = formattedRecord
	}

	if err := writer.WriteAll(formattedRecords); err != nil {
		return fmt.Errorf("error writing records to %s: %w", filename, err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing writer for %s: %w", filename, err)
	}

	return nil
}
