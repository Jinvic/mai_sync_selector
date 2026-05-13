package handler

import (
	"mai_sync_selector/internal/db"
	"mai_sync_selector/internal/dto"
	"mai_sync_selector/internal/model"
	"mai_sync_selector/internal/service"

	"github.com/gin-gonic/gin"
)

func SelectSong(c *gin.Context) {
	var request dto.SelectSongRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		Failed(c, err)
		return
	}

	response, err := service.SelectSong(request)
	if err != nil {
		Failed(c, err)
		return
	}
	SuccessWithData(c, response)
}

func GetFromList(c *gin.Context) {
	fromList, err := model.GetFromList(db.DB)
	if err != nil {
		Failed(c, err)
		return
	}
	SuccessWithData(c, gin.H{
		"from_list": fromList,
	})
}

func GetGenreList(c *gin.Context) {
	genreList, err := model.GetGenreList(db.DB)
	if err != nil {
		Failed(c, err)
		return
	}

	// 去掉宴会場
	for i, genre := range genreList {
		if genre == "宴会場" {
			genreList = append(genreList[:i], genreList[i+1:]...)
			break
		}
	}

	SuccessWithData(c, gin.H{
		"genre_list": genreList,
	})
}

func GetLevelList(c *gin.Context) {
	SuccessWithData(c, gin.H{
		"level_list": model.LevelList,
	})
}
