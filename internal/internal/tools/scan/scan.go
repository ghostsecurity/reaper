package scan

import (
	"errors"
	"log/slog"

	"gorm.io/gorm"

	"github.com/ghostsecurity/reaper/internal/database/models"
)

// updateDomainHostCount updates the host count for a domain
func updateDomainHostCount(domain *models.Domain, db *gorm.DB) error {
	res := db.First(&domain)
	if res.RowsAffected == 0 {
		slog.Error("[subfinder]", "message", "failed to update domain host count", "domain", domain.Name, "error", "domain not found")
		return errors.New("domain not found")
	}
	count := int64(0)

	resp := db.
		Model(&models.Host{}).
		Where(&models.Host{DomainID: domain.ID}).
		Count(&count)
	if resp.Error != nil {
		slog.Error("[subfinder]", "message", "failed to update domain host count", "domain", domain.Name, "error", resp.Error)
		return resp.Error
	}
	domain.HostCount = int(count)
	db.Model(&domain).Updates(domain)
	return nil
}
