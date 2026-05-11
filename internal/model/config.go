package model

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Config struct {
	Key   string `json:"key" gorm:"primaryKey;index"`
	Value string `json:"value"`
}

func GetEtag(tx *gorm.DB) (string, error) {
	var config Config
	if err := tx.Where("key = ?", "etag").First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return config.Value, nil
}

func SetEtag(tx *gorm.DB, etag string) error {
	return tx.Model(&Config{}).
		Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns([]string{"value"})}).
		Create(&Config{Key: "etag", Value: etag}).Error
}
