package handler

import (
	"mai_sync_selector/internal/dto"
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
