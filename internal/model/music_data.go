package model

import (
	"mai_sync_selector/divingfish"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MaiMaiMusicData struct {
	// 主键
	ID string `gorm:"primaryKey;type:varchar(100);not null;comment:音乐ID" json:"id"`

	// 基础信息
	Title    string `gorm:"type:varchar(500);index:idx_title;comment:标题" json:"title"`
	Type     string `gorm:"type:varchar(50);index:idx_type;default:'';comment:类型" json:"type"`
	Artist   string `gorm:"type:varchar(200);index:idx_artist;comment:艺术家" json:"artist"`
	Genre    string `gorm:"type:varchar(100);index:idx_genre;comment:流派" json:"genre"`
	From     string `gorm:"type:varchar(50);index:idx_from;default:'';comment:来源" json:"from"`
	CoverURL string `gorm:"type:varchar(500);index:idx_cover_url;comment:封面URL" json:"cover_url"`

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
	m.CoverURL = divingfish.GetCoverUrl(musicData.ID)
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
func BatchUpsert(tx *gorm.DB, dataList []MaiMaiMusicData) error {
	if len(dataList) == 0 {
		return nil
	}

	return tx.Transaction(func(tx *gorm.DB) error {
		for i := range dataList {
			// 使用 Save 自动判断插入或更新
			if err := tx.Save(&dataList[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

type DsFilter struct {
	MinDS float64
	MaxDS float64
}

type LevelFilter struct {
	LevelRange []string
}

// 根据过滤条件查询歌曲
func SelectSongByFilter(tx *gorm.DB,
	fromList []string,
	genreList []string,
	dsFilters []DsFilter,
	levelFilters []LevelFilter,
	page int,
	pageSize int) (musicDataList []MaiMaiMusicData, total int64, err error) {
	query := tx.Debug().Model(&MaiMaiMusicData{})

	// 宴会場不参与筛选
	query = query.Where("genre != '宴会場'")

	if len(fromList) > 0 {
		query = query.Where("from IN (?)", fromList)
	}
	if len(genreList) > 0 {
		query = query.Where("genre IN (?)", genreList)
	}

	for _, dsFilter := range dsFilters {
		parts := make([]string, 0, len(DifficultyList))
		args := make([]interface{}, 0, len(DifficultyList)*2)
		for _, difficulty := range DifficultyList {
			parts = append(parts, "(ds->>? >= ? AND ds->>? <= ?) ")
			args = append(args, difficulty, dsFilter.MinDS, difficulty, dsFilter.MaxDS)
		}
		query = query.Where("("+strings.Join(parts, " OR ")+")", args...)
	}

	for _, levelFilter := range levelFilters {
		parts := make([]string, 0, len(DifficultyList))
		args := make([]interface{}, 0, len(DifficultyList)*2)
		for _, difficulty := range DifficultyList {
			parts = append(parts, "level->>? IN (?) ")
			args = append(args, difficulty, levelFilter.LevelRange)
		}
		query = query.Where("("+strings.Join(parts, " OR ")+")", args...)
	}

	err = query.Session(&gorm.Session{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&musicDataList)
	return musicDataList, total, nil
}

// ------------------------------------------------------------

type FromList struct {
	From      string    `json:"from" gorm:"primaryKey;uniqueIndex"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (FromList) TableName() string {
	return "from_list"
}

type GenreList struct {
	Genre     string    `json:"genre" gorm:"primaryKey;uniqueIndex"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (GenreList) TableName() string {
	return "genre_list"
}

func GetFromList(tx *gorm.DB) ([]string, error) {
	var fromList []string
	if err := tx.Model(&FromList{}).Pluck("from", &fromList).Error; err != nil {
		return nil, err
	}
	return fromList, nil
}

func InsertFromListIfNotExists(tx *gorm.DB, fromList []FromList) error {
	return tx.Model(&FromList{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&fromList).Error
}

func GetGenreList(tx *gorm.DB) ([]string, error) {
	var genreList []string
	if err := tx.Model(&GenreList{}).Pluck("genre", &genreList).Error; err != nil {
		return nil, err
	}
	return genreList, nil
}

func InsertGenreListIfNotExists(tx *gorm.DB, genreList []GenreList) error {
	return tx.Model(&GenreList{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&genreList).Error
}
