package cmd

import (
	"fmt"
	"os"

	"github.com/3-2-1-contact/csv-tools/internal/version"
	"github.com/spf13/cobra"
)

var verbose bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("csvtools version %s\n", version.Version)
		fmt.Printf("commit: %s\n", version.CommitHash)
		fmt.Printf("built: %s\n", version.BuildTime)
	},
}

var rootCmd = &cobra.Command{
	Use:   "csvtools",
	Short: "A CSV processing tool",
	Long: `csvtools is a command line tool for processing CSV files.
It supports merging multiple CSV files and splitting them based on criteria.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.AddCommand(versionCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
}
