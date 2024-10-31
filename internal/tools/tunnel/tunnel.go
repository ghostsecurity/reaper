package tunnel

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

type Tunnel struct {
	URL      *string
	Failed   bool
	shutdown chan struct{}
}

func NewTunnel() *Tunnel {
	return &Tunnel{
		shutdown: make(chan struct{}),
	}
}

func (t *Tunnel) Start() error {
	slog.Info("Starting tunnel")
	ctx := context.Background()

	// local listener
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	backend, _ := url.Parse(fmt.Sprintf("http://%s:%s", host, port))

	slog.Info("Tunnel local listener", "backend", backend)

	go func() {
		fwd, err := ngrok.ListenAndForward(
			ctx,
			backend,
			config.HTTPEndpoint(),
			// DEBUG
			ngrok.WithLogger(&customLogger{logger: slog.New(slog.NewTextHandler(os.Stdout, nil))}),
			ngrok.WithAuthtokenFromEnv(),
		)
		if err != nil {
			slog.Error("Error starting tunnel", "error", err)
			return
		}

		session := fwd.Session()
		for _, v := range session.Warnings() {
			slog.Warn("Tunnel warning", "warning", v)
		}

		url := fwd.URL()

		if url == "" {
			slog.Error("Tunnel could not be started")
			t.Failed = true // tunnel failed to start
			session.Close()
			return
		}

		slog.Info("Tunnel started", "url", url)
		t.URL = &url

		<-t.shutdown
		if err := session.Close(); err != nil {
			slog.Error("Error stopping tunnel", "error", err)
		}
		slog.Info("Tunnel stopped")
	}()

	return nil
}

func (t *Tunnel) Stop() {
	slog.Info("Stopping tunnel")
	close(t.shutdown)
}

type customLogger struct {
	logger *slog.Logger
}

func (l *customLogger) Log(ctx context.Context, level int, message string, data map[string]interface{}) {
	// Implement logging logic here
	l.logger.Log(ctx, slog.Level(level), message, "data", data)
}
