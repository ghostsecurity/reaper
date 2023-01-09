package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/ghostsecurity/reaper/backend/log"
	"github.com/ghostsecurity/reaper/backend/settings"
)

type Proxy struct {
	client    *http.Client
	server    *http.Server
	reqFuncs  []ProxyRequestFunc
	respFuncs []ProxyResponseFunc
	logger    *log.Logger
	addr      string
}

type proxyLogger struct {
	logger *log.Logger
}

func (l *proxyLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(log.LevelDebug, format, v...)
}

type ProxyRequestFunc func(*http.Request, int64) (*http.Request, *http.Response)
type ProxyResponseFunc func(*http.Response, int64) *http.Response

func New(userSettings *settings.Provider, logger *log.Logger) (*Proxy, error) {
	// TODO: allow user to specify which ip to bind to
	addr := fmt.Sprintf("127.0.0.1:%d", userSettings.Get().ProxyPort)
	logger.Infof("Creating proxy on %s...", addr)
	proxy := goproxy.NewProxyHttpServer()
	ca, err := tls.X509KeyPair(userSettings.Get().CACert, userSettings.Get().CAKey)
	if err != nil {
		return nil, err
	}
	if ca.Leaf, err = x509.ParseCertificate(ca.Certificate[0]); err != nil {
		return nil, err
	}
	logger.Infof("Using CA from organisation '%s'", strings.Join(ca.Leaf.Subject.Organization, ", "))
	tlsConfig := goproxy.TLSConfigFromCA(&ca)

	// the proxy module we're using requires setting this bunch of globals to use our own CA - kind of gross
	goproxy.OkConnect = &goproxy.ConnectAction{
		Action:    goproxy.ConnectAccept,
		TLSConfig: tlsConfig,
	}
	goproxy.MitmConnect = &goproxy.ConnectAction{
		Action:    goproxy.ConnectMitm,
		TLSConfig: tlsConfig,
	}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{
		Action:    goproxy.ConnectHTTPMitm,
		TLSConfig: tlsConfig,
	}
	goproxy.RejectConnect = &goproxy.ConnectAction{
		Action:    goproxy.ConnectReject,
		TLSConfig: tlsConfig,
	}

	// this at least lets us set things up a little without overriding another goproxy global
	proxy.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&ca)}, host
	}))
	proxy.Verbose = userSettings.Get().LogLevel <= log.LevelDebug
	proxy.Logger = &proxyLogger{logger: logger.WithPrefix("inner")}
	srv := &http.Server{
		Addr:    addr,
		Handler: proxy,
	}
	p := &Proxy{
		client: &http.Client{},
		server: srv,
		logger: logger,
		addr:   addr,
	}

	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		for _, f := range p.reqFuncs {
			modified, resp := f(r, ctx.Session)
			if resp != nil {
				return modified, resp
			}
			r = modified
		}
		return r, nil
	})
	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		for _, f := range p.respFuncs {
			r = f(r, ctx.Session)
		}
		return r
	})

	return p, nil
}

func (p *Proxy) Addr() string {
	return p.addr
}

func (p *Proxy) OnRequest(f ProxyRequestFunc) *Proxy {
	p.reqFuncs = append(p.reqFuncs, f)
	return p
}

func (p *Proxy) OnResponse(f ProxyResponseFunc) *Proxy {
	p.respFuncs = append(p.respFuncs, f)
	return p
}

func (p *Proxy) Run() error {
	p.logger.Infof("Listening on %s...", p.server.Addr)
	if err := p.server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}
	return nil
}

func (p *Proxy) Close() error {
	return p.server.Close()
}
