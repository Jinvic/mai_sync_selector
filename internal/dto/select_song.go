package dto

import "mai_sync_selector/internal/model"

type SelectSongRequest struct {
	Filter1  SelectSongFilter `json:"filter1"` // 1P
	Filter2  SelectSongFilter `json:"filter2"` // 2P
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

type SelectSongFilter struct {
	FromList  []string `json:"from_list"`
	GenreList []string `json:"genre_list"`
	MinDS     float64  `json:"min_ds"`
	MaxDS     float64  `json:"max_ds"`
	MinLevel  string   `json:"min_level"`
	MaxLevel  string   `json:"max_level"`
}

type SelectSongResponse struct {
	SongList []model.MaiMaiMusicData `json:"song_list"`
	Total    int64                   `json:"total"`
}
