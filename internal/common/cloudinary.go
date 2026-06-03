package common

import (
	"context"
	"log"
	"example.com/m/internal/config"
	"github.com/cloudinary/cloudinary-go/v2"
)

func InitCloudinary() *cloudinary.Cloudinary {
	cfg := config.LoadDotenv()
	cld, err := cloudinary.NewFromParams(cfg.CLOUDINARY_CLOUD_NAME, cfg.CLOUDINARY_API_KEY, cfg.CLOUDINARY_API_SECRET)
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary: ", err)
	}

	_ = context.Background()

	return cld
}