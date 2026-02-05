package repository

import (
	"errors"

	"anyrouter-checkin/internal/model"

	"gorm.io/gorm"
)

func ListConfigs(category string) ([]model.Config, error) {
	var configs []model.Config
	if err := DB.Where("category = ?", category).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

func GetConfigValue(key string) (string, error) {
	var cfg model.Config
	if err := DB.Where("`key` = ?", key).First(&cfg).Error; err != nil {
		return "", err
	}
	return cfg.Value, nil
}

func SetConfigValue(key, value, category string) error {
	var cfg model.Config
	result := DB.Where("`key` = ?", key).First(&cfg)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			cfg = model.Config{Key: key, Value: value, Category: category}
			return DB.Create(&cfg).Error
		}
		return result.Error
	}
	cfg.Value = value
	return DB.Save(&cfg).Error
}
