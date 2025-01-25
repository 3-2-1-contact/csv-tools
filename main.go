package main

import (
	"os"

	"github.com/3-2-1-contact/csv-tools/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
