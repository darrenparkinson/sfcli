package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/darrenparkinson/sfcli/pkg/salesforce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bulkInsertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Bulk Insert a CSV File",
	Run:   bulkInsert,
}

func init() {
	bulkCmd.AddCommand(bulkInsertCmd)

	bulkInsertCmd.Flags().StringVarP(&file, "file", "f", "", "CSV File")
	viper.BindPFlag("file", bulkInsertCmd.Flags().Lookup("file"))

	bulkInsertCmd.Flags().StringVarP(&sobject, "sobject", "s", "", "Type of Object for Insert, e.g. Account, Contact, Opportunity")
	viper.BindPFlag("sobject", bulkInsertCmd.Flags().Lookup("sobject"))

	bulkInsertCmd.Flags().BoolVarP(&crlfLineEnding, "crlf", "c", false, "Specify CRLF Line Ending (default is LF)")
	viper.BindPFlag("crlf", bulkInsertCmd.Flags().Lookup("crlf"))

	bulkInsertCmd.Flags().Int64VarP(&refreshTimer, "refresh", "r", 0, "Refresh timer in seconds, when set will check for status updates")
	viper.BindPFlag("refresh", bulkInsertCmd.Flags().Lookup("refresh"))
}

func bulkInsert(cmd *cobra.Command, args []string) {

	filename := viper.GetString("file")
	if filename == "" {
		fmt.Fprintln(os.Stderr, "Error executing CLI: file is required")
		os.Exit(1)
	}
	object := viper.GetString("sobject")
	if object == "" {
		fmt.Fprintln(os.Stderr, "Error executing CLI: object type is required")
		os.Exit(1)
	}

	crlf := viper.GetBool("crlf")

	// check file exists
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// create a job
	br := salesforce.BulkRequest{
		Object:      object,
		ContentType: "CSV", // TODO: Make this a parameter?
		Operation:   "insert",
	}
	if crlf {
		br.LineEnding = "CRLF"
	}
	job, err := app.sc.BulkService.CreateJob(context.Background(), br)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Job Created for insert: %s (%s)\n", job.ID, job.State)

	// upload the csv
	err = app.sc.BulkService.UploadCSV(context.Background(), job.ID, file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("File content uploaded, starting job...", job.ID)

	// begin the job
	res, err := app.sc.BulkService.ProcessJob(context.Background(), salesforce.BulkTypeIngest, job.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Job: %s; Status: %s\n", res.ID, res.State)

	// check for status updates if a refresh timer is specified
	app.continuouslyUpdateStatusOrExit(job.ID, res.State)

}
