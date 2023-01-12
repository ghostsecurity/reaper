package app

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"

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
	ctx          context.Context
	proxy        *proxy.Proxy
	logger       *log.Logger
	userSettings *settings.Provider
	workspaceMu  sync.RWMutex
	workspace    *workspace.Workspace
	interceptor  *interceptor.Interceptor
	proxyMu      sync.Mutex
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

	a.interceptor = interceptor.New(a.logger.WithPrefix("interceptor"), func(req *http.Request, id int64) {
		if packaged, err := packaging.PackageHttpRequest(req, a.proxy.ID(), id); err != nil {
			a.logger.Errorf("Error packaging request: %s", err)
		} else {
			runtime.EventsEmit(a.ctx, EventInterceptedRequest, packaged)
		}
	})

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

	runtime.EventsOn(ctx, EventInterceptionEnabledChange, func(args ...interface{}) {
		if len(args) != 1 {
			a.logger.Errorf("Expected 1 argument, got %d", len(args))
			return
		}
		enabled, ok := args[0].(bool)
		if !ok {
			a.logger.Errorf("Expected bool, got %T", args[0])
			return
		}
		a.interceptor.SetEnabled(enabled)
	})

	a.logger.Infof("Startup complete!")
}

// Shutdown is called when the app is shutting down - we can cleanly stop the proxy here
func (a *App) Shutdown(_ context.Context) {
	a.logger.Infof("App closed, shutting down...")
	if err := a.stopProxy(); err != nil {
		a.logger.Errorf("Failed to stop proxy: %s", err)
	}
}

func (a *App) HighlightCode(code string) string {
	return highlight.HTTP(code)
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
