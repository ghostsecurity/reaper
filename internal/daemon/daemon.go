package daemon

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/ghostsecurity/reaper/internal/proxy"
	"github.com/ghostsecurity/reaper/internal/storage"
	"github.com/ghostsecurity/reaper/version"
)

type Config struct {
	Domains []string
	Hosts   []string
	Port    int
	Daemon  bool
}

func DataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}

	dir := filepath.Join(home, ".ghost", "reaper")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("creating data directory: %w", err)
	}

	return dir, nil
}

func Run(cfg Config) error {
	dataDir, err := DataDir()
	if err != nil {
		return err
	}

	// Check for existing daemon
	sockPath := filepath.Join(dataDir, "reaper.sock")
	if _, err := os.Stat(sockPath); err == nil {
		return fmt.Errorf("daemon already running (socket exists at %s). Use 'reaper shutdown' first", sockPath)
	}

	// Generate in-memory CA
	ca, err := proxy.GenerateCA()
	if err != nil {
		return fmt.Errorf("generating CA: %w", err)
	}

	// Init storage
	dbPath := filepath.Join(dataDir, "reaper.db")
	store, err := storage.NewSQLiteStore(dbPath)
	if err != nil {
		return fmt.Errorf("opening storage: %w", err)
	}
	defer store.Close()

	// Create proxy
	scope := proxy.NewScope(cfg.Domains, cfg.Hosts)
	p := &proxy.Proxy{
		Scope: scope,
		Store: store,
		CA:    ca,
	}
	if !cfg.Daemon {
		p.OnEvent = printEvent
	}

	// Start IPC server
	shutdown := make(chan struct{})
	ipcServer, err := NewIPCServer(dataDir, store, shutdown)
	if err != nil {
		return fmt.Errorf("starting IPC server: %w", err)
	}
	defer ipcServer.Close()
	go ipcServer.Serve()

	// Write PID file
	pidPath := filepath.Join(dataDir, "reaper.pid")
	_ = os.WriteFile(pidPath, []byte(strconv.Itoa(os.Getpid())), 0600)
	defer os.Remove(pidPath)
	defer os.Remove(sockPath)

	// Start HTTP proxy server
	addr := fmt.Sprintf(":%d", cfg.Port)
	server := &http.Server{ //nolint:gosec
		Addr:    addr,
		Handler: p,
	}

	// Print banner
	printBanner(cfg, dataDir)

	// Signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigCh:
		case <-shutdown:
		}
		fmt.Println("\nshutting down...")
		server.Close()
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("proxy server: %w", err)
	}

	return nil
}

func printEvent(e proxy.Event) {
	ts := time.Now().Format("15:04:05")
	tag := "="
	if e.Intercepted {
		tag = "⇄"
	}

	url := fmt.Sprintf("%s://%s%s", e.Scheme, e.Host, e.Path)
	fmt.Printf("%s %s %s %s %d %dms\n", ts, tag, e.Method, url, e.StatusCode, e.DurationMs)
}

func printBanner(cfg Config, dataDir string) {
	fmt.Printf("reaper %s\n", version.Version)
	fmt.Printf("proxy listening on :%d\n", cfg.Port)
	fmt.Printf("data directory: %s\n", dataDir)
	if len(cfg.Domains) > 0 {
		fmt.Printf("domains: %v\n", cfg.Domains)
	}
	if len(cfg.Hosts) > 0 {
		fmt.Printf("hosts: %v\n", cfg.Hosts)
	}
	fmt.Printf("started at %s\n\n", time.Now().Format(time.DateTime))
}
