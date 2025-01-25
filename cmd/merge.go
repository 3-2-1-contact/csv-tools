package cmd

import (
	"github.com/3-2-1-contact/csv-tools/internal/merge"

	"github.com/spf13/cobra"
)

var mergeArgs merge.Flags

var mergeCmd = &cobra.Command{
	Use:     "merge",
	Aliases: []string{"m"},
	Short:   "Merge multiple CSV files",
	Long: `Merge multiple CSV files into a single CSV file.
Optionally add a new column with values derived from filenames.`,
	Example: `  csvtools merge --col county --in "comal.csv,hays.csv" --out all.csv
  csvtools merge --col state --in "texas.csv,florida.csv" | csvtools split --filter "category=Aircraft" --split state`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return merge.Execute(mergeArgs)
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	mergeCmd.Flags().StringVar(&mergeArgs.InputFiles, "in", "", "Input CSV files (comma-separated)")
	mergeCmd.Flags().StringVar(&mergeArgs.OutputFile, "out", "", "Output CSV file (optional)")
	mergeCmd.Flags().StringVar(&mergeArgs.NewColumn, "col", "", "Name of new column to add with filename values")

	mergeCmd.MarkFlagRequired("in")
}
