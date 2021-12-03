package domain

import (
	"context"
)

// Uplod ...
type Upload struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	FileName  string `json:"file_name"`
	Headers   string `json:"headers"`
	Size      int64  `json:"size"`
	Extension string `json:"extension"`
}

type UploadEntity interface {
	Upload(ctx context.Context, upl Upload) error
}

type UploadRepository interface {
	Upload(ctx context.Context, upl Upload) error
}
