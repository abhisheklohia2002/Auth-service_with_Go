package interfaces

import (
	"context"
	"mime/multipart"
)

type UploadResult struct {
	URL      string
	PublicID string
	FileType string
	FileSize int64
}

type FileUploader interface {
	UploadPDF(ctx context.Context, file multipart.File, filename string, folder string) (*UploadResult, error)
	Delete(ctx context.Context, publicID string) error
}