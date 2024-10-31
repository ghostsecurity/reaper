package service

import (
	"github.com/ghostsecurity/reaper/internal/database/models"
	"gorm.io/gorm"
)

func GetSettingByKey(key string, db *gorm.DB) (*string, error) {
	setting := models.Setting{}
	res := db.Where(models.Setting{
		Key: key,
	}).First(&setting)

	if res.Error != nil {
		return nil, res.Error
	}

	return &setting.Value, nil
}

func SetSettingByKey(key string, value string, db *gorm.DB) error {
	setting := models.Setting{
		Key:   key,
		Value: value,
	}
	res := db.Create(&setting)

	return res.Error
}

func DeleteSettingByKey(key string, db *gorm.DB) error {
	res := db.Where(models.Setting{
		Key: key,
	}).Delete(&models.Setting{})

	return res.Error
}
