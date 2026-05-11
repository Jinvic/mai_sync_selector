package service

import (
	"context"
	"mai_sync_selector/divingfish"
	"mai_sync_selector/internal/db"
	"mai_sync_selector/internal/model"

	"gorm.io/gorm"
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
	maiMaiMusicDataList := model.BatchFromDivingFish(divingfishMusicDataList)

	fromListSet := make(map[string]bool)
	genreListSet := make(map[string]bool)
	for _, musicData := range maiMaiMusicDataList {
		fromListSet[musicData.From] = true
		genreListSet[musicData.Genre] = true
	}

	fromList := make([]model.FromList, 0, len(fromListSet))
	genreList := make([]model.GenreList, 0, len(genreListSet))
	for from := range fromListSet {
		fromList = append(fromList, model.FromList{From: from})
	}
	for genre := range genreListSet {
		genreList = append(genreList, model.GenreList{Genre: genre})
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		err = model.SetEtag(tx, etag)
		if err != nil {
			return err
		}
		err = model.BatchUpsert(tx, maiMaiMusicDataList)
		if err != nil {
			return err
		}
		err = model.InsertFromListIfNotExists(tx, fromList)
		if err != nil {
			return err
		}
		err = model.InsertGenreListIfNotExists(tx, genreList)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
