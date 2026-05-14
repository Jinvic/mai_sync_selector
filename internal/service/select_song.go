package service

import (
	"mai_sync_selector/internal/db"
	"mai_sync_selector/internal/dto"
	"mai_sync_selector/internal/model"
)

func SelectSong(request dto.SelectSongRequest) (dto.SelectSongResponse, error) {
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = 10
	}

	formList := UnionList(request.Filter1.FromList, request.Filter2.FromList)
	genreList := UnionList(request.Filter1.GenreList, request.Filter2.GenreList)
	dsFilters := []model.DsFilter{
		{MinDS: request.Filter1.MinDS, MaxDS: request.Filter1.MaxDS},
		{MinDS: request.Filter2.MinDS, MaxDS: request.Filter2.MaxDS},
	}
	levelFilters := []model.LevelFilter{
		{LevelRange: model.GetLevelRange(request.Filter1.MinLevel, request.Filter1.MaxLevel)},
		{LevelRange: model.GetLevelRange(request.Filter2.MinLevel, request.Filter2.MaxLevel)},
	}
	musicDataList, total, err := model.SelectSongByFilter(db.DB, formList, genreList, dsFilters, levelFilters, request.Page, request.PageSize)
	if err != nil {
		return dto.SelectSongResponse{}, err
	}
	return dto.SelectSongResponse{
		SongList: musicDataList,
		Total:    total,
	}, nil
}

func UnionList[T comparable](list1 []T, list2 []T) []T {
	set := make(map[T]bool)
	for _, item := range list1 {
		set[item] = true
	}
	for _, item := range list2 {
		set[item] = true
	}

	result := make([]T, 0, len(set))
	for item := range set {
		result = append(result, item)
	}
	return result
}