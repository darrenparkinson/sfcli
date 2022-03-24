package cmd

import (
	"github.com/spf13/cobra"
)

// used for both bulk insert and bulk upsert
var file string
var sobject string
var crlfLineEnding bool

var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Bulk API V2 Commands",
}

func init() {
	rootCmd.AddCommand(bulkCmd)
}
