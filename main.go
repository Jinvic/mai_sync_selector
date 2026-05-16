package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mai_sync_selector/internal/db"
	"mai_sync_selector/internal/handler"
	"mai_sync_selector/internal/model"
	"mai_sync_selector/internal/service"

	"github.com/gin-gonic/gin"
)

const version = "v1.2.1"

func main() {
	dbPath := "sqlite.db"
	if v := os.Getenv("DB_PATH"); v != "" {
		dbPath = v
	}

	port := "8080"
	if v := os.Getenv("SERVER_PORT"); v != "" {
		port = v
	}

	db.InitDB(dbPath)
	model.InitTables(db.DB)
	cronScheduler := service.StartCron()
	service.SyncData(context.Background())

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		b, err := readEmbeddedWeb("index.html")
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", b)
	})
	router.GET("/app.css", func(c *gin.Context) {
		b, err := readEmbeddedWeb("app.css")
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "text/css; charset=utf-8", b)
	})
	router.GET("/app.js", func(c *gin.Context) {
		b, err := readEmbeddedWeb("app.js")
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "application/javascript; charset=utf-8", b)
	})

	api := router.Group("/api")
	{
		api.GET("/sync-data", handler.SyncData)
		api.POST("/select-song", handler.SelectSong)
		api.GET("/from", handler.GetFromList)
		api.GET("/genre", handler.GetGenreList)
		api.GET("/level", handler.GetLevelList)
		api.GET("/version", func(c *gin.Context) {
			handler.SuccessWithData(c, gin.H{
				"version": version,
			})
		})
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("HTTP server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	signal.Stop(sigCh)

	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown: %v", err)
	}

	cronScheduler.Stop()

	if sqlDB, err := db.DB.DB(); err == nil {
		_ = sqlDB.Close()
	}

	log.Println("exit")
}
