package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/darrenparkinson/sfcli/pkg/salesforce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bulkUpsertCmd = &cobra.Command{
	Use:   "upsert",
	Short: "Bulk Upsert a CSV File",
	Run:   bulkUpsert,
}

func init() {
	bulkCmd.AddCommand(bulkUpsertCmd)

	bulkUpsertCmd.Flags().StringP("file", "f", "", "CSV File")
	viper.BindPFlag("file", bulkUpsertCmd.Flags().Lookup("file"))

	bulkUpsertCmd.Flags().StringP("object", "o", "", "Type of Object for Insert, e.g. Account, Contact, Opportunity")
	viper.BindPFlag("object", bulkUpsertCmd.Flags().Lookup("object"))

	bulkUpsertCmd.Flags().StringP("external", "e", "", "External ID Field")
	viper.BindPFlag("external", bulkUpsertCmd.Flags().Lookup("external"))

}

func bulkUpsert(cmd *cobra.Command, args []string) {

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
	external := viper.GetString("external")
	if external == "" {
		fmt.Fprintln(os.Stderr, "Error executing CLI: external id is required")
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
		Object:              object,
		ContentType:         "CSV", // TODO: Make this a parameter?
		Operation:           "upsert",
		ExternalIDFieldName: external,
	}
	job, err := app.sc.BulkService.CreateJob(context.Background(), br)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Job Created for upsert: %s (%s)\n", job.ID, job.State)

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
	fmt.Printf("Started: %s; Status: %s\n", res.ID, res.State)

}
