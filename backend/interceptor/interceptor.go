package interceptor

import (
	"net/http"
	"sync"

	"github.com/ghostsecurity/reaper/backend/log"
	"github.com/ghostsecurity/reaper/backend/workspace"
)

type Interceptor struct {
	sync.Mutex         // protects the 'enabled' field
	interceptStartFunc func(req *http.Request, id int64)
	queueChangeFunc    func(len int)
	queue              queue
	logger             *log.Logger
	flightMu           sync.Mutex // ensures only one message can be in flight at a time
	scope              workspace.Scope
}

type InterceptionResponse struct {
	Req  *http.Request
	Resp *http.Response
}

type queuedRequest struct {
	id  int64
	req *http.Request
	c   chan *InterceptionResponse
}

func New(logger *log.Logger, scope workspace.Scope, start func(*http.Request, int64), queueChange func(int)) *Interceptor {
	return &Interceptor{
		interceptStartFunc: start,
		queueChangeFunc:    queueChange,
		logger:             logger,
		scope:              scope,
	}
}

func (i *Interceptor) SetScope(scope workspace.Scope) {
	i.Lock()
	defer i.Unlock()
	i.scope = scope
}

func (i *Interceptor) isInScope(req *http.Request) bool {
	i.Lock()
	defer i.Unlock()
	return len(i.scope.Include) > 0 && i.scope.Includes(req)
}

func (i *Interceptor) Intercept(req *http.Request, id int64) (*http.Request, *http.Response) {
	if !i.isInScope(req) {
		return req, nil
	}
	resultChan := i.queue.Add(req, id)
	if i.queueChangeFunc != nil {
		i.queueChangeFunc(i.queue.Len())
	}

	i.flightMu.Lock()
	defer i.flightMu.Unlock()
	if i.interceptStartFunc != nil {
		i.interceptStartFunc(req, id)
	}

	i.logger.Infof("Waiting for intercepted request %d to be processed...", id)
	interception := <-resultChan
	if interception == nil {
		return req, nil
	}
	return interception.Req, interception.Resp
}

func (i *Interceptor) HandleCallback(req *http.Request, id int64, resp *http.Response) {
	i.queue.Consume(req, id, resp)
	if i.queueChangeFunc != nil {
		i.queueChangeFunc(i.queue.Len())
	}
}

func (i *Interceptor) Flush() {
	i.Lock()
	defer i.Unlock()
	i.queue.Flush()
	if i.queueChangeFunc != nil {
		i.queueChangeFunc(i.queue.Len())
	}
}
