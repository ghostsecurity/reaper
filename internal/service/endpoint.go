package service

import (
	"github.com/ghostsecurity/reaper/internal/database/models"
	"gorm.io/gorm"
)

type EndpointInput struct {
	Method   string `json:"method"`
	Hostname string `json:"hostname"`
	Path     string `json:"path"`
	Params   string `json:"params"`
}

// CreateOrUpdateEndpoint creates a new endpoint or updates an existing one with params
func CreateOrUpdateEndpoint(db *gorm.DB, ep EndpointInput) (models.Endpoint, error) {
	endpoint := models.Endpoint{
		Hostname: ep.Hostname,
		Method:   ep.Method,
		Path:     ep.Path,
		Params:   ep.Params,
	}

	resp := db.Where(models.Endpoint{
		Hostname: ep.Hostname,
		Path:     ep.Path,
		Method:   ep.Method,
	}).First(&endpoint)
	if resp.RowsAffected == 0 {
		db.Create(&endpoint)
	}
	if resp.RowsAffected == 1 {
		endpoint.Params = ep.Params
		db.Save(&endpoint)
	}

	return endpoint, nil
}
