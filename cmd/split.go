package cmd

import (
	"github.com/spf13/cobra"

	"github.com/3-2-1-contact/csv-tools/internal/split"
)

var splitArgs split.Flags

var splitCmd = &cobra.Command{
	Use:     "split",
	Aliases: []string{"s"},
	Short:   "Split CSV file based on criteria",
	Long: `Split a CSV file into multiple files based on specified criteria.
Filter rows and split them into separate files based on column values.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return split.Execute(splitArgs)
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	splitCmd.Flags().StringVar(&splitArgs.InputFile, "in", "", "Input CSV file (optional, reads from stdin if not provided)")
	splitCmd.Flags().StringVar(&splitArgs.Filter, "filter", "", "Filter in format column=value1,value2,...")
	splitCmd.Flags().StringVar(&splitArgs.SplitCol, "split", "", "Column to split by")
	splitCmd.Flags().StringVar(&splitArgs.FormatStr, "format", "",
		`Column format specifications (e.g., "Total Value:decimal:6,'Account Number':int")
Available formats: string, int, decimal:<precision>`)

	err := splitCmd.MarkFlagRequired("filter")
	if err != nil {
		return
	}
	err = splitCmd.MarkFlagRequired("split")
	if err != nil {
		return
	}
}
