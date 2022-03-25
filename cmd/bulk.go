package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/darrenparkinson/sfcli/pkg/salesforce"
	"github.com/spf13/cobra"
)

// used for both bulk insert and bulk upsert
var file string
var sobject string
var crlfLineEnding bool
var refreshTimer int64

var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Bulk API V2 Commands",
}

func init() {
	rootCmd.AddCommand(bulkCmd)
}

func (app *App) continuouslyUpdateStatusOrExit(id, initialState string) {
	if refreshTimer > 0 && initialState != "JobComplete" {
		for {
			time.Sleep(time.Duration(refreshTimer) * time.Second)
			bs, err := app.sc.BulkService.GetJob(context.Background(), salesforce.BulkTypeIngest, id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing CLI: Problem getting job status for ingest job: %s\n", err)
				os.Exit(1)
			}
			printJobStatus(bs)
			if bs.State == "JobComplete" {
				break
			}
		}
	}
}
