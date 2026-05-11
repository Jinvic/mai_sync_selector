package model

import (
	"database/sql/driver"
	"encoding/json"
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

var DifficultyList = []string{"basic", "advanced", "expert", "master", "re_master"}
var BanquetDifficultyList = []string{"banquet", "banquet2"}

type DS map[string]float64
type Level map[string]string
type Cids map[string]int

// 将数组按难度映射为 map
func mapByDifficulty[T float64 | string | int](row []T, genre string) map[string]T {
	result := make(map[string]T)

	var difficultyList []string
	if genre == "宴会場" {
		difficultyList = BanquetDifficultyList
	} else {
		difficultyList = DifficultyList
	}
	for i := range row {
		difficulty := difficultyList[i]
		result[difficulty] = row[i]
	}
	return result
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

func (d DS) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	return json.Marshal(d)
}

func (d *DS) Scan(value interface{}) error {
	if value == nil {
		*d = make(DS)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, d)
}

func (l Level) Value() (driver.Value, error) {
	if l == nil {
		return nil, nil
	}
	return json.Marshal(l)
}

func (l *Level) Scan(value interface{}) error {
	if value == nil {
		*l = make(Level)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, l)
}

func (c Cids) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

func (c *Cids) Scan(value interface{}) error {
	if value == nil {
		*c = make(Cids)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, c)
}
