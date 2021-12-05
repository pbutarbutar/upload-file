package repository

import (
	"context"

	"github.com/pbutarbutar/upload-file/domain"

	"github.com/jinzhu/gorm"
)

type UploadRepository struct {
	Conn *gorm.DB
}

//set DB connector in params
func NewUploadRepository(Conn *gorm.DB) domain.UploadRepository {
	return &UploadRepository{Conn}
}

func (u UploadRepository) Close() {
	u.Conn.Close()
}

func (u *UploadRepository) Upload(ctx context.Context, upld domain.Upload) error {
	err := u.Conn.Create(&upld).Error
	return err
}
