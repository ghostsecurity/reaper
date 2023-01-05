package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ghostsecurity/reaper/backend/highlight"
	interceptor2 "github.com/ghostsecurity/reaper/backend/interceptor"
	"github.com/ghostsecurity/reaper/backend/log"
	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/proxy"
	"github.com/ghostsecurity/reaper/backend/settings"
	"github.com/ghostsecurity/reaper/backend/workspace"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
	"os"
	"sync"
	"time"
)

// App struct
type App struct {
	ctx          context.Context
	proxy        *proxy.Proxy
	logger       *log.Logger
	userSettings *settings.Provider
	workspaceMu  sync.RWMutex
	workspace    *workspace.Workspace
	readyChan    chan struct{}
	interceptor  *interceptor2.Interceptor
}

// New creates a new App application struct
func New(logger *log.Logger, settingsProvider *settings.Provider, ws *workspace.Workspace) *App {
	return &App{
		logger:       logger,
		readyChan:    make(chan struct{}),
		userSettings: settingsProvider,
		workspace:    ws,
	}
}

func (a *App) notifyUser(msg string, dialogType runtime.DialogType) {
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    dialogType,
		Title:   "Reaper",
		Message: msg,
	})
}

func (a *App) restartWithNewSettings(ctx context.Context) error {

	provider := a.userSettings

	a.logger.Infof("Creating proxy...")
	var err error
	a.proxy, err = proxy.New(provider, a.logger.WithPrefix("proxy"))
	if err != nil {
		return err
	}

	// tell the frontend about any custom user settings
	a.sendSettingsToFrontend()

	a.logger.Infof("Setting up proxy handlers...")
	a.proxy.OnRequest(func(request *http.Request, _ int64) (*http.Request, *http.Response) {
		if request.Host == provider.Get().ProxyHost {
			return request, a.handleLocalRequest(request)
		}
		return request, nil
	})
	a.proxy.OnRequest(func(request *http.Request, id int64) (*http.Request, *http.Response) {
		a.workspaceMu.RLock()
		defer a.workspaceMu.RUnlock()
		if !a.workspace.Scope.Includes(request) {
			return request, nil
		}
		runtime.EventsEmit(ctx, "OnHttpRequest", packaging.PackageHttpRequest(request, id))
		// update workspace tree
		runtime.EventsEmit(ctx, "OnTreeUpdate", a.workspace.UpdateTree(request).Structure())
		return request, nil
	})
	a.proxy.OnRequest(func(request *http.Request, id int64) (*http.Request, *http.Response) {
		a.workspaceMu.RLock()
		defer a.workspaceMu.RUnlock()
		if !a.workspace.Scope.Includes(request) {
			return request, nil
		}
		// TODO: check interception scope here as well
		return a.interceptor.Intercept(request, id)
	})
	a.proxy.OnResponse(func(response *http.Response, id int64) *http.Response {
		if response == nil {
			return nil
		}
		a.workspaceMu.RLock()
		defer a.workspaceMu.RUnlock()
		if !a.workspace.Scope.Includes(response.Request) {
			return response
		}
		runtime.EventsEmit(ctx, "OnHttpResponse", packaging.PackageHttpResponse(response, id))
		return response
	})

	a.logger.Infof("Starting proxy...")
	if err := a.proxy.Run(); err != nil {
		return err
	}

	a.logger.Infof("Proxy shut down cleanly.")
	return nil
}

// TODO: trigger this from the frontend
//func (a *App) setWorkspace(workspace *workspace.Workspace) {
//	a.workspaceMu.Lock()
//	defer a.workspaceMu.Unlock()
//	a.workspace = workspace
//}

func (a *App) isProxyRestartRequired(oldSettings settings.Settings, newSettings settings.Settings) bool {
	return oldSettings.ProxyPort != newSettings.ProxyPort ||
		oldSettings.LogLevel != newSettings.LogLevel ||
		string(oldSettings.CAKey) != string(newSettings.CAKey)
}

