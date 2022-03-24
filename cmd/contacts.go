package cmd

import (
	"github.com/spf13/cobra"
)

var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "contact related commands",
}

func init() {
	// rootCmd.AddCommand(contactsCmd)
}
