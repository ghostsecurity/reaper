package handlers

import (
	"log"

	"gorm.io/gorm"

	ws "github.com/ghostsecurity/reaper/internal/handlers/websocket"
	"github.com/ghostsecurity/reaper/internal/tools/proxy"
	"github.com/ghostsecurity/reaper/internal/tools/tunnel"
)

type Handler struct {
	proxy  *proxy.Proxy
	tunnel *tunnel.Tunnel
	pool   *ws.Pool
	db     *gorm.DB
}

func NewHandler(pool *ws.Pool, db *gorm.DB) *Handler {
	proxy := proxy.NewProxy(pool, db)
	// start proxy by default
	err := proxy.Start()
	if err != nil {
		log.Fatalf("error starting proxy: %v", err)
	}
	return &Handler{
		proxy: proxy,
		pool:  pool,
		db:    db,
	}
}

func NewWsHandler(pool *ws.Pool) *Handler {
	return &Handler{
		pool: pool,
	}
}
