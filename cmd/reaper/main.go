package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/ghostsecurity/reaper/internal/browser"
	"github.com/ghostsecurity/reaper/internal/database"
	"github.com/ghostsecurity/reaper/internal/handlers"
	ws "github.com/ghostsecurity/reaper/internal/handlers/websocket"
	"github.com/ghostsecurity/reaper/internal/middleware"
	"github.com/ghostsecurity/reaper/internal/tools/proxy"
)

//go:embed dist/* frontend/browser/*
var static embed.FS

func main() {
	// We don't need the time field in the local logs
	opts := slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))
	slog.SetDefault(logger)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		// TODO: dynamically set this based on the ngrok URL
		AllowOriginsFunc: func(origin string) bool {
			return origin == "http://localhost:5173" || strings.HasSuffix(origin, ".ngrok.app")
		},
		AllowMethods:     "GET,DELETE,POST,OPTIONS",
		AllowHeaders:     fmt.Sprintf("%s, %s", fiber.HeaderContentType, middleware.AuthTokenHeader),
		AllowCredentials: true,
	}))

	// Websocket client pool
	pool := ws.NewPool()
	go pool.Start()

	// Websocket handler
	wsh := handlers.NewWsHandler(pool)
	app.Use("/ws", wsh.WebSocketUpgrade)
	app.Get("/ws", websocket.New(wsh.WebSocket, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}))

	// database
	database.Migrate()
	db := database.Connect()

	// Initialize proxy
	proxy := proxy.NewProxy(pool, db)

	// Initialize browser
	browserInstance := browser.NewBrowser(proxy)

	// handler
	h := handlers.NewHandler(pool, db)
	bh := handlers.NewBrowserHandler(browserInstance)

	// status
	app.Get("/status", h.Status)
	app.Post("/register", h.Register)
	api := app.Group("/api", middleware.TokenAuth(db))

	api.Post("/navigation", h.Navigation)

	// Browser routes
	browserGroup := api.Group("/browser")
	browserGroup.Post("/start", bh.HandleStart)
	browserGroup.Post("/stop", bh.HandleStop)
	browserGroup.Post("/navigate", bh.HandleNavigate)
	browserGroup.Post("/reload", bh.HandleReload)

	// Browser WebSocket for VNC
	app.Get("/browser/vnc", websocket.New(bh.HandleVNC))

	// scan
	api.Post("/scan/domains", h.CreateDomain)
	api.Get("/scan/domains", h.GetDomains)
	api.Get("/scan/domains/:id", h.GetDomain)
	api.Delete("/scan/domains/:id", h.DeleteDomain)
	api.Get("/scan/domains/:id/hosts", h.GetDomainHosts)
	// explore
	api.Get("/proxy/status", h.ProxyStatus)
	api.Post("/proxy/start", h.ProxyStart)
	api.Post("/proxy/stop", h.ProxyStop)
	api.Get("/explore/host", h.ExploreHostExample)
	api.Get("/explore/endpoint", h.ExploreEndpointExample)
	// crawl
	// replay
	api.Get("/requests", h.GetRequests)
	api.Get("/requests/:id", h.GetRequest)
	api.Post("/replay", h.Replay)
	// attack
	api.Get("/endpoints", h.GetEndpoints)
	api.Get("/endpoints/:id", h.GetEndpoint)
	api.Post("/attack", h.CreateAttack)
	// api.Get("/attacks", h.GetAttacks)
	// api.Get("/attacks/:id", h.GetAttack)
	api.Delete("/attack/:id/results", h.DeleteAttackResults)
	// fuzz
	// automate
	// collaborate
	api.Get("/tunnel/status", h.TunnelStatus)
	api.Post("/tunnel/start", h.TunnelStart)
	api.Post("/tunnel/stop", h.TunnelStop)
	// ai assist
	api.Get("/agent/sessions", h.GetSessions)
	api.Get("/agent/sessions/:id", h.GetSession)
	api.Post("/agent/sessions", h.CreateSession)
	api.Delete("/agent/sessions/:id", h.DeleteSession)
	api.Get("/agent/sessions/:id/messages", h.GetSessionMessages)
	api.Post("/agent/sessions/:id/messages", h.CreateSessionMessage)
	// reports
	api.Get("/reports", h.GetReports)
	api.Get("/reports/:id", h.GetReport)
	api.Post("/reports", h.CreateReport)
	api.Delete("/reports/:id", h.DeleteReport)
	// settings

	// serve static frontend files
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(static),
		PathPrefix:   "dist",
		Browse:       true,
		NotFoundFile: "dist/index.html",
	}))
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(static),
		PathPrefix: "dist/assets",
		Browse:     true,
	}))
	app.Use("/browser", filesystem.New(filesystem.Config{
		Root:       http.FS(static),
		PathPrefix: "frontend/browser",
		Browse:     true,
	}))

	// Start server
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	listener := fmt.Sprintf("%s:%s", host, port)
	err := app.Listen(listener)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
