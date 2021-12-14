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

var bulkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the last 1000 bulk jobs",
	Run:   bulkList,
}

func init() {
	bulkCmd.AddCommand(bulkListCmd)

	bulkListCmd.Flags().BoolP("ingest", "i", false, "List Ingest jobs")
	bulkListCmd.Flags().BoolP("query", "q", false, "List Query jobs")
	viper.BindPFlag("ingest", bulkListCmd.Flags().Lookup("ingest"))
	viper.BindPFlag("query", bulkListCmd.Flags().Lookup("query"))

}

func bulkList(cmd *cobra.Command, args []string) {

	ingest, query := viper.GetBool("ingest"), viper.GetBool("query")
	if !ingest && !query {
		fmt.Println("No job type defined, defaulting to ingest jobs")
		ingest = true
	}
	if ingest {
		lr, err := app.sc.BulkService.ListJobs(context.Background(), salesforce.BulkTypeIngest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing CLI: Problem listing bulk ingest jobs: %s\n", err)
			os.Exit(1)
		}
		printBulkJobs(lr.Records, "INGEST JOBS")
	}
	if query {
		lr, err := app.sc.BulkService.ListJobs(context.Background(), salesforce.BulkTypeQuery)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing CLI: Problem listing bulk query jobs: %s\n", err)
			os.Exit(1)
		}
		printBulkJobs(lr.Records, "QUERY JOBS")
	}

}

func printBulkJobs(jobs []salesforce.JobInfo, title string) {
	// TODO: Get Submitted By Name
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	blue := color.New(color.FgHiBlue)
	fmt.Println()
	blue.Println(title)
	tblIngestJobs := table.New("Job ID", "Status", "Job Type", "Operation", "Object", "Submitted By ID", "Start Time")
	tblIngestJobs.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, r := range jobs {
		tblIngestJobs.AddRow(r.ID, r.State, r.JobType, r.Operation, r.Object, r.CreatedByID, r.CreatedDate)
	}
	tblIngestJobs.Print()
}
