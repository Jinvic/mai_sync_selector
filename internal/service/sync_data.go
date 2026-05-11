package service

import (
	"context"
	"mai_sync_selector/divingfish"
	"mai_sync_selector/internal/db"
	"mai_sync_selector/internal/model"
)

func SyncData(ctx context.Context) error {
	etag, err := model.GetEtag(db.DB)
	if err != nil {
		return err
	}
	divingfishClient := divingfish.NewClient(etag)
	divingfishMusicDataList, etag, err := divingfishClient.GetMaiMaiMusicData(ctx)
	if err != nil {
		return err
	}
	err = model.SetEtag(db.DB, etag)
	if err != nil {
		return err
	}
	maiMaiMusicDataList := model.BatchFromDivingFish(divingfishMusicDataList)
	err = model.BatchUpsert(db.DB, maiMaiMusicDataList)
	if err != nil {
		return err
	}
	return nil
}
