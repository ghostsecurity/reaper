package browser

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/ghostsecurity/reaper/internal/tools/proxy"
)

type Browser struct {
	sync.Mutex
	cmd         *exec.Cmd
	vncCmd      *exec.Cmd
	proxy       *proxy.Proxy
	displayNum  int
	chromePid   int
	currentURL  string
	shutdownCtx context.Context
	shutdown    context.CancelFunc
}

func NewBrowser(proxy *proxy.Proxy) *Browser {
	ctx, cancel := context.WithCancel(context.Background())
	return &Browser{
		proxy:       proxy,
		displayNum:  99, // Use a high display number to avoid conflicts
		shutdownCtx: ctx,
		shutdown:    cancel,
	}
}

func (b *Browser) Start() error {
	b.Lock()
	defer b.Unlock()

	// Start Xvfb
	display := fmt.Sprintf(":%d", b.displayNum)
	xvfbCmd := exec.Command("Xvfb", display, "-screen", "0", "1280x800x24")
	if err := xvfbCmd.Start(); err != nil {
		return fmt.Errorf("failed to start Xvfb: %v", err)
	}

	// Start VNC server
	vncCmd := exec.Command("x11vnc", "-display", display, "-forever", "-shared", "-rfbport", "5900")
	if err := vncCmd.Start(); err != nil {
		xvfbCmd.Process.Kill()
		return fmt.Errorf("failed to start VNC server: %v", err)
	}
	b.vncCmd = vncCmd

	// Get proxy address
	proxyHost := os.Getenv("HOST")
	proxyPort := os.Getenv("PROXY_PORT")
	proxyAddr := fmt.Sprintf("http://%s:%s", proxyHost, proxyPort)

	// Start Chrome with proxy settings
	userDataDir := filepath.Join(os.TempDir(), "reaper-chrome-profile")
	chromeArgs := []string{
		"--no-sandbox",
		"--disable-gpu",
		"--headless=new",
		"--remote-debugging-port=9222",
		fmt.Sprintf("--user-data-dir=%s", userDataDir),
		fmt.Sprintf("--proxy-server=%s", proxyAddr),
		"about:blank",
	}

	cmd := exec.Command("google-chrome", chromeArgs...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("DISPLAY=%s", display))

	if err := cmd.Start(); err != nil {
		b.cleanup()
		return fmt.Errorf("failed to start Chrome: %v", err)
	}

	b.cmd = cmd
	b.chromePid = cmd.Process.Pid

	// Monitor processes
	go b.monitorProcesses()

	return nil
}

func (b *Browser) Navigate(urlStr string) error {
	b.Lock()
	defer b.Unlock()

	if b.cmd == nil {
		return fmt.Errorf("browser not started")
	}

	// Validate URL
	_, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	// Use Chrome DevTools Protocol to navigate
	// For simplicity, we'll use a direct command here
	script := fmt.Sprintf("window.location.href = '%s'", urlStr)
	evalCmd := exec.Command("google-chrome", "--remote-debugging-port=9222", "--headless=new",
		"--disable-gpu", "--repl", fmt.Sprintf("--eval=%s", script))

	if err := evalCmd.Run(); err != nil {
		return fmt.Errorf("failed to navigate: %v", err)
	}

	b.currentURL = urlStr
	return nil
}

func (b *Browser) Reload() error {
	b.Lock()
	defer b.Unlock()

	if b.cmd == nil {
		return fmt.Errorf("browser not started")
	}

	script := "window.location.reload()"
	evalCmd := exec.Command("google-chrome", "--remote-debugging-port=9222", "--headless=new",
		"--disable-gpu", "--repl", fmt.Sprintf("--eval=%s", script))

	if err := evalCmd.Run(); err != nil {
		return fmt.Errorf("failed to reload: %v", err)
	}

	return nil
}

func (b *Browser) Stop() error {
	b.Lock()
	defer b.Unlock()

	if b.shutdown != nil {
		b.shutdown()
	}

	return b.cleanup()
}

func (b *Browser) cleanup() error {
	var errors []error

	if b.cmd != nil && b.cmd.Process != nil {
		if err := b.cmd.Process.Kill(); err != nil {
			errors = append(errors, fmt.Errorf("failed to kill Chrome: %v", err))
		}
	}

	if b.vncCmd != nil && b.vncCmd.Process != nil {
		if err := b.vncCmd.Process.Kill(); err != nil {
			errors = append(errors, fmt.Errorf("failed to kill VNC server: %v", err))
		}
	}

	// Kill Xvfb
	display := fmt.Sprintf(":%d", b.displayNum)
	if err := exec.Command("pkill", "-f", fmt.Sprintf("Xvfb %s", display)).Run(); err != nil {
		errors = append(errors, fmt.Errorf("failed to kill Xvfb: %v", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %v", errors)
	}

	return nil
}

func (b *Browser) monitorProcesses() {
	select {
	case <-b.shutdownCtx.Done():
		slog.Info("Browser shutdown requested")
		b.cleanup()
	}
}
