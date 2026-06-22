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
	UploadImage(
		ctx context.Context,
		file multipart.File,
		fileName string,
		folder string,
	) (*UploadResult, error)

	UploadVideo(
		ctx context.Context,
		file multipart.File,
		fileName string,
		folder string,
	) (*UploadResult, error)

	Delete(ctx context.Context, publicID string) error

	DeleteVideo(ctx context.Context, publicID string) error
}
