package api

import (
	"github.com/ghostsecurity/reaper/backend/packaging"
)

func (a *API) ModifyInterceptedRequest(modified packaging.HttpRequest) {
	final, err := packaging.UnpackageHttpRequest(&modified)
	if err != nil {
		a.logger.Errorf("failed to unpack incoming request from frontend: %s", err)
		return
	}
	a.interceptor.HandleCallback(final, modified.LocalID, nil)
}

func (a *API) DropInterceptedRequest(modified packaging.HttpRequest) {
	final, err := packaging.UnpackageHttpRequest(&modified)
	if err != nil {
		a.logger.Errorf("failed to unpack incoming request from frontend: %s", err)
		return
	}
	a.interceptor.HandleCallback(final, modified.LocalID, a.createReaperMessageResponse(final, "Request dropped."))
}

func (a *API) SendRequest(req packaging.HttpRequest) {
	a.sendRequest(req)
}
