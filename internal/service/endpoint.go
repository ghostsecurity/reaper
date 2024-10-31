package service

import (
	"github.com/ghostsecurity/reaper/internal/database/models"
	"gorm.io/gorm"
)

type EndpointInput struct {
	Method   string `json:"method"`
	Hostname string `json:"hostname"`
	Path     string `json:"path"`
}

func CreateEndpoint(db *gorm.DB, ep EndpointInput) (models.Endpoint, error) {
	endpoint := models.Endpoint{
		Hostname: ep.Hostname,
		Method:   ep.Method,
		Path:     ep.Path,
	}

	resp := db.Where(models.Endpoint{
		Hostname: ep.Hostname,
		Path:     ep.Path,
		Method:   ep.Method,
	}).FirstOrCreate(&endpoint)
	if resp.Error != nil {
		return models.Endpoint{}, resp.Error
	}

	return endpoint, nil
}
