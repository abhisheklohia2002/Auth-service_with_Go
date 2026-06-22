package cloudinary_service

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	interfaces "example.com/m/internal/interface"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryUploader struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryUploader(cld *cloudinary.Cloudinary) interfaces.FileUploader {
	return &CloudinaryUploader{
		cld: cld,
	}
}

func (u *CloudinaryUploader) UploadPDF(
	ctx context.Context,
	file multipart.File,
	filename string,
	folder string,
) (*interfaces.UploadResult, error) {
	ext := filepath.Ext(filename)
	publicID := strings.TrimSuffix(filename, ext)

	result, err := u.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:       folder,
		ResourceType: "raw",
		PublicID:     publicID,
	})

	if err != nil {
		return nil, err
	}

	return &interfaces.UploadResult{
		URL:      result.SecureURL,
		PublicID: result.PublicID,
		FileType: "pdf",
	}, nil
}

func (u *CloudinaryUploader) Delete(ctx context.Context, publicID string) error {
	_, err := u.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "raw",
	})

	return err
}

func (u *CloudinaryUploader) UploadImage(
	ctx context.Context,
	file multipart.File,
	fileName string,
	folder string,
) (*interfaces.UploadResult, error) {
	if file == nil {
		return nil, errors.New("image file is required")
	}

	uploadResult, err := u.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:       folder,
		ResourceType: "image",
		PublicID:     strings.TrimSuffix(fileName, filepath.Ext(fileName)),
	})
	if err != nil {
		return nil, err
	}

	return &interfaces.UploadResult{
		URL:      uploadResult.SecureURL,
		PublicID: uploadResult.PublicID,
	}, nil
}

func (u *CloudinaryUploader) UploadVideo(
	ctx context.Context,
	file multipart.File,
	fileName string,
	folder string,
) (*interfaces.UploadResult, error) {
	result, err := u.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:       folder,
		ResourceType: "video",
		PublicID:     strings.TrimSuffix(fileName, filepath.Ext(fileName)),
	})
	if err != nil {
		return nil, err
	}

	return &interfaces.UploadResult{
		URL:      result.SecureURL,
		PublicID: result.PublicID,
	}, nil
}

func (u *CloudinaryUploader) DeleteVideo(ctx context.Context, publicID string) error {
	if publicID == "" {
		return nil
	}

	_, err := u.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "video",
	})

	return err
}
