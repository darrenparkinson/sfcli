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

func init() {
	bulkCmd.AddCommand(bulkStatusCmd)

	bulkStatusCmd.Flags().StringP("id", "i", "", "Job ID")
	viper.BindPFlag("bulkStatusID", bulkStatusCmd.Flags().Lookup("id"))

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
