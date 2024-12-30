package types

import (
	"net/http"
	"time"
)

type UserRole string

const (
	UserRoleAgent  UserRole = "agent"
	UserRoleAdmin  UserRole = "admin"
	UserRoleViewer UserRole = "viewer"
)

type DomainStatus string

const (
	DomainStatusPending  DomainStatus = "pending"
	DomainStatusScanning DomainStatus = "scanning"
	DomainStatusProbing  DomainStatus = "probing"
	DomainStatusIdle     DomainStatus = "idle"
)

type RequestSource string

const (
        RequestSourceUser      RequestSource = "user"
        RequestSourceProxy     RequestSource = "proxy"
        RequestSourceSubfinder RequestSource = "subfinder"
        RequestSourceHttpx     RequestSource = "httpx"
        RequestSourceFuzzer    RequestSource = "fuzzer"
        RequestSourceReplay    RequestSource = "replay"
        RequestSourceBrute     RequestSource = "brute"
)

type ReportStatus string

const (
	ReportStatusCreated ReportStatus = "created"
	ReportStatusPending ReportStatus = "pending"
	ReportStatusDone    ReportStatus = "done"
)

type ReplayInput struct {
	Method  string `json:"method"`
	URL     string `json:"url"`
	Headers string `json:"headers"`
	Body    string `json:"body"`
}

type ReplayResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Headers    http.Header `json:"headers"`
	Body       []byte      `json:"body"`
}

// AgentSession tracks agent interactions with a chat session
type AgentSession struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    string    `json:"status"`
}
