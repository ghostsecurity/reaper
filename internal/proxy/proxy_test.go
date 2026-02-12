package proxy

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ghostsecurity/reaper/internal/storage"
)

type nullStore struct{}

func (s *nullStore) Save(e *storage.Entry) error                                  { return nil }
func (s *nullStore) Get(id int64) (*storage.Entry, error)                         { return nil, fmt.Errorf("not found") }
func (s *nullStore) List(l, o int) ([]*storage.Entry, error)                      { return nil, nil }
func (s *nullStore) Search(p storage.SearchParams) ([]*storage.Entry, error)      { return nil, nil }
func (s *nullStore) ListAfter(afterID int64, limit int) ([]*storage.Entry, error) { return nil, nil }
func (s *nullStore) Clear() error                                                 { return nil }
func (s *nullStore) Close() error                                                 { return nil }

func startTestProxy(t *testing.T, domains []string, transport http.RoundTripper) (*Proxy, net.Listener) {
	t.Helper()

	ca, err := GenerateCA()
	if err != nil {
		t.Fatal(err)
	}

	p := &Proxy{
		Scope:     NewScope(domains, nil),
		Store:     &nullStore{},
		CA:        ca,
		Transport: transport,
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	srv := &http.Server{Handler: p}
	go srv.Serve(ln)
	t.Cleanup(func() { srv.Close(); ln.Close() })

	return p, ln
}

func TestProxyMITMGzip(t *testing.T) {
	// Upstream server that always responds with gzip
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json")
		gz := gzip.NewWriter(w)
		gz.Write([]byte(`{"status":"ok"}`))
		gz.Close()
	}))
	defer upstream.Close()

	upstreamURL, _ := url.Parse(upstream.URL)
	host, _, _ := net.SplitHostPort(upstreamURL.Host)

	// Proxy transport that trusts the test server
	proxyTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	_, proxyLn := startTestProxy(t, []string{host}, proxyTransport)
	proxyAddr := proxyLn.Addr().String()

	// Client that uses the proxy and trusts the MITM cert
	proxyURL, _ := url.Parse("http://" + proxyAddr)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 5 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(upstream.URL + "/test")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("request failed (took %v): %v", elapsed, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}

	if string(body) != `{"status":"ok"}` {
		t.Errorf("body = %q, want %q", body, `{"status":"ok"}`)
	}

	if elapsed > 3*time.Second {
		t.Errorf("request took %v, expected < 3s", elapsed)
	}

	t.Logf("completed in %v", elapsed)
}

func TestProxyHTTP(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"hello":"world"}`))
	}))
	defer upstream.Close()

	_, proxyLn := startTestProxy(t, []string{"127.0.0.1"}, nil)
	proxyAddr := proxyLn.Addr().String()

	proxyURL, _ := url.Parse("http://" + proxyAddr)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 5 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(upstream.URL + "/test")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}

	if string(body) != `{"hello":"world"}` {
		t.Errorf("body = %q, want %q", body, `{"hello":"world"}`)
	}

	if elapsed > 3*time.Second {
		t.Errorf("request took %v, expected < 3s", elapsed)
	}

	t.Logf("completed in %v", elapsed)
}
