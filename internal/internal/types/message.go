package types

import (
	"time"
)

type MessageType string

const (
	MessageTypeExploreHost         MessageType = "explore_host"
	MessageTypeExploreEndpoint     MessageType = "explore_endpoint"
	MessageTypeScanDomain          MessageType = "scan.domain"
	MessageTypeScanDomainResult    MessageType = "scan.domain.result"
	MessageTypeScanSyncDomain      MessageType = "scan.domain.sync"
	MessageTypeNavigationFollow    MessageType = "navigation.follow"
	MessageTypeExploreRequest      MessageType = "explore.request"
	MessageTypeExploreResponse     MessageType = "explore.response"
	MessageTypeAttackResult        MessageType = "attack.result"
	MessageTypeAttackResultClear   MessageType = "attack.result.clear"
	MessageTypeAttackComplete      MessageType = "attack.complete"
	MessageTypeReportStatus        MessageType = "report.status"
	MessageTypeReportSync          MessageType = "report.sync"
	MessageTypeAgentSessionMessage MessageType = "agent.session.message"
)

type ExploreHostMessage struct {
	Type      MessageType `json:"type"`
	Name      string      `json:"name"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
}

func (h *ExploreHostMessage) GetType() MessageType {
	return h.Type
}

type ExploreEndpointMessage struct {
	Type      MessageType `json:"type"`
	Path      string      `json:"path"`
	Host      string      `json:"host"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
}

func (h *ExploreEndpointMessage) GetType() MessageType {
	return h.Type
}

type ProxyMessage struct {
	Type   MessageType `json:"type"`
	Host   string      `json:"host"`
	Method string      `json:"method"`
	Path   string      `json:"path"`
	Status int         `json:"status"`
}

func (p *ProxyMessage) GetType() MessageType {
	return p.Type
}

type DomainMessage struct {
	Type      MessageType `json:"type"`
	Domain    string      `json:"domain"`
	Source    string      `json:"source"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
}

func (d *DomainMessage) GetType() MessageType {
	return d.Type
}

type HostMessage struct {
	Type      MessageType `json:"type"`
	Host      string      `json:"host"`
	IP        string      `json:"ip"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
}

func (h *HostMessage) GetType() MessageType {
	return h.Type
}

type SubfinderResultMessage struct {
	Type   MessageType `json:"type"`
	Domain string      `json:"domain"`
	Host   string      `json:"host"`
	Source string      `json:"source"`
}

func (m *SubfinderResultMessage) GetType() MessageType {
	return m.Type
}

type AttackCompleteMessage struct {
	Type MessageType `json:"type"`
}

func (m *AttackCompleteMessage) GetType() MessageType {
	return m.Type
}

type AttackResultMessage struct {
	Type             MessageType `json:"type"`
	EndpointID       uint        `json:"endpoint_id"` // reaper endpoint id
	TemplateName     string      `json:"template_name"`
	TemplateAuthor   string      `json:"template_author"`
	TemplateSeverity string      `json:"template_severity"`
	TemplateTags     string      `json:"template_tags"`
	TemplateType     string      `json:"template_type"`
	Hostname         string      `json:"hostname"`
	Port             string      `json:"port"`
	Scheme           string      `json:"scheme"`
	URL              string      `json:"url"`
	Endpoint         string      `json:"endpoint"`
	Request          string      `json:"request"`
	Response         string      `json:"response"`
	IpAddress        string      `json:"ip_address"`
	Command          string      `json:"command"` // curl command to reproduce the request
	Timestamp        time.Time   `json:"timestamp,omitempty"`
}

func (a *AttackResultMessage) GetType() MessageType {
	return a.Type
}

type DomainSyncMessage struct {
	Type          MessageType  `json:"type"`
	ID            uint         `json:"id"`
	Status        DomainStatus `json:"status"`
	HostCount     int          `json:"host_count"`
	LastScannedAt *time.Time   `json:"last_scanned_at"`
}

func (m *DomainSyncMessage) GetType() MessageType {
	return m.Type
}

type NavigationFollowMessage struct {
	Type      MessageType `json:"type"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
}

func (m *NavigationFollowMessage) GetType() MessageType {
	return m.Type
}

type ReportMessage struct {
	Type      MessageType `json:"type"`
	Domain    string      `json:"domain"`
	Markdown  string      `json:"markdown"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
}

func (m *ReportMessage) GetType() MessageType {
	return m.Type
}

type AgentSessionMessage struct {
	Type       MessageType `json:"type"`
	SessionID  uint        `json:"session_id"`
	AuthorID   uint        `json:"author_id"`
	AuthorRole UserRole    `json:"author_role"`
	Content    string      `json:"content"`
	Timestamp  time.Time   `json:"timestamp,omitempty"`
}

func (m *AgentSessionMessage) GetType() MessageType {
	return m.Type
}
