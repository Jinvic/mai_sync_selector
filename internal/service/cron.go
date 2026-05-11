package service

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"
)

func StartCron() *cron.Cron {
	c := cron.New()
	// 每周一凌晨3点
	c.AddFunc("0 3 * * 1", func() {
		err := SyncData(context.Background())
		if err != nil {
			log.Println("SyncData error:", err)
		}
	})
	c.Start()
	return c
}
