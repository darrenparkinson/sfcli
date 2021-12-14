package cmd

import (
	"github.com/spf13/cobra"
)

var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Bulk API V2 Commands",
}

func init() {
	rootCmd.AddCommand(bulkCmd)
}
