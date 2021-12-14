package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/darrenparkinson/sfcli/pkg/salesforce"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

//TODO: Abstract this to enable description of anything...

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "list field names for the various objects",
}

var describeAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "list account fields",
	Run:   accountDescribe,
}
var describeContactCmd = &cobra.Command{
	Use:   "contact",
	Short: "list contact fields",
	Run:   contactDescribe,
}
var describeOpportunityCmd = &cobra.Command{
	Use:   "opportunity",
	Short: "list opportunity fields",
	Run:   opportunityDescribe,
}
var describeUserCmd = &cobra.Command{
	Use:   "user",
	Short: "list user fields",
	Run:   userDescribe,
}

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.AddCommand(describeAccountCmd)
	describeCmd.AddCommand(describeContactCmd)
	describeCmd.AddCommand(describeOpportunityCmd)
	describeCmd.AddCommand(describeUserCmd)
}

func accountDescribe(cmd *cobra.Command, args []string) {
	dr, err := app.sc.AccountService.Describe(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err)
		os.Exit(1)
	}
	printDescription(dr)
}
func contactDescribe(cmd *cobra.Command, args []string) {
	dr, err := app.sc.ContactService.Describe(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err)
		os.Exit(1)
	}
	printDescription(dr)
}
func opportunityDescribe(cmd *cobra.Command, args []string) {
	dr, err := app.sc.OpportunityService.Describe(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err)
		os.Exit(1)
	}
	printDescription(dr)
}
func userDescribe(cmd *cobra.Command, args []string) {
	dr, err := app.sc.UserService.Describe(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err)
		os.Exit(1)
	}
	printDescription(dr)
}

func printDescription(dr *salesforce.DescribeResponse) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	blue := color.New(color.FgHiBlue)
	fmt.Println()
	blue.Println(dr.Label, "Fields")
	tblFields := table.New("Name", "Label", "Type", "Length", "Unique", "Updateable", "ID Lookup", "Relationship Name", "Reference To")
	tblFields.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for _, f := range dr.Fields {
		rels := ""
		if f.RelationshipName != "" {
			rels = strings.Join(f.ReferenceTo, "|")
		}
		tblFields.AddRow(f.Name, f.Label, f.Type, f.Length, f.Unique, f.Updateable, f.IDLookup, f.RelationshipName, rels)
	}
	tblFields.Print()
	fmt.Println()
	blue.Println(dr.Label, "Record Types")
	tblRecordTypes := table.New("ID", "Name", "Developer Name", "Available", "Active")
	tblRecordTypes.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for _, t := range dr.RecordTypeInfos {
		tblRecordTypes.AddRow(t.RecordTypeID, t.Name, t.DeveloperName, t.Available, t.Active)
	}
	tblRecordTypes.Print()
}
