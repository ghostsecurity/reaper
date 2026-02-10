package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ghostsecurity/reaper/internal/daemon"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show recent proxy log entries",
	RunE:  runLogs,
}

var logsN int

func init() {
	logsCmd.Flags().IntVarP(&logsN, "number", "n", 50, "Number of entries to show")
	rootCmd.AddCommand(logsCmd)
}

func runLogs(cmd *cobra.Command, args []string) error {
	dataDir, err := daemon.DataDir()
	if err != nil {
		return err
	}

	params, _ := json.Marshal(daemon.LogsParams{Limit: logsN})
	client := daemon.NewClient(dataDir)
	resp, err := client.Send(daemon.Request{Command: "logs", Params: params})
	if err != nil {
		return fmt.Errorf("no running daemon found: %w", err)
	}
	if !resp.OK {
		return fmt.Errorf("%s", resp.Error)
	}

	var entries []entryRow
	if err := json.Unmarshal(resp.Data, &entries); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	printTable(entries)
	return nil
}
