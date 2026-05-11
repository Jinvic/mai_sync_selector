package model

import "gorm.io/gorm"

func InitTables(db *gorm.DB) {
	db.AutoMigrate(&MaiMaiMusicData{})
	db.AutoMigrate(&Config{})
}
