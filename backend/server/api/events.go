package api

const (
	EventHttpRequest  = "HttpRequest"
	EventHttpResponse = "HttpResponse"

	// EventProxyStatus is emitted when the proxy status changes, also includes the proxy address and a message
	EventProxyStatus                   = "ProxyStatusChange"
	EventInterceptedRequest            = "InterceptedRequest"
	EventTreeUpdate                    = "TreeUpdate"
	EventInterceptedRequestQueueChange = "InterceptedRequestQueueChange"
	EventWorkflowStarted               = "WorkflowStarted"
	EventWorkflowFinished              = "WorkflowFinished"
	EventWorkflowUpdate                = "WorkflowUpdated"
	EventWorkflowOutput                = "WorkflowOutput"
	EventNotifyUser                    = "NotifyUser"
)
