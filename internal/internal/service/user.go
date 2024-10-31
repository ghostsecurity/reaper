package service

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"

	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/types"
	"gorm.io/gorm"
)

func GetUserByToken(token string, db *gorm.DB) (*models.User, error) {
	user := models.User{}
	res := db.Where(models.User{
		Token: token,
	}).First(&user)

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func FirstAdmin(db *gorm.DB) (*models.User, error) {
	user := models.User{}

	res := db.Where(models.User{
		Role: types.UserRoleAdmin,
	}).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func CreateAgentUser(db *gorm.DB) (*models.User, error) {
	user := models.User{
		Name:  "Reaper AI Agent",
		Role:  types.UserRoleAgent,
		Token: generateUserAuthToken(),
	}

	res := db.Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func GetAgentToken(db *gorm.DB) (string, error) {
	user := models.User{}
	res := db.Where(models.User{
		Role: types.UserRoleAgent,
	}).First(&user)
	if res.Error != nil {
		return "", res.Error
	}
	return user.Token, nil
}

func CreateAdminUser(username string, db *gorm.DB) (*models.User, error) {
	user := models.User{
		Name:  sanitizeUsername(username),
		Role:  types.UserRoleAdmin,
		Token: generateUserAuthToken(),
	}

	res := db.Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func CreateGuestUser(username string, db *gorm.DB) (*models.User, error) {
	user := models.User{
		Name:  sanitizeUsername(username),
		Role:  types.UserRoleViewer,
		Token: generateUserAuthToken(),
	}

	res := db.Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func sanitizeUsername(username string) string {
	// regex to strip all non-alphanumeric characters, plus hyphens, underscores, and dots
	username = regexp.MustCompile(`[^a-zA-Z0-9-_.]+`).ReplaceAllString(username, "")

	if len(username) > 20 {
		return username[:20]
	}

	return username
}

func generateUserAuthToken() string {
	token := make([]byte, 20)
	_, err := rand.Read(token)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(token)
}
