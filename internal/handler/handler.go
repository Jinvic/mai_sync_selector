package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Failed(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status": 0,
		"error":  err.Error(),
	})
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": 1,
	})
}

func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   data,
	})
}
