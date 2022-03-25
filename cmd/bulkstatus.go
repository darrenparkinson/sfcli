package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/darrenparkinson/sfcli/pkg/salesforce"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bulkStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a specific job",
	Run:   bulkStatus,
}
var bulkSuccessResultsCmd = &cobra.Command{
	Use:   "success",
	Short: "Download the successful results",
	Run:   bulkSuccessResults,
}
var bulkErrorResultsCmd = &cobra.Command{
	Use:   "errors",
	Short: "Download the error results",
	Run:   bulkErrorResults,
}

func init() {
	bulkCmd.AddCommand(bulkStatusCmd)

	bulkStatusCmd.AddCommand(bulkSuccessResultsCmd)
	bulkStatusCmd.AddCommand(bulkErrorResultsCmd)

	bulkStatusCmd.Flags().StringP("id", "i", "", "Job ID")
	viper.BindPFlag("bulkStatusID", bulkStatusCmd.Flags().Lookup("id"))
	bulkStatusCmd.Flags().Int64VarP(&refreshTimer, "refresh", "r", 0, "Refresh timer in seconds, when set will check for status updates")
	viper.BindPFlag("refresh", bulkStatusCmd.Flags().Lookup("refresh"))

	bulkSuccessResultsCmd.Flags().StringP("id", "i", "", "Job ID")
	viper.BindPFlag("bulkSuccessID", bulkSuccessResultsCmd.Flags().Lookup("id"))

	bulkErrorResultsCmd.Flags().StringP("id", "i", "", "Job ID")
	viper.BindPFlag("bulkErrorID", bulkErrorResultsCmd.Flags().Lookup("id"))

}

func bulkSuccessResults(cmd *cobra.Command, args []string) {
	id := viper.GetString("bulkSuccessID")
	if id == "" {
		fmt.Fprintln(os.Stderr, "Error executing CLI: ID is required")
		os.Exit(1)
	}
	res, err := app.sc.BulkService.GetSuccessfulResults(context.Background(), salesforce.BulkTypeIngest, id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: Problem getting success results: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}
func bulkErrorResults(cmd *cobra.Command, args []string) {
	id := viper.GetString("bulkErrorID")
	if id == "" {
		fmt.Fprintln(os.Stderr, "Error executing CLI: ID is required")
		os.Exit(1)
	}
	res, err := app.sc.BulkService.GetFailedResults(context.Background(), salesforce.BulkTypeIngest, id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: Problem getting error results: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}

func bulkStatus(cmd *cobra.Command, args []string) {
	id := viper.GetString("bulkStatusID")

	if id == "" {
		fmt.Fprintln(os.Stderr, "Error executing CLI: ID is required")
		os.Exit(1)
	}
	bs, err := app.sc.BulkService.GetJob(context.Background(), salesforce.BulkTypeIngest, id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: Problem getting job status for ingest job: %s\n", err)
		os.Exit(1)
	}
	printJobStatus(bs)

	// check for status updates if a refresh timer is specified
	app.continuouslyUpdateStatusOrExit(id, bs.State)

}

func printJobStatus(job *salesforce.JobInfo) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	blue := color.New(color.FgHiBlue)
	fmt.Println()
	blue.Println("Job ID:", job.ID)
	tblJob := table.New("Field", "Value")
	tblJob.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tblJob.AddRow("ID", job.ID)
	tblJob.AddRow("Operation", job.Operation)
	tblJob.AddRow("Object", job.Object)
	tblJob.AddRow("CreatedById", job.CreatedByID)
	tblJob.AddRow("CreatedDate", job.CreatedDate)
	tblJob.AddRow("JobType", job.JobType)
	tblJob.AddRow("Status", job.State)
	tblJob.AddRow("RecordsProcessed", job.NumberRecordsProcessed)
	tblJob.AddRow("RecordsFailed", job.NumberRecordsFailed)
	tblJob.AddRow("Retries", job.Retries)
	tblJob.Print()
}
