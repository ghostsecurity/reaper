package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ghostsecurity/reaper/internal/daemon"
)

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Show full request and response for an entry",
	Args:  cobra.ExactArgs(1),
	RunE:  runGet,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) error {
	return fetchAndPrint(args[0], "get")
}

func fetchAndPrint(idStr, command string) error {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid entry ID: %s", idStr)
	}

	dataDir, err := daemon.DataDir()
	if err != nil {
		return err
	}

	params, _ := json.Marshal(daemon.GetParams{ID: id})
	client := daemon.NewClient(dataDir)
	resp, err := client.Send(daemon.Request{Command: command, Params: params})
	if err != nil {
		return fmt.Errorf("no running daemon found: %w", err)
	}
	if !resp.OK {
		return fmt.Errorf("%s", resp.Error)
	}

	var result struct {
		Command string    `json:"command"`
		Entry   entryFull `json:"entry"`
	}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	e := result.Entry
	switch command {
	case "get":
		printRawRequest(e)
		fmt.Println()
		printRawResponse(e)
	case "req":
		printRawRequest(e)
	case "res":
		printRawResponse(e)
	}

	return nil
}
