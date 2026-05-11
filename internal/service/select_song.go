package service

import (
	"mai_sync_selector/internal/db"
	"mai_sync_selector/internal/dto"
	"mai_sync_selector/internal/model"
	"math"
)

func SelectSong(request dto.SelectSongRequest) (dto.SelectSongResponse, error) {
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = 10
	}

	filter := combineFilter(request.Filter1, request.Filter2)
	levelList := model.GetLevelList(filter.MinLevel, filter.MaxLevel)
	musicDataList, total, err := model.SelectSongByFilter(db.DB, filter.FromList, filter.GenreList, filter.MinDS, filter.MaxDS, levelList, request.Page, request.PageSize)
	if err != nil {
		return dto.SelectSongResponse{}, err
	}
	return dto.SelectSongResponse{
		SongList: musicDataList,
		Total:    total,
	}, nil
}

func combineFilter(filter1 dto.SelectSongFilter, filter2 dto.SelectSongFilter) dto.SelectSongFilter {
	var minLevel, maxLevel string
	// 取交集：min取两者较大值，max取两者较小值
	if model.LevelSortFunc(filter1.MinLevel, filter2.MinLevel) < 0 {
		minLevel = filter2.MinLevel // filter2的min更大
	} else {
		minLevel = filter1.MinLevel // filter1的min更大或相等
	}
	if model.LevelSortFunc(filter1.MaxLevel, filter2.MaxLevel) > 0 {
		maxLevel = filter2.MaxLevel // filter2的max更小
	} else {
		maxLevel = filter1.MaxLevel // filter1的max更小或相等
	}

	// 检查交集是否有效
	if model.LevelSortFunc(minLevel, maxLevel) > 0 {
		// 交集无效，置空
		minLevel = ""
		maxLevel = ""
	}

	// 版本和曲包去重并集
	fromSet := make(map[string]bool)
	for _, f := range append(filter1.FromList, filter2.FromList...) {
		fromSet[f] = true
	}
	fromList := make([]string, 0, len(fromSet))
	for f := range fromSet {
		fromList = append(fromList, f)
	}

	genreSet := make(map[string]bool)
	for _, g := range append(filter1.GenreList, filter2.GenreList...) {
		genreSet[g] = true
	}
	genreList := make([]string, 0, len(genreSet))
	for g := range genreSet {
		genreList = append(genreList, g)
	}

	return dto.SelectSongFilter{
		// 版本和曲包取并集
		FromList:  append(filter1.FromList, filter2.FromList...),
		GenreList: append(filter1.GenreList, filter2.GenreList...),
		// 定数等级取交集
		MinDS:    math.Max(filter1.MinDS, filter2.MinDS),
		MaxDS:    math.Min(filter1.MaxDS, filter2.MaxDS),
		MinLevel: minLevel,
		MaxLevel: maxLevel,
	}
}
