package scan

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"gorm.io/gorm"

	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/handlers/websocket"
)

// FindSubdomains runs subfinder and passes the results to a callback function
func FindSubdomains(domain models.Domain, db *gorm.DB, ws *websocket.Pool) {
	// status
	slog.Info("[subfinder]", "message", "scanning domain", "domain", domain.Name)
	domain.Status = "scanning"
	db.Save(&domain)
	domain.BroadcastSync(ws)

	results := []string{}

	// run subfinder
	subfinderOpts := &runner.Options{
		Threads:            10, // Thread controls the number of threads to use for active enumerations
		Timeout:            30, // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 10, // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
		ResultCallback: func(result *resolve.HostEntry) {
			handleSubdomainResult(domain, result, db, ws, &results)
		},
		// ProviderConfig: "your_provider_config.yaml",
	}

	subfinder, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		slog.Error("[subfinder]", "message", "failed to create subfinder runner", "domain", domain.Name, "error", err)
		return
	}

	output := &bytes.Buffer{}
	// To run subdomain enumeration on a single domain
	if err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), domain.Name, []io.Writer{output}); err != nil {
		slog.Error("[subfinder]", "message", "failed to enumerate domain", "domain", domain.Name, "error", err)
		return
	}

	err = updateDomainHostCount(&domain, db)
	if err != nil {
		slog.Error("[subfinder]", "message", "failed to update domain host count", "domain", domain.Name, "error", err)
		return
	}

	slog.Info("[subfinder]", "message", "probing hosts", "domain", domain.Name, "count", len(results))

	ProbeHosts(&domain, results, db, ws)
}

// handleSubdomainResult writes the subdomain result to the database and sends a broadcast websocket message
func handleSubdomainResult(domain models.Domain, result *resolve.HostEntry, db *gorm.DB, ws *websocket.Pool, results *[]string) {
	slog.Info("[subfinder]", "message", "subdomain result", "result", result)

	*results = append(*results, result.Host)

	if ws == nil {
		slog.Error("[subfinder]", "message", "no websocket pool available", "domain", domain.Name)
		return
	}

	if db == nil {
		slog.Error("[subfinder]", "message", "no database connection available", "domain", domain.Name)
		return
	}

	host := models.Host{
		ProjectID: domain.ProjectID,
		DomainID:  domain.ID,
		Name:      result.Host,
		Source:    result.Source,
	}

	res := db.Where(models.Host{
		ProjectID: host.ProjectID,
		DomainID:  host.DomainID,
		Name:      host.Name,
	}).FirstOrCreate(&host)
	if res.Error != nil {
		slog.Error("[subfinder]", "message", "failed to create host", "error", res.Error)
		return
	}
}
