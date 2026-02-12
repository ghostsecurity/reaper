package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ghostsecurity/reaper/internal/daemon"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all proxy log entries",
	RunE:  runClear,
}

func init() {
	rootCmd.AddCommand(clearCmd)
}

func runClear(cmd *cobra.Command, args []string) error {
	dataDir, err := daemon.DataDir()
	if err != nil {
		return err
	}

	client := daemon.NewClient(dataDir)
	resp, err := client.Send(daemon.Request{Command: "clear"})
	if err != nil {
		return fmt.Errorf("no running daemon found: %w", err)
	}
	if !resp.OK {
		return fmt.Errorf("clear failed: %s", resp.Error)
	}

	fmt.Println("all entries cleared")
	return nil
}
