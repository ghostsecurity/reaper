package workflow

import (
	"fmt"

	"github.com/ghostsecurity/reaper/backend/packaging"
)

type TransmissionType uint8

const (
	TransmissionTypeUnknown TransmissionType = iota
	TransmissionTypeRequest
	TransmissionTypeRequestAndResponse
	TransmissionTypeAny
)

type Transmission interface {
	Type() TransmissionType
}

func (t TransmissionType) validate() error {
	switch t {
	case TransmissionTypeRequest:
	case TransmissionTypeRequestAndResponse:
	default:
		return fmt.Errorf("invalid transmission type: %d", t)
	}
	return nil
}

type RequestProvider interface {
	Request() packaging.HttpRequest
}

type ResponseProvider interface {
	Response() packaging.HttpResponse
}

type RequestAndResponseProvider interface {
	RequestProvider
	ResponseProvider
	ParameterSetProvider
}

type ListProvider interface {
	List() []string
}

type ParameterSetProvider interface {
	ParameterSet() map[string]string
}

type requestTransmission struct {
	request packaging.HttpRequest
}

func NewRequestTransmission(req packaging.HttpRequest) Transmission {
	return &requestTransmission{
		request: req,
	}
}

func (r *requestTransmission) Type() TransmissionType {
	return TransmissionTypeRequest
}

func (r *requestTransmission) Request() packaging.HttpRequest {
	return r.request
}

type requestAndResponseTransmission struct {
	request  packaging.HttpRequest
	response packaging.HttpResponse
	params   map[string]string
}

func NewRequestAndResponseTransmission(req packaging.HttpRequest, resp packaging.HttpResponse, params map[string]string) Transmission {
	return &requestAndResponseTransmission{
		request:  req,
		response: resp,
		params:   params,
	}
}

func (r *requestAndResponseTransmission) Type() TransmissionType {
	return TransmissionTypeRequestAndResponse
}

func (r *requestAndResponseTransmission) Request() packaging.HttpRequest {
	return r.request
}

func (r *requestAndResponseTransmission) Response() packaging.HttpResponse {
	return r.response
}

func (r *requestAndResponseTransmission) ParameterSet() map[string]string {
	return r.params
}
