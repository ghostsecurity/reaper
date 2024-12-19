package models

import (
	"time"

	"github.com/ghostsecurity/reaper/internal/types"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Role      types.UserRole `json:"role" gorm:"default:viewer"`
	Enabled   bool           `json:"enabled" gorm:"default:true"`
	Token     string         `json:"token"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Setting struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Project struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Domain struct {
	ID            uint               `json:"id" gorm:"primaryKey"`
	ProjectID     uint               `json:"project_id"`
	Project       Project            `json:"project"`
	Name          string             `json:"name" gorm:"uniqueIndex:idx_project_domain"`
	Status        types.DomainStatus `json:"status"`
	HostCount     int                `json:"host_count"`
	LastScannedAt *time.Time         `json:"last_scanned_at"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

type Host struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	DomainID    uint      `json:"domain_id"`
	ProjectID   uint      `json:"project_id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	StatusCode  int       `json:"status_code"`  // httpx
	Source      string    `json:"source"`       // subfinder
	Scheme      string    `json:"scheme"`       // httpx
	ContentType string    `json:"content_type"` // httpx
	CDNName     string    `json:"cdn_name"`     // httpx
	CDNType     string    `json:"cdn_type"`     // httpx
	Webserver   string    `json:"webserver"`    // httpx
	Tech        string    `json:"tech"`         // httpx
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Endpoint struct {
	ID uint `json:"id" gorm:"primaryKey"`
	// TOOD: use HostID and reference Host table
	// HostID    uint      `json:"host_id"`
	Hostname  string    `json:"hostname"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Params    string    `json:"params"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Request struct {
	ID            uint                `json:"id" gorm:"primaryKey"`
	Source        types.RequestSource `json:"source"`
	Method        string              `json:"method"`         // http.Request.Method
	Host          string              `json:"host"`           // http.Request.Host
	URL           string              `json:"url"`            // http.Request.URL
	Headers       string              `json:"headers"`        // http.Request.Header
	Proto         string              `json:"proto"`          // http.Request.Proto
	ProtoMajor    int                 `json:"proto_major"`    // http.Request.ProtoMajor
	ProtoMinor    int                 `json:"proto_minor"`    // http.Request.ProtoMinor
	ContentType   string              `json:"content_type"`   // http.Request.Header.Get("Content-Type")
	ContentLength int64               `json:"content_length"` // http.Request.ContentLength
	HeaderKeys    string              `json:"header_keys"`
	ParamKeys     string              `json:"param_keys"`
	BodyKeys      string              `json:"body_keys"`
	Body          string              `json:"body"` // http.Request.Body
	Response      Response            `json:"response"`
	CreatedAt     time.Time           `json:"created_at"`
}

type Response struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	RequestID     uint      `json:"request_id" gorm:"foreignKey:ID"`
	Status        string    `json:"status"`         // http.Response.Status
	StatusCode    int       `json:"status_code"`    // http.Response.StatusCode
	ContentType   string    `json:"content_type"`   // http.Response.Header.Get("Content-Type")
	ContentLength int64     `json:"content_length"` // http.Response.ContentLength
	Headers       string    `json:"headers"`        // http.Response.Header
	Body          string    `json:"body"`           // http.Response.Body
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TODO: clean up this mess
type FuzzAttack struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Type      string    `json:"type"` // header, body, param
	Headers   string    `json:"headers"`
	Keys      string    `json:"keys"`
	Params    string    `json:"params"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FuzzResult struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FuzzAttackID uint      `json:"fuzz_attack_id" gorm:"foreignKey:ID"`
	Hostname     string    `json:"hostname"`
	IpAddress    string    `json:"ip_address"`
	Port         string    `json:"port"`
	Scheme       string    `json:"scheme"`
	URL          string    `json:"url"`
	Endpoint     string    `json:"endpoint"`
	Request      string    `json:"request"`
	Response     string    `json:"response"`
	StatusCode   int       `json:"status_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Report struct {
	ID        uint               `json:"id" gorm:"primaryKey"`
	Domain    string             `json:"domain"`
	Title     string             `json:"title"`
	Markdown  string             `json:"markdown"`
	Status    types.ReportStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type AgentSession struct {
	ID          uint                  `json:"id" gorm:"primaryKey"`
	Description string                `json:"description"`
	Messages    []AgentSessionMessage `json:"messages"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type AgentSessionMessage struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	AuthorID       uint           `json:"author_id" gorm:"foreignKey:ID,tablename:users"`
	AuthorRole     types.UserRole `json:"author_role"`
	AgentSessionID uint           `json:"agent_session_id" gorm:"foreignKey:ID"`
	Content        string         `json:"content"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
