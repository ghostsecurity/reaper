package scan

import (
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/httpx/runner"
	"gorm.io/gorm"

	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/handlers/websocket"
	"github.com/ghostsecurity/reaper/internal/types"
)

func ProbeHosts(domain *models.Domain, hosts []string, db *gorm.DB, ws *websocket.Pool) {
	// status
	if domain != nil {
		slog.Info("[httpx]", "message", "probing hosts", "domain", domain.Name, "hosts", len(hosts))
		domain.Status = types.DomainStatusProbing
		db.Save(domain)
		domain.BroadcastSync(ws)
	}

	// run httpx
	options := runner.Options{
		DisableStdin:    true, // Running as a server, no stdin
		Methods:         "GET",
		InputTargetHost: goflags.StringSlice(hosts),
		OnResult: func(r runner.Result) {
			handleHttpxResult(r.Input, r, db, ws)
		},
		FollowRedirects: true,
		MaxRedirects:    5,
		Timeout:         10,
	}

	if err := options.ValidateOptions(); err != nil {
		slog.Error("[httpx]", "message", "invalid httpx options", "error", err)
		return
	}

	httpxRunner, err := runner.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	defer httpxRunner.Close()

	httpxRunner.RunEnumeration()

	if domain != nil {
		t := time.Now()
		domain.Status = "completed"
		domain.LastScannedAt = &t
		db.Save(domain)
		domain.BroadcastSync(ws)
	}
}

func handleHttpxResult(host string, r runner.Result, db *gorm.DB, ws *websocket.Pool) {
	// DEBUG full detail
	// slog.Info("[httpx]", "message", "httpx result", "host", host, "result", r)
	// DEBUG just the host
	slog.Info("[httpx]", "message", "httpx result", "host", host)

	if r.Err != nil {
		slog.Error("[httpx]", "message", "httpx error", "host", host, "error", r.Err)
		return
	}

	if ws == nil {
		slog.Error("[httpx]", "message", "no websocket pool available", "host", host)
		return
	}

	if db == nil {
		slog.Error("[httpx]", "message", "no database connection available", "host", host)
		return
	}

	h := models.Host{
		Name:        host,
		Status:      "live",
		StatusCode:  r.StatusCode,
		Scheme:      r.Scheme,
		ContentType: r.ContentType,
		CDNName:     r.CDNName,
		CDNType:     r.CDNType,
		Webserver:   r.WebServer,
		Tech:        strings.Join(r.Technologies, "|"),
	}

	res := db.Model(&h).Where(&models.Host{Name: host}).Updates(h)
	if res.Error != nil {
		slog.Error("[httpx]", "message", "failed to update host", "host", host, "error", res.Error)
	}
}
