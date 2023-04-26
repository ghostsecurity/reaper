package transmission

import (
	"encoding/json"

	"github.com/ghostsecurity/reaper/backend/packaging"
)

type Request packaging.HttpRequest

func NewRequest(req packaging.HttpRequest) *Request {
	r := Request(req)
	return &r
}

func (r *Request) Type() Type {
	return NewType(TypeRequest, InternalTypeNone)
}

func (r *Request) Request() packaging.HttpRequest {
	if r == nil {
		return packaging.HttpRequest{}
	}
	return packaging.HttpRequest(*r)
}

func (r *Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Request())
}

func (r *Request) UnmarshalJSON(data []byte) error {
	var v packaging.HttpRequest
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*r = Request(v)
	return nil
}

type Response packaging.HttpResponse

func (r *Response) Type() Type {
	return NewType(TypeResponse, InternalTypeNone)
}

func (r *Response) Response() packaging.HttpResponse {
	if r == nil {
		return packaging.HttpResponse{}
	}
	return packaging.HttpResponse(*r)
}

func (r *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Response())
}

func (r *Response) UnmarshalJSON(data []byte) error {
	var v packaging.HttpResponse
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*r = Response(v)
	return nil
}

type RequestResponsePair struct {
	request  Request
	response Response
}

func NewRequestResponsePair(req packaging.HttpRequest, resp packaging.HttpResponse) *RequestResponsePair {
	return &RequestResponsePair{
		request:  Request(req),
		response: Response(resp),
	}
}

func (p *RequestResponsePair) Type() Type {
	return NewType(TypeRequest|TypeResponse, InternalTypeNone)
}

func (p *RequestResponsePair) Request() packaging.HttpRequest {
	return p.request.Request()
}

func (p *RequestResponsePair) Response() packaging.HttpResponse {
	return p.response.Response()
}

func (p *RequestResponsePair) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"request":  p.request,
		"response": p.response,
	})
}

func (p *RequestResponsePair) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v["request"] != nil {
		if err := json.Unmarshal(v["request"].([]byte), &p.request); err != nil {
			return err
		}
	}
	if v["response"] != nil {
		if err := json.Unmarshal(v["response"].([]byte), &p.response); err != nil {
			return err
		}
	}
	return nil
}

type RequestResponsePairWithMap struct {
	request  Request
	response Response
	params   map[string]string
}

func NewRequestResponsePairWithMap(req packaging.HttpRequest, resp packaging.HttpResponse, params map[string]string) *RequestResponsePairWithMap {
	return &RequestResponsePairWithMap{
		request:  Request(req),
		response: Response(resp),
		params:   params,
	}
}

func (p *RequestResponsePairWithMap) Type() Type {
	return NewType(TypeRequest|TypeResponse|TypeMap, InternalTypeNone)
}

func (p *RequestResponsePairWithMap) Request() packaging.HttpRequest {
	return p.request.Request()
}

func (p *RequestResponsePairWithMap) Response() packaging.HttpResponse {
	return p.response.Response()
}

func (p *RequestResponsePairWithMap) Map() map[string]string {
	return p.params
}

func (p *RequestResponsePairWithMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"request":  p.request,
		"response": p.response,
		"params":   p.params,
	})
}

func (p *RequestResponsePairWithMap) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v["request"] != nil {
		if err := json.Unmarshal(v["request"].([]byte), &p.request); err != nil {
			return err
		}
	}
	if v["response"] != nil {
		if err := json.Unmarshal(v["response"].([]byte), &p.response); err != nil {
			return err
		}
	}
	if v["params"] != nil {
		if err := json.Unmarshal(v["params"].([]byte), &p.params); err != nil {
			return err
		}
	}
	return nil
}
