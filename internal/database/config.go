package database

import (
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/google/uuid"

	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/service"
	"github.com/ghostsecurity/reaper/internal/types"
)

const (
	file = "reaper.db"
	// file
	url = file
	// in memory
	// url = "file::memory:?cache=shared"
)

func Connect() *gorm.DB {
	// clean up file if we're using in-memory db
	if strings.Contains(url, "memory") {
		_ = os.Remove(file)
	}
	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{
		// DEBUG
		// Logger: logger.Default.LogMode(logger.Info), // log all queries to console
	})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func Migrate() {
	db := Connect()

	_ = db.AutoMigrate(
		&models.User{},
		&models.Setting{},
		&models.Project{},
		&models.Domain{},
		&models.Host{},
		&models.Endpoint{},
		&models.Request{},
		&models.Response{},
		&models.FuzzAttack{},
		&models.FuzzResult{},
		&models.Report{},
		&models.AgentSession{},
		&models.AgentSessionMessage{},
	)

	project := models.Project{
		Name: "Default",
	}
	db.Create(&project)

	// create an agent user if one doesn't exist
	res := db.Model(&models.User{}).
		Where(&models.User{Role: types.UserRoleAgent}).
		First(&models.User{})
	if res.RowsAffected == 0 {
		_, _ = service.CreateAgentUser(db)
	}

	adminToken, _ := service.GetSettingByKey("admin_token", db)
	guestToken, _ := service.GetSettingByKey("guest_token", db)

	if adminToken == nil {
		db.Create(&models.Setting{
			Key:   "admin_token",
			Value: uuid.New().String(),
		})
	}
	if guestToken == nil {
		db.Create(&models.Setting{
			Key:   "guest_token",
			Value: uuid.New().String(),
		})
	}
}
