package handler

import (
	"log"
	"mai_sync_selector/internal/service"

	"github.com/gin-gonic/gin"
)

func SyncData(c *gin.Context) {
	err := service.SyncData(c.Request.Context())
	if err != nil {
		log.Printf("SyncData error: %v", err)
		Failed(c, err)
	}

	Success(c)
}
