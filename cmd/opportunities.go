package cmd

import (
	"github.com/spf13/cobra"
)

var opportunitiesCmd = &cobra.Command{
	Use:   "opportunities",
	Short: "opportunity related commands",
}

func init() {
	rootCmd.AddCommand(opportunitiesCmd)
}
