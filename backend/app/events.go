package app

import "github.com/wailsapp/wails/v2/pkg/runtime"

const (
	EventHttpRequest               = "HttpRequest"
	EventHttpResponse              = "HttpResponse"
	EventProxyStatus               = "ProxyStatusChange"
	EventCAExport                  = "CAExport"
	EventInterceptRequestModified  = "InterceptedRequestChange"
	EventInterceptRequestDropped   = "InterceptedRequestDrop"
	EventInterceptionEnabledChange = "InterceptionEnabledChange"
	EventInterceptedRequest        = "InterceptedRequest"
	EventTreeUpdate                = "TreeUpdate"
)

func (a *App) emitProxyStatus(status bool, addr, message string) {
	runtime.EventsEmit(a.ctx, EventProxyStatus, status, addr, message)
}
