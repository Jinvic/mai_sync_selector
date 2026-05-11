package handler

import (
	"log"
	"mai_sync_selector/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SyncData(c *gin.Context) {
	err := service.SyncData(c.Request.Context())
	if err != nil {
		log.Printf("SyncData error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
