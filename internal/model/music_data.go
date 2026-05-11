package model

import (
	"mai_sync_selector/divingfish"
	"time"

	"gorm.io/gorm"
)

type MaiMaiMusicData struct {
	// 主键
	ID string `gorm:"primaryKey;type:varchar(100);not null;comment:音乐ID" json:"id"`

	// 基础信息
	Title  string `gorm:"type:varchar(500);index:idx_title;comment:标题" json:"title"`
	Type   string `gorm:"type:varchar(50);index:idx_type;default:'';comment:类型" json:"type"`
	Artist string `gorm:"type:varchar(200);index:idx_artist;comment:艺术家" json:"artist"`
	Genre  string `gorm:"type:varchar(100);index:idx_genre;comment:流派" json:"genre"`
	From   string `gorm:"type:varchar(50);index:idx_from;default:'';comment:来源" json:"from"`

	// JSON 存储字段
	DS    DS    `gorm:"type:text;comment:难度定数数据" json:"ds"`
	Level Level `gorm:"type:text;comment:难度等级" json:"level"`
	Cids  Cids  `gorm:"type:text;comment:谱面ID" json:"cids"`

	// 数值字段
	BPM   int  `gorm:"index:idx_bpm;comment:每分钟节拍数" json:"bpm"`
	IsNew bool `gorm:"index:idx_is_new;default:false;comment:是否新曲" json:"is_new"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (MaiMaiMusicData) TableName() string {
	return "music_data"
}

func (m *MaiMaiMusicData) FromDivingFish(musicData divingfish.MaiMaiMusicData) {
	m.ID = musicData.ID
	m.Title = musicData.Title
	m.Type = musicData.Type
	m.DS = mapByDifficulty(musicData.DS, musicData.BasicInfo.Genre)
	m.Level = mapByDifficulty(musicData.Level, musicData.BasicInfo.Genre)
	m.Cids = mapByDifficulty(musicData.Cids, musicData.BasicInfo.Genre)
	m.Artist = musicData.BasicInfo.Artist
	m.Genre = musicData.BasicInfo.Genre
	m.BPM = musicData.BasicInfo.BPM
	m.From = musicData.BasicInfo.From
	m.IsNew = musicData.BasicInfo.IsNew
}

func BatchFromDivingFish(musicDataList []divingfish.MaiMaiMusicData) []MaiMaiMusicData {
	newList := make([]MaiMaiMusicData, len(musicDataList))
	for i, musicData := range musicDataList {
		newList[i].FromDivingFish(musicData)
	}
	return newList
}

// BatchUpsert 批量插入或更新
func BatchUpsert(db *gorm.DB, dataList []MaiMaiMusicData) error {
	if len(dataList) == 0 {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for i := range dataList {
			// 使用 Save 自动判断插入或更新
			if err := tx.Save(&dataList[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
