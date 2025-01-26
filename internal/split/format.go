package split

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseFormatString parses format string like "Total Value:decimal:6,'Account Number':int".
func ParseFormatString(formatStr string) ([]ColumnFormat, error) {
	if formatStr == "" {
		return nil, nil
	}

	// Split on comma but not within quotes
	specs := splitUnquoted(formatStr, ',')
	formats := make([]ColumnFormat, len(specs), 0)

	for _, spec := range specs {
		parts := splitUnquoted(spec, ':')

		// Remove quotes from column name if present
		colName := strings.Trim(parts[0], `"' `)
		if colName == "" {
			return nil, fmt.Errorf("empty column name in format specification: %s", spec)
		}

		format := ColumnFormat{
			Name: colName,
			Type: FormatString, // default type
		}

		if len(parts) > 1 {
			formatType := strings.ToLower(strings.TrimSpace(parts[1]))
			switch formatType {
			case "string":
				format.Type = FormatString
			case "int":
				format.Type = FormatInt
			case "decimal":
				format.Type = FormatDecimal
				if len(parts) > 2 {
					precision, err := strconv.Atoi(strings.TrimSpace(parts[2]))
					if err != nil {
						return nil, fmt.Errorf("invalid precision for column %s: %s", colName, parts[2])
					}
					format.Precision = precision
				} else {
					format.Precision = 2 // default precision if not specified
				}
			default:
				return nil, fmt.Errorf("invalid format type for column %s: %s", colName, formatType)
			}
		}

		formats = append(formats, format)
	}

	return formats, nil
}

// Helper function to split string on delimiter but respect quotes.
func splitUnquoted(s string, delimiter rune) []string {
	var result []string
	var current strings.Builder
	inQuotes := false

	for _, r := range s {
		switch {
		case r == '"' || r == '\'':
			inQuotes = !inQuotes
		case r == delimiter && !inQuotes:
			if current.Len() > 0 {
				result = append(result, strings.TrimSpace(current.String()))
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}

	return result
}

// FormatValue formats a single value according to the specified format.
func FormatValue(value string, format ColumnFormat) (string, error) {
	switch format.Type {
	case FormatString:
		return value, nil
	case FormatInt:
		i, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return "", fmt.Errorf("invalid integer value: %s", value)
		}
		return strconv.Itoa(i), nil
	case FormatDecimal:
		f, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return "", fmt.Errorf("invalid decimal value: %s", value)
		}
		return fmt.Sprintf("%.*f", format.Precision, f), nil
	default:
		return value, nil
	}
}
