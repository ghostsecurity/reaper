package cli

import (
	"encoding/json"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/ghostsecurity/reaper/internal/daemon"
)

var tailCmd = &cobra.Command{
	Use:          "tail",
	Short:        "Stream new proxy log entries in real-time",
	SilenceUsage: true,
	RunE:         runTail,
}

func init() {
	rootCmd.AddCommand(tailCmd)
}

func runTail(cmd *cobra.Command, args []string) error {
	dataDir, err := daemon.DataDir()
	if err != nil {
		return err
	}

	client := daemon.NewClient(dataDir)

	// Establish high-water mark from existing entries
	params, _ := json.Marshal(daemon.LogsParams{Limit: 1})
	resp, err := client.Send(daemon.Request{Command: "logs", Params: params})
	if err != nil {
		return fmt.Errorf("no running daemon found: %w", err)
	}
	if !resp.OK {
		return fmt.Errorf("%s", resp.Error)
	}

	var lastID int64
	var existing []entryRow
	if err := json.Unmarshal(resp.Data, &existing); err == nil && len(existing) > 0 {
		lastID = existing[0].ID // logs returns DESC order, so first is highest
	}

	ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Println("tailing proxy logs... (ctrl+c to stop)")

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			tailParams, _ := json.Marshal(daemon.TailParams{AfterID: lastID})
			resp, err := client.Send(daemon.Request{Command: "tail", Params: tailParams})
			if err != nil {
				return fmt.Errorf("daemon connection lost")
			}
			if !resp.OK {
				return fmt.Errorf("%s", resp.Error)
			}

			var entries []entryRow
			if err := json.Unmarshal(resp.Data, &entries); err != nil {
				continue
			}

			for _, e := range entries {
				ts := e.Timestamp.Local().Format("15:04:05")
				url := fmt.Sprintf("%s://%s%s", e.Scheme, e.Host, e.Path)
				fmt.Printf("%s %s %s %d %dms\n", ts, e.Method, url, e.StatusCode, e.DurationMs)
				if e.ID > lastID {
					lastID = e.ID
				}
			}
		}
	}
}
