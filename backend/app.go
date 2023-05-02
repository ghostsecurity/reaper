package backend

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/ghostsecurity/reaper/backend/highlight"
	"github.com/ghostsecurity/reaper/backend/interceptor"
	"github.com/ghostsecurity/reaper/backend/log"
	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/proxy"
	"github.com/ghostsecurity/reaper/backend/settings"
	"github.com/ghostsecurity/reaper/backend/workspace"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx                   context.Context
	proxy                 *proxy.Proxy
	logger                *log.Logger
	userSettings          *settings.Provider
	workspaceMu           sync.RWMutex
	workspace             *workspace.Workspace
	interceptor           *interceptor.Interceptor
	proxyMu               sync.Mutex
	workflowContextCancel context.CancelFunc
	runningWorkflowID     uuid.UUID
	workflowMu            sync.Mutex
}

// New creates a new App application struct
func New(logger *log.Logger, settingsProvider *settings.Provider) *App {
	return &App{
		logger:       logger,
		userSettings: settingsProvider,
	}
}

func (a *App) notifyUser(msg string, dialogType runtime.DialogType) {
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    dialogType,
		Title:   "Reaper",
		Message: msg,
	})
}

func (a *App) isProxyRestartRequired(oldSettings settings.Settings, newSettings settings.Settings) bool {
	return oldSettings.ProxyPort != newSettings.ProxyPort ||
		oldSettings.LogLevel != newSettings.LogLevel ||
		string(oldSettings.CAKey) != string(newSettings.CAKey)
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {

	a.ctx = ctx

	runtime.EventsOn(ctx, EventCAExport, func(_ ...interface{}) {
		a.logger.Infof("Exporting CA...")
		path, err := runtime.SaveFileDialog(ctx, runtime.SaveDialogOptions{
			DefaultFilename:      "ghost-ca.crt",
			Title:                "Select CA file location",
			CanCreateDirectories: true,
		})
		if err != nil {
			a.logger.Errorf("Error opening save dialog: %s", err)
			return
		}
		if path == "" {
			a.logger.Infof("No path selected, skipping export")
			return
		}
		a.logger.Infof("Saving CA to %s...", path)
		if err := os.WriteFile(path, a.userSettings.Get().CACert, 0644); err != nil {
			a.logger.Errorf("Error saving CA cert: %s", err)
			return
		}
	})

	a.interceptor = interceptor.New(
		a.logger.WithPrefix("interceptor"),
		workspace.Scope{},
		func(req *http.Request, id int64) {
			if packaged, err := packaging.PackageHttpRequest(req, a.proxy.ID(), id); err != nil {
				a.logger.Errorf("Error packaging request: %s", err)
			} else {
				runtime.EventsEmit(a.ctx, EventInterceptedRequest, packaged)
			}
		},
		func(length int) {
			runtime.EventsEmit(a.ctx, EventInterceptedRequestQueueChange, length)
		},
	)

	runtime.EventsOn(a.ctx, EventInterceptRequestModified, func(args ...interface{}) {

		if len(args) != 1 {
			a.logger.Errorf("OnInterceptRequestModified: Expected 1 argument, got %d", len(args))
			return
		}

		raw, err := json.Marshal(args[0])
		if err != nil {
			a.logger.Errorf("failed to marshal incoming request from frontend: %s", err)
			return
		}

		var modified packaging.HttpRequest
		if err := json.Unmarshal(raw, &modified); err != nil {
			a.logger.Errorf("failed to unmarshal incoming request from frontend: %s", err)
			return
		}

		final, err := packaging.UnpackageHttpRequest(&modified)
		if err != nil {
			a.logger.Errorf("failed to unpack incoming request from frontend: %s", err)
			return
		}

		// TODO: log return value
		a.interceptor.HandleCallback(final, modified.LocalID, nil)

	})

	runtime.EventsOn(a.ctx, EventInterceptRequestDropped, func(args ...interface{}) {

		if len(args) != 1 {
			a.logger.Errorf("OnInterceptRequestDropped: Expected 1 argument, got %d", len(args))
			return
		}

		raw, err := json.Marshal(args[0])
		if err != nil {
			a.logger.Errorf("failed to marshal incoming request from frontend: %s", err)
			return
		}

		var modified packaging.HttpRequest
		if err := json.Unmarshal(raw, &modified); err != nil {
			a.logger.Errorf("failed to unmarshal incoming request from frontend: %s", err)
			return
		}

		final, err := packaging.UnpackageHttpRequest(&modified)
		if err != nil {
			a.logger.Errorf("failed to unpack incoming request from frontend: %s", err)
			return
		}

		a.interceptor.HandleCallback(final, modified.LocalID, a.createReaperMessageResponse(final, "Request dropped."))
	})

	runtime.EventsOn(ctx, EventSendRequest, func(args ...interface{}) {
		for _, arg := range args {
			if req, ok := arg.(packaging.HttpRequest); ok {
				a.sendRequest(req)
			} else if m, ok := arg.(map[string]interface{}); ok {
				data, err := json.Marshal(m)
				if err != nil {
					a.logger.Errorf("failed to marshal request: %s", err)
					continue
				}
				var req packaging.HttpRequest
				if err := json.Unmarshal(data, &req); err != nil {
					a.logger.Errorf("failed to unmarshal request: %s", err)
					continue
				}
				a.sendRequest(req)
			} else {
				a.logger.Errorf("Expected HttpRequest, got %T: %#v", arg, arg)
			}
		}
	})

	a.logger.Infof("Startup complete!")
}

func (a *App) sendRequest(request packaging.HttpRequest) {
	req, err := packaging.UnpackageHttpRequest(&request)
	if err != nil {
		a.logger.Errorf("failed to unpack request: %s", err)
		return
	}
	port := a.userSettings.Get().ProxyPort
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy: http.ProxyURL(
			&url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("127.0.0.1:%d", port),
			},
		),
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	proxyClient := http.Client{
		Timeout:   time.Minute,
		Transport: transport,
	}
	if _, err := proxyClient.Do(req); err != nil {
		a.logger.Errorf("failed to send request: %s", err)
		return
	}
}

// Shutdown is called when the app is shutting down - we can cleanly stop the proxy here
func (a *App) Shutdown(_ context.Context) {
	a.logger.Infof("App closed, shutting down...")
	if err := a.stopProxy(); err != nil {
		a.logger.Errorf("Failed to stop proxy: %s", err)
	}
}

func (a *App) HighlightHTTP(code string) string {
	return highlight.HTTP(code)
}

func (a *App) HighlightBody(body, contentType string) string {
	return highlight.Body(body, contentType)
}

func (a *App) GenerateID() string {
	return uuid.New().String()
}

func (a *App) Confirm(title, msg string) bool {
	result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         title,
		Message:       msg,
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "No",
	})
	if err != nil {
		a.logger.Errorf("Error showing confirmation: %s", err)
	}
	return err == nil && result == "Yes"
}

func (a *App) Notify(title, msg string) {
	a.message(title, msg, runtime.InfoDialog)
}

func (a *App) Warn(title, msg string) {
	a.message(title, msg, runtime.WarningDialog)
}

func (a *App) Error(title, msg string) {
	a.message(title, msg, runtime.ErrorDialog)
}

func (a *App) message(title, msg string, dialogType runtime.DialogType) {
	_, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    dialogType,
		Title:   title,
		Message: msg,
	})
	if err != nil {
		a.logger.Errorf("Error showing notification: %s", err)
	}
}
