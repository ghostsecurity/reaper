package models

import (
	"github.com/ghostsecurity/reaper/internal/handlers/websocket"
	"github.com/ghostsecurity/reaper/internal/types"
)

// Sync methods are used to sync data between the database and websocket clients

func (d *Domain) BroadcastSync(ws *websocket.Pool) {
	msg := &types.DomainSyncMessage{
		Type:          types.MessageTypeScanSyncDomain,
		ID:            d.ID,
		Status:        d.Status,
		HostCount:     d.HostCount,
		LastScannedAt: d.LastScannedAt,
	}
	ws.Broadcast <- msg
}
