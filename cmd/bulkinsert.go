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

	bulkInsertCmd.Flags().StringP("file", "f", "", "CSV File")
	viper.BindPFlag("file", bulkInsertCmd.Flags().Lookup("file"))

	bulkInsertCmd.Flags().StringP("object", "o", "", "Type of Object for Insert, e.g. Account, Contact, Opportunity")
	viper.BindPFlag("object", bulkInsertCmd.Flags().Lookup("object"))

}

func bulkInsert(cmd *cobra.Command, args []string) {

	filename := viper.GetString("file")
	if filename == "" {
		fmt.Fprintln(os.Stderr, "Error executing CLI: file is required")
		os.Exit(1)
	}
	object := viper.GetString("object")
	if object == "" {
		fmt.Fprintln(os.Stderr, "Error executing CLI: object type is required")
		os.Exit(1)
	}

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

}