func (a *App) OnDomReady(ctx context.Context) {
	a.readyChan <- struct{}{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {

	runtime.EventsOn(ctx, "OnAppReady", func(_ ...interface{}) {
		a.logger.Infof("Frontend App component reported in!")
		a.sendSettingsToFrontend()
	})

	a.ctx = ctx

	runtime.EventsOn(ctx, "OnExportCA", func(_ ...interface{}) {
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

	a.interceptor = interceptor2.New(a.logger.WithPrefix("interceptor"), func(req *http.Request, id int64) {
		runtime.EventsEmit(a.ctx, "OnInterceptRequest", packaging.PackageHttpRequest(req, id))
	})

	runtime.EventsOn(a.ctx, "OnInterceptRequestModified", func(args ...interface{}) {

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
		a.interceptor.HandleCallback(final, modified.ID, nil)

	})

	runtime.EventsOn(a.ctx, "OnInterceptRequestDropped", func(args ...interface{}) {

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

		a.interceptor.HandleCallback(final, modified.ID, a.createReaperMessageResponse(final, "Request dropped."))
	})

	runtime.EventsOn(ctx, "OnInterceptionEnabledChange", func(args ...interface{}) {
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

	runtime.EventsOn(ctx, "OnHighlightRequest", func(args ...interface{}) {
		if len(args) != 1 {
			a.logger.Errorf("Invalid number of arguments for OnHighlightRequest")
			return
		}
		raw, ok := args[0].(string)
		if !ok {
			a.logger.Errorf("Invalid argument type for OnHighlightRequest")
			return
		}
		go func() {
			highlighted := highlight.HTTP(raw)
			if highlighted != "" {
				runtime.EventsEmit(ctx, "OnHighlightResponse", highlighted, raw)
			}
		}()
	})

	settingsUpdateChan := make(chan struct{}, 10)

	// ask the frontend to tell us when the user changes any settings
	a.logger.Infof("Subscribing to frontend events...")
	runtime.EventsOn(ctx, "OnSettingsSave", func(data ...interface{}) {

		a.logger.Debugf("Received settings update from frontend: %#v", data)

		if len(data) != 1 {
			a.logger.Errorf("OnSettingsSave: Expected 1 argument, got %d", len(data))
			return
		}

		raw, err := json.Marshal(data[0])
		if err != nil {
			a.notifyUser(fmt.Sprintf("Failed to save settings (marshal failed): %s", err), runtime.ErrorDialog)
			return
		}

		var newSettings settings.Settings
		if err := json.Unmarshal(raw, &newSettings); err != nil {
			a.notifyUser(fmt.Sprintf("Failed to save settings (unmarshal failed): %s", err), runtime.ErrorDialog)
			return
		}

		if err := settings.Save(&newSettings); err != nil {
			a.notifyUser(fmt.Sprintf("Failed to save settings (save failed): %s", err), runtime.ErrorDialog)
			return
		}

		oldSettings := a.userSettings.Get()

		if err := a.userSettings.Modify(func(s *settings.Settings) {
			*s = newSettings
		}); err != nil {
			a.notifyUser(fmt.Sprintf("Failed to update settings for current session: %s", err), runtime.ErrorDialog)
			return
		}

		if a.isProxyRestartRequired(oldSettings, newSettings) {
			a.logger.Infof("Proxy settings have changed, restarting proxy...")
			settingsUpdateChan <- struct{}{}
			_ = a.proxy.Close()
		} else {
			a.logger.Infof("Settings change does not require a proxy restart.")
			// tell the frontend about any custom user settings
			a.sendSettingsToFrontend()
		}
	})

	a.logger.Infof("Starting main run loop on separate routine...")
	go a.runLoop(ctx, settingsUpdateChan)

	a.logger.Infof("Startup complete!")
}

func (a *App) sendSettingsToFrontend() {
	a.logger.Infof("Sending settings to frontend...")
	runtime.EventsEmit(a.ctx, "OnSettingsLoad", a.userSettings.Get())
}

func (a *App) runLoop(ctx context.Context, settingsUpdateChan chan struct{}) {
	select {
	case <-ctx.Done():
		a.logger.Warnf("Context cancelled, shutting down...")
		return
	case <-a.readyChan:
		a.logger.Infof("Frontend is ready, starting proxy...")
	}
	go func() {
		defer a.logger.Info("Reload watcher routine exited cleanly")
		for range a.readyChan {
			a.logger.Warning("Frontend was reloaded!")
			a.sendSettingsToFrontend()
		}
	}()
	defer close(a.readyChan)
	for {
		if err := a.restartWithNewSettings(ctx); err != nil {
			a.notifyUser(fmt.Sprintf("Error starting proxy: %s - will retry shortly...", err), runtime.ErrorDialog)
			time.Sleep(5 * time.Second)
			continue
		}
		select {
		case <-ctx.Done():
			a.logger.Warnf("Context cancelled, shutting down...")
			return
		case <-settingsUpdateChan:
			a.logger.Infof("User settings updated, restarting proxy...")
		default:
			a.logger.Infof("Run loop is complete, exiting...")
			return
		}
	}
}

// Shutdown is called when the app is shutting down - we can cleanly stop the proxy here
func (a *App) Shutdown(_ context.Context) {
	a.logger.Infof("Shutting down...")
	if a.proxy != nil {
		a.logger.Infof("Shutting down proxy...")
		_ = a.proxy.Close()
	}
}
