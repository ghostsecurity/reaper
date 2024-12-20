package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/ghostsecurity/reaper/internal/database"
	"github.com/ghostsecurity/reaper/internal/handlers"
	ws "github.com/ghostsecurity/reaper/internal/handlers/websocket"
	"github.com/ghostsecurity/reaper/internal/middleware"
)

//go:embed dist/*
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

	app := fiber.New(fiber.Config{
		AppName: "Reaper",
	})
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

	// handler
	h := handlers.NewHandler(pool, db)

	// status
	app.Get("/status", h.Status)
	app.Post("/register", h.Register)
	api := app.Group("/api", middleware.TokenAuth(db))

	api.Post("/navigation", h.Navigation)

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
	api.Get("/attacks", h.GetAttacks)
	api.Get("/attacks/:id", h.GetAttack)
	api.Get("/attacks/:id/results", h.GetAttackResults)
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

	// Generate OpenAPI spec dynamically
	app.Get("/openapi.json", func(c *fiber.Ctx) error {
		return c.JSON(generateOpenAPISpec(app))
	})

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

	// Start server
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	listener := fmt.Sprintf("%s:%s", host, port)
	err := app.Listen(listener)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func generateOpenAPISpec(app *fiber.App) OpenAPISpec {
	routes := app.GetRoutes()
	var openAPIPaths []OpenAPIPath
	excludedMethods := []string{fiber.MethodConnect, fiber.MethodOptions, fiber.MethodTrace}

	for _, route := range routes {
		if slices.Contains(excludedMethods, route.Method) {
			continue
		}
		description := extractAnnotation(route.Path)
		operationObject := OpenAPIOperation{
			Summary:     description,
			Description: description,
		}

		operationObject.Summary = route.Path
		operationObject.Description = fmt.Sprintf("description: %s %s", route.Method, route.Path)

		path := OpenAPIPath{
			Description: description,
		}

		// skip all  middleware handlers; e.g. handler starts with "github.com/gofiber/fiber/v2/middleware"
		handlers := parseHandlers(route.Handlers)
		isMiddleware := false
		for _, handler := range handlers {
			if strings.Contains(handler, "/middleware/") {
				isMiddleware = true
				break
			}
		}
		if isMiddleware {
			continue
		}

		switch route.Method {
		case fiber.MethodGet:
			path.Get = &operationObject
		case fiber.MethodPost:
			path.Post = &operationObject
		case fiber.MethodDelete:
			path.Delete = &operationObject
		case fiber.MethodPut:
			path.Put = &operationObject
		case fiber.MethodPatch:
			path.Patch = &operationObject
		}

		openAPIPaths = append(openAPIPaths, path)
	}

	// sort paths by summary
	slices.SortFunc(openAPIPaths, func(a, b OpenAPIPath) int {
		return strings.Compare(a.Summary, b.Summary)
	})

	return OpenAPISpec{
		Version: "3.1.1",
		Info: OpenAPIInfo{
			Title:       app.Config().AppName,
			Description: "Reaper API",
			Version:     "1.0.0",
		},
		Paths: openAPIPaths,
	}
}

func extractAnnotation(routePath string) string {
	return fmt.Sprintf("%s - path description for", routePath)
}

func parseHandlers(handlers []func(*fiber.Ctx) error) []string {
	var handlerNames []string
	for _, handler := range handlers {
		handlerNames = append(handlerNames, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
	}
	return handlerNames
}

type OpenAPIPath struct {
	Get         *OpenAPIOperation `json:"get,omitempty"`
	Post        *OpenAPIOperation `json:"post,omitempty"`
	Delete      *OpenAPIOperation `json:"delete,omitempty"`
	Put         *OpenAPIOperation `json:"put,omitempty"`
	Patch       *OpenAPIOperation `json:"patch,omitempty"`
	Options     *OpenAPIOperation `json:"options,omitempty"`
	Head        *OpenAPIOperation `json:"head,omitempty"`
	Connect     *OpenAPIOperation `json:"connect,omitempty"`
	Trace       *OpenAPIOperation `json:"trace,omitempty"`
	Summary     string            `json:"summary"`     // An optional string summary, intended to apply to all operations in this path.
	Description string            `json:"description"` // An optional string description, intended to apply to all operations in this path. [CommonMark] syntax MAY be used for rich text representation.
}

type OpenAPIOperation struct {
	ID          string         `json:"operationId"`
	Tags        []string       `json:"tags"`        // A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.
	Summary     string         `json:"summary"`     // A short summary of what the operation does.
	Description string         `json:"description"` // A verbose explanation of the operation behavior. [CommonMark] syntax MAY be used for rich text representation.
	Params      []OpenAPIParam `json:"parameters"`
}

type OpenAPIParam struct {
	Name        string `json:"name"` // REQUIRED. The name of the parameter. Parameter names are case sensitive.
	In          string `json:"in"`   // REQUIRED. The location of the parameter. Possible values are "query", "header", "path" or "cookie".
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type OpenAPISpec struct {
	Version string        `json:"openapi"`
	Info    OpenAPIInfo   `json:"info"`
	Title   string        `json:"title"`
	Paths   []OpenAPIPath `json:"paths"`
}

type OpenAPIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}
