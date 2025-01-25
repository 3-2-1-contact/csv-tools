package split

type FormatType string

const (
	FormatString  FormatType = "string"
	FormatInt     FormatType = "int"
	FormatDecimal FormatType = "decimal"
)

type ColumnFormat struct {
	Name      string
	Type      FormatType
	Precision int
}

type Flags struct {
	InputFile string
	Filter    string
	SplitCol  string
	FormatStr string         // Raw format string from command line
	Formats   []ColumnFormat // Parsed format specifications
}
