package interceptor

import (
	"net/http"
	"sync"

	"github.com/ghostsecurity/reaper/backend/log"
)

type Interceptor struct {
	sync.Mutex         // protects the 'enabled' field
	enabled            bool
	interceptStartFunc func(req *http.Request, id int64)
	queue              queue
	logger             *log.Logger
	flightMu           sync.Mutex // ensures only one message can be in flight at a time
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

func New(logger *log.Logger, start func(req *http.Request, id int64)) *Interceptor {
	return &Interceptor{
		enabled:            false,
		interceptStartFunc: start,
		logger:             logger,
	}
}

// access .enabled via mutex for safety
func (i *Interceptor) isEnabled() bool {
	i.Lock()
	defer i.Unlock()
	return i.enabled
}

func (i *Interceptor) Intercept(req *http.Request, id int64) (*http.Request, *http.Response) {
	if !i.isEnabled() {
		return req, nil
	}
	resultChan := i.queue.Add(req, id)

	i.flightMu.Lock()
	defer i.flightMu.Unlock()
	if i.interceptStartFunc != nil {
		if !i.isEnabled() {
			return req, nil
		}
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
}

func (i *Interceptor) SetEnabled(enabled bool) {
	i.Lock()
	defer i.Unlock()
	i.enabled = enabled
	if !enabled {
		i.queue.Flush()
	}
}
