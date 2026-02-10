package proxy

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ghostsecurity/reaper/internal/storage"
)

// Event describes a proxied request for display logging.
type Event struct {
	ID          int64
	Method      string
	Scheme      string
	Host        string
	Path        string
	StatusCode  int
	DurationMs  int64
	Intercepted bool
}

// Proxy is an HTTP/HTTPS MITM proxy.
type Proxy struct {
	Scope     *Scope
	Store     storage.Store
	CA        *CA
	Transport http.RoundTripper // optional; defaults to http.DefaultTransport
	OnEvent   func(Event)       // optional callback for live activity display

	certCache sync.Map // host → *tls.Certificate
}

func (p *Proxy) transport() http.RoundTripper {
	if p.Transport != nil {
		return p.Transport
	}
	return http.DefaultTransport
}

func (p *Proxy) emit(e Event) {
	if p.OnEvent != nil {
		p.OnEvent(e)
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.handleConnect(w, r)
		return
	}
	p.handleHTTP(w, r)
}

func (p *Proxy) handleConnect(w http.ResponseWriter, r *http.Request) {
	hij, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "hijacking not supported", http.StatusInternalServerError)
		return
	}

	host := r.Host
	if !strings.Contains(host, ":") {
		host += ":443"
	}

	clientConn, buf, err := hij.Hijack()
	if err != nil {
		return
	}

	// Send 200 after hijacking so we control the flush
	buf.WriteString("HTTP/1.1 200 Connection Established\r\n\r\n")
	buf.Flush()

	hostname := r.Host
	if idx := strings.Index(hostname, ":"); idx != -1 {
		hostname = hostname[:idx]
	}

	if !p.Scope.InScope(hostname) {
		p.emit(Event{
			Scheme:      "https",
			Host:        hostname,
			Intercepted: false,
		})
		p.blindRelay(clientConn, host)
		return
	}

	p.mitmConnect(clientConn, hostname, host)
}

func (p *Proxy) blindRelay(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	upstream, err := net.DialTimeout("tcp", targetAddr, 10*time.Second)
	if err != nil {
		return
	}
	defer upstream.Close()

	done := make(chan struct{})
	go func() {
		io.Copy(upstream, clientConn)
		close(done)
	}()
	io.Copy(clientConn, upstream)
	<-done
}

func (p *Proxy) mitmConnect(clientConn net.Conn, hostname, targetAddr string) {
	defer clientConn.Close()

	tlsCert, err := p.getCertForHost(hostname)
	if err != nil {
		return
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*tlsCert},
	}

	tlsConn := tls.Server(clientConn, tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		return
	}
	defer tlsConn.Close()

	reader := bufio.NewReader(tlsConn)

	for {
		req, err := http.ReadRequest(reader)
		if err != nil {
			return
		}

		req.URL.Scheme = "https"
		req.URL.Host = targetAddr
		req.RequestURI = ""

		p.proxyAndLog(tlsConn, req, "https", hostname)
	}
}

func (p *Proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	if !r.URL.IsAbs() {
		http.Error(w, "non-proxy request", http.StatusBadRequest)
		return
	}

	hostname := r.URL.Hostname()
	inScope := p.Scope.InScope(hostname)

	// Forward the request
	start := time.Now()
	r.RequestURI = ""
	r.Header.Del("Accept-Encoding")

	resp, err := p.transport().RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	duration := time.Since(start).Milliseconds()

	if inScope {
		reqDump, _ := httputil.DumpRequest(r, true)
		entry := &storage.Entry{
			Method:          r.Method,
			Scheme:          "http",
			Host:            hostname,
			Path:            r.URL.Path,
			Query:           r.URL.RawQuery,
			RequestHeaders:  r.Header.Clone(),
			RequestBody:     extractBody(reqDump),
			StatusCode:      resp.StatusCode,
			ResponseHeaders: resp.Header.Clone(),
			ResponseBody:    body,
			Timestamp:       time.Now(),
			DurationMs:      duration,
		}
		p.Store.Save(entry)
		p.emit(Event{
			ID:          entry.ID,
			Method:      r.Method,
			Scheme:      "http",
			Host:        hostname,
			Path:        r.URL.Path,
			StatusCode:  resp.StatusCode,
			DurationMs:  duration,
			Intercepted: true,
		})
	} else {
		p.emit(Event{
			Method:      r.Method,
			Scheme:      "http",
			Host:        hostname,
			Path:        r.URL.Path,
			StatusCode:  resp.StatusCode,
			DurationMs:  duration,
			Intercepted: false,
		})
	}

	// Copy response headers
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func (p *Proxy) proxyAndLog(clientConn net.Conn, req *http.Request, scheme, hostname string) {
	start := time.Now()

	reqBody, _ := io.ReadAll(req.Body)
	req.Body.Close()

	// Forward to upstream
	upstreamReq, err := http.NewRequest(req.Method, req.URL.String(), strings.NewReader(string(reqBody)))
	if err != nil {
		return
	}
	upstreamReq.Header = req.Header.Clone()
	upstreamReq.Header.Del("Accept-Encoding")

	resp, err := p.transport().RoundTrip(upstreamReq)
	if err != nil {
		// Write error response to client
		fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\nContent-Length: 0\r\n\r\n")
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	duration := time.Since(start).Milliseconds()

	entry := &storage.Entry{
		Method:          req.Method,
		Scheme:          scheme,
		Host:            hostname,
		Path:            req.URL.Path,
		Query:           req.URL.RawQuery,
		RequestHeaders:  req.Header.Clone(),
		RequestBody:     reqBody,
		StatusCode:      resp.StatusCode,
		ResponseHeaders: resp.Header.Clone(),
		ResponseBody:    respBody,
		Timestamp:       time.Now(),
		DurationMs:      duration,
	}
	p.Store.Save(entry)
	p.emit(Event{
		ID:          entry.ID,
		Method:      entry.Method,
		Scheme:      scheme,
		Host:        hostname,
		Path:        req.URL.Path,
		StatusCode:  resp.StatusCode,
		DurationMs:  duration,
		Intercepted: true,
	})

	// Write response back to client manually
	resp.Header.Del("Transfer-Encoding")
	resp.Header.Del("Content-Encoding")
	resp.Header.Set("Content-Length", strconv.Itoa(len(respBody)))

	var sb strings.Builder
	fmt.Fprintf(&sb, "HTTP/%d.%d %s\r\n", resp.ProtoMajor, resp.ProtoMinor, resp.Status)
	resp.Header.Write(&sb)
	sb.WriteString("\r\n")

	clientConn.Write([]byte(sb.String()))
	clientConn.Write(respBody)
}

func (p *Proxy) getCertForHost(host string) (*tls.Certificate, error) {
	if cached, ok := p.certCache.Load(host); ok {
		return cached.(*tls.Certificate), nil
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: host,
		},
		DNSNames:    []string{host},
		NotBefore:   time.Now().Add(-time.Hour),
		NotAfter:    time.Now().Add(24 * time.Hour),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, p.CA.Cert, &certKey.PublicKey, p.CA.Key)
	if err != nil {
		return nil, err
	}

	tlsCert := &tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  certKey,
	}

	p.certCache.Store(host, tlsCert)
	return tlsCert, nil
}

func extractBody(dump []byte) []byte {
	idx := strings.Index(string(dump), "\r\n\r\n")
	if idx == -1 {
		return nil
	}
	body := dump[idx+4:]
	if len(body) == 0 {
		return nil
	}
	return body
}

// StatusText returns standard status text — used by response formatting.
func StatusText(code int) string {
	return http.StatusText(code)
}
