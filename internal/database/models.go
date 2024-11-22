package database

import (
	"time"
)

// OpenAIConfig stores OpenAI API configuration
type OpenAIConfig struct {
	ID        uint      `gorm:"primarykey"`
	APIKey    string    `gorm:"uniqueIndex;not null"`
	Model     string    `gorm:"not null;default:'gpt-4-turbo-preview'"`
	LastUsed  time.Time `gorm:"index"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
}
