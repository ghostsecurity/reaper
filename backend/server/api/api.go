package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/ghostsecurity/reaper/backend/format"
	"github.com/ghostsecurity/reaper/backend/highlight"
	"github.com/ghostsecurity/reaper/backend/interceptor"
	"github.com/ghostsecurity/reaper/backend/log"
	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/proxy"
	"github.com/ghostsecurity/reaper/backend/settings"
	"github.com/ghostsecurity/reaper/backend/workflow"
	"github.com/ghostsecurity/reaper/backend/workflow/node"
	"github.com/ghostsecurity/reaper/backend/workspace"
)

type API struct {
	logger                *log.Logger
	userSettings          *settings.Provider
	ctx                   context.Context
	cancel                context.CancelFunc
	interceptor           *interceptor.Interceptor
	workspaceMu           sync.RWMutex
	proxyMu               sync.Mutex
	workflowMu            sync.Mutex
	eventTrigger          func(event string, args ...interface{}) error
	runningWorkflowID     uuid.UUID
	proxy                 *proxy.Proxy
	workspace             *workspace.Workspace
	workflowContextCancel context.CancelFunc
}

func New(ctx context.Context, logger *log.Logger, settingsProvider *settings.Provider, eventTrigger func(event string, args ...interface{}) error) *API {
	ctx, cancel := context.WithCancel(ctx)

	a := &API{
		ctx:          ctx,
		cancel:       cancel,
		logger:       logger,
		userSettings: settingsProvider,
		eventTrigger: eventTrigger,
	}

	a.interceptor = interceptor.New(
		logger.WithPrefix("interceptor"),
		workspace.Scope{},
		func(req *http.Request, id int64) {
			if packaged, err := packaging.PackageHttpRequest(req, a.proxy.ID(), id); err != nil {
				logger.Errorf("Error packaging request: %s", err)
			} else {
				_ = eventTrigger(EventInterceptedRequest, packaged)
			}
		},
		func(length int) {
			_ = eventTrigger(EventInterceptedRequestQueueChange, length)
		},
	)

	return a
}

func (a *API) BindingsOnly(x workflow.UpdateM, y node.OutputM) {

}

func (a *API) Test(input string) string {
	return input
}

// Close can be called to shut down the api
func (a *API) Close() {
	a.cancel()
	a.logger.Infof("App closed, shutting down...")
	if err := a.StopProxy(); err != nil {
		a.logger.Errorf("Failed to stop proxy: %s", err)
	}
}

func (a *API) notify(format string, args ...interface{}) {
	_ = a.eventTrigger(EventNotifyUser, fmt.Sprintf(format, args...))
}

func (a *API) HighlightHTTP(code string) string {
	return highlight.HTTP(code)
}

func (a *API) HighlightBody(body, contentType string) string {
	return highlight.Body(body, contentType)
}

func (a *API) GenerateID() string {
	return uuid.New().String()
}

func (a *API) FormatCode(msg, contentType string) string {
	return format.Code(msg, contentType)
}

func (a *API) sendRequest(request packaging.HttpRequest) {
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

func (a *API) isProxyRestartRequired(oldSettings settings.Settings, newSettings settings.Settings) bool {
	return oldSettings.ProxyPort != newSettings.ProxyPort ||
		oldSettings.LogLevel != newSettings.LogLevel ||
		string(oldSettings.CAKey) != string(newSettings.CAKey)
}
