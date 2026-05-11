package main

import (
	"mai_sync_selector/internal/db"
	"mai_sync_selector/internal/model"
	"mai_sync_selector/internal/service"
)

func main() {
	db.InitDB()
	model.InitTables(db.DB)
	service.StartCron()
}
