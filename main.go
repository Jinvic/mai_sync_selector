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

	router := gin.New()
	gin.SetMode(gin.ReleaseMode)
	router.Use(gin.Logger(), gin.Recovery())
	router.GET("/sync-data", handler.SyncData)
	router.POST("/select-song", handler.SelectSong)
	router.GET("/from", handler.GetFromList)
	router.GET("/genre", handler.GetGenreList)
	router.GET("/level", handler.GetLevelList)
	router.GET("/", func(c *gin.Context) {
		c.File("select.html")
	})

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
