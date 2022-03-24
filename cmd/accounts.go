package cmd

import (
	"github.com/spf13/cobra"
)

// accountsCmd represents the accounts command
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "account related commands",
}

func init() {
	// rootCmd.AddCommand(accountsCmd)
}
