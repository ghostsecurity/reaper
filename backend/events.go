package backend

import "github.com/wailsapp/wails/v2/pkg/runtime"

const (
	EventHttpRequest                   = "HttpRequest"
	EventHttpResponse                  = "HttpResponse"
	EventProxyStatus                   = "ProxyStatusChange"
	EventCAExport                      = "CAExport"
	EventInterceptRequestModified      = "InterceptedRequestChange"
	EventInterceptRequestDropped       = "InterceptedRequestDrop"
	EventInterceptedRequest            = "InterceptedRequest"
	EventTreeUpdate                    = "TreeUpdate"
	EventSendRequest                   = "SendRequest"
	EventInterceptedRequestQueueChange = "InterceptedRequestQueueChange"
	EventWorkflowStarted               = "WorkflowStarted"
	EventWorkflowFinished              = "WorkflowFinished"
)

func (a *App) emitProxyStatus(status bool, addr, message string) {
	runtime.EventsEmit(a.ctx, EventProxyStatus, status, addr, message)
}
