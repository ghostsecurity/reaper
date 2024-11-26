package database

import (
	"errors"
	"os"
	"time"

	"github.com/ghostsecurity/reaper/internal/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// Initialize sets up the database connection and runs migrations
func Initialize() error {
	var err error
	db, err = gorm.Open(sqlite.Open("reaper.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// Run migrations
	err = db.AutoMigrate(&models.OpenAIConfig{})
	if err != nil {
		return err
	}

	// Check for API key in environment and store if not exists
	return initializeOpenAIConfig()
}

// initializeOpenAIConfig checks for API key in environment and stores it if not exists
func initializeOpenAIConfig() error {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil // No API key provided, skip initialization
	}

	var config models.OpenAIConfig
	result := db.Where("is_active = ?", true).First(&config)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new config with environment API key
			config = models.OpenAIConfig{
				APIKey:   apiKey,
				Model:    "gpt-4o", // Default model updated
				LastUsed: time.Now(),
				IsActive: true,
			}
			return db.Create(&config).Error
		}
		return result.Error
	}

	return nil // Config already exists
}

// GetOpenAIConfig retrieves the active OpenAI configuration
func GetOpenAIConfig() (*models.OpenAIConfig, error) {
	var config models.OpenAIConfig
	result := db.Where("is_active = ?", true).First(&config)
	if result.Error != nil {
		return nil, result.Error
	}
	return &config, nil
}

// isValidModel checks if the provided model is supported
func isValidModel(model string) bool {
	validModels := []string{
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-4-turbo",
		"gpt-4",
		"gpt-3.5-turbo",
	}
	for _, validModel := range validModels {
		if model == validModel {
			return true
		}
	}
	return false
}

// UpdateOpenAIConfig updates the OpenAI configuration
func UpdateOpenAIConfig(apiKey string, model string) error {
	if !isValidModel(model) {
		return errors.New("invalid model specified")
	}

	// Deactivate current active config
	err := db.Model(&models.OpenAIConfig{}).Where("is_active = ?", true).Update("is_active", false).Error
	if err != nil {
		return err
	}

	// Create new active config
	config := models.OpenAIConfig{
		APIKey:   apiKey,
		Model:    model,
		LastUsed: time.Now(),
		IsActive: true,
	}
	return db.Create(&config).Error
}

// UpdateLastUsed updates the last used timestamp for the active configuration
func UpdateLastUsed() error {
	return db.Model(&models.OpenAIConfig{}).Where("is_active = ?", true).Update("last_used", time.Now()).Error
}
