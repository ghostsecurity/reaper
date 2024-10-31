package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"

	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/handlers/websocket"
	"github.com/ghostsecurity/reaper/internal/service"
	"github.com/ghostsecurity/reaper/internal/types"
)

type Proxy struct {
	proxy    *goproxy.ProxyHttpServer
	pool     *websocket.Pool
	db       *gorm.DB
	shutdown chan struct{}
}

func NewProxy(pool *websocket.Pool, db *gorm.DB) *Proxy {
	p := goproxy.NewProxyHttpServer()
	p.Verbose = true
	return &Proxy{
		proxy:    p,
		pool:     pool,
		db:       db,
		shutdown: make(chan struct{}),
	}
}

// Start starts the proxy in the background
func (p *Proxy) Start() error {
	host := os.Getenv("HOST")
	port := os.Getenv("PROXY_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	slog.Info("Starting proxy", "address", addr)

	// TODO: clean this up
	go func() {
		server := &http.Server{Addr: addr, Handler: p.proxy}
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("Proxy server error", "error", err)
			}
		}()

		initializeCA()

		// https CONNECT
		p.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

		// request
		p.proxy.OnRequest().DoFunc(p.httpRequestHandler)

		// response
		p.proxy.OnResponse().DoFunc(p.httpResponseHandler)

		<-p.shutdown
		if err := server.Close(); err != nil {
			slog.Error("Error stopping proxy server", "error", err)
		}
		slog.Info("Proxy stopped")
	}()

	return nil
}

// Stop stops the proxy
func (p *Proxy) Stop() {
	slog.Info("Stopping proxy")
	close(p.shutdown)
}

func (p *Proxy) httpRequestHandler(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	uri := req.URL
	slog.Info("Proxy HTTP request", "remote", uri)

	// save req to db
	r, err := p.saveReq(req, ctx)
	if err != nil {
		slog.Error("could not save request to db", "error", err.Error())
	}

	m := &types.ProxyMessage{
		Type: types.MessageTypeExploreRequest,
		Host: uri.Hostname(),
		Path: uri.Path,
	}

	p.pool.Broadcast <- m

	return r, nil
}

func (p *Proxy) httpResponseHandler(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	slog.Info("Proxy HTTP response", "remote", resp.Request.RemoteAddr)

	m := &types.ProxyMessage{
		Type:   types.MessageTypeExploreResponse,
		Host:   resp.Request.URL.Hostname(),
		Method: resp.Request.Method,
		Path:   resp.Request.URL.Path,
		Status: resp.StatusCode,
	}

	p.pool.Broadcast <- m

	// save resp to db
	r, err := p.saveResp(resp, ctx)
	if err != nil {
		slog.Error("could not save response to db", "error", err.Error())
	}

	return r
}

// saveReq saves the request to the database, saves the request ID to the context, and returns the original request
func (p *Proxy) saveReq(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, error) {
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		slog.Error("[proxy req handler]", "msg", "error reading request body", "error", err)
		return nil, err
	}

	request := models.Request{
		Source:        types.RequestSourceProxy,
		Method:        req.Method,
		Host:          req.Host,
		URL:           req.URL.String(),
		Headers:       headersToString(req.Header),
		Proto:         req.Proto,
		ProtoMajor:    req.ProtoMajor,
		ProtoMinor:    req.ProtoMinor,
		ContentType:   req.Header.Get("Content-Type"),
		ContentLength: req.ContentLength,
		HeaderKeys:    keysToString(req.Header),
		ParamKeys:     paramKeysToString(req.URL.Query()),
		BodyKeys:      bodyKeysToString(requestBody),
		Body:          string(requestBody),
	}

	result := p.db.Create(&request)
	if result.Error != nil {
		slog.Error("[proxy req handler]", "msg", "error writing request to db", "error", result.Error)
	}

	// save request ID to ctx
	ctx.UserData = request.ID

	// replace the original request body
	req.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	return req, nil
}

// saveResp saves the response to the database, along with the request ID from the context, and returns the original response
func (p *Proxy) saveResp(resp *http.Response, ctx *goproxy.ProxyCtx) (*http.Response, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("[proxy resp handler]", "msg", "error reading response body", "error", err)
		return nil, err
	}

	// only save if the response is json
	body := `not saved`
	json := strings.Contains(resp.Header.Get(fiber.HeaderContentType), fiber.MIMEApplicationJSON)
	if json {
		body = string(responseBody)
	}

	response := models.Response{
		RequestID:     ctx.UserData.(uint), // link to request
		Status:        resp.Status,
		StatusCode:    resp.StatusCode,
		ContentType:   resp.Header.Get(fiber.HeaderContentType),
		ContentLength: resp.ContentLength,
		Headers:       headersToString(resp.Header),
		Body:          body,
	}

	result := p.db.Create(&response)
	if result.Error != nil {
		slog.Error("[proxy resp handler]", "msg", "error writing response to db", "error", result.Error)
	}

	service.CreateEndpoint(p.db, service.EndpointInput{
		Hostname: resp.Request.URL.Hostname(),
		Path:     resp.Request.URL.Path,
		Method:   resp.Request.Method,
	})

	// Replace the original response body
	resp.Body = io.NopCloser(bytes.NewBuffer(responseBody))

	return resp, nil
}

// headersToString converts a map of headers to a newline delimited string
func headersToString(headers http.Header) string {
	var sb strings.Builder
	for key, values := range headers {
		for _, value := range values {
			sb.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
	}
	return sb.String()
}

// keysToString converts a map of headers to a comma-separated string of the header keys
func keysToString(headers http.Header) string {
	var slice []string
	for key := range headers {
		slice = append(slice, key)
	}
	return strings.Join(slice, ",")
}

// paramKeysToString converts a map of URL values (query params) to a comma-separated string of the keys
func paramKeysToString(values url.Values) string {
	var slice []string
	for key := range values {
		slice = append(slice, key)
	}
	return strings.Join(slice, ",")
}

// bodyKeysToString converts a JSON body to a comma-separated string of the keys
func bodyKeysToString(body []byte) string {
	var m map[string]string
	json.Unmarshal(body, &m)
	return strings.Join(maps.Keys(m), ",")
}
