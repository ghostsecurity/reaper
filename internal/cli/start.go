package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ghostsecurity/reaper/internal/daemon"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the proxy",
	RunE:  runStart,
}

var (
	startDomains  []string
	startHosts    []string
	startPort     int
	startDaemon   bool
	startInternal bool
)

func init() {
	startCmd.Flags().StringSliceVar(&startDomains, "domains", nil, "Domain suffixes to intercept (e.g. example.com)")
	startCmd.Flags().StringSliceVar(&startHosts, "hosts", nil, "Exact hostnames to intercept (e.g. api.example.com)")
	startCmd.Flags().IntVar(&startPort, "port", 8443, "Proxy listen port")
	startCmd.Flags().BoolVarP(&startDaemon, "daemon", "d", false, "Run as background daemon")
	startCmd.Flags().BoolVar(&startInternal, "internal", false, "Internal flag for daemon child process")
	_ = startCmd.Flags().MarkHidden("internal")

	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) error {
	if len(startDomains) == 0 && len(startHosts) == 0 {
		return fmt.Errorf("at least one --domains or --hosts flag is required")
	}

	cfg := daemon.Config{
		Domains: startDomains,
		Hosts:   startHosts,
		Port:    startPort,
		Daemon:  startInternal,
	}

	if startDaemon && !startInternal {
		return daemonize(cfg)
	}

	return daemon.Run(cfg)
}

func daemonize(cfg daemon.Config) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolving executable: %w", err)
	}

	daemonArgs := []string{"start", "--internal", "--port", fmt.Sprintf("%d", cfg.Port)}
	if len(cfg.Domains) > 0 {
		daemonArgs = append(daemonArgs, "--domains", strings.Join(cfg.Domains, ","))
	}
	if len(cfg.Hosts) > 0 {
		daemonArgs = append(daemonArgs, "--hosts", strings.Join(cfg.Hosts, ","))
	}

	proc, err := os.StartProcess(exe, append([]string{exe}, daemonArgs...), &os.ProcAttr{
		Dir:   "/",
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Sys:   daemonSysProcAttr(),
	})
	if err != nil {
		return fmt.Errorf("starting daemon: %w", err)
	}

	// Detach from child — let it run independently
	if err := proc.Release(); err != nil {
		return fmt.Errorf("releasing daemon process: %w", err)
	}

	// Verify daemon started by checking socket
	dataDir, err := daemon.DataDir()
	if err != nil {
		return err
	}

	if err := daemon.WaitForSocket(dataDir); err != nil {
		return fmt.Errorf("daemon failed to start: %w", err)
	}

	fmt.Println("reaper daemon started")
	return nil
}
