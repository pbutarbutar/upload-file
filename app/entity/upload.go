package entity

import (
	"context"

	"github.com/pbutarbutar/upload-file/domain"
)

type UploadEntity struct {
	uploadRepo domain.UploadRepository
}

// NewArticleUsecase will create new an articleUsecase object representation of domain.ArticleUsecase interface
func NewUploadEntity(a domain.UploadRepository) domain.UploadEntity {
	return &UploadEntity{
		uploadRepo: a,
	}
}
func (a *UploadEntity) Upload(c context.Context, up domain.Upload) error {
	err := a.uploadRepo.Upload(c, up)
	if err != nil {
		return err
	}
	return nil
}
