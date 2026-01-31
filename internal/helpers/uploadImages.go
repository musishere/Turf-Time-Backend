package helpers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/config"
)

const maxFileSize = 5 * 1024 * 1024 // 5MB

var (
	allowedImageExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	safePublicIDRe   = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
)

type ImageUploader struct {
	cld           *cloudinary.Cloudinary
	cloudName     string
	defaultFolder string
}

func NewImageUploader() (*ImageUploader, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("cloudinary credentials not configured")
	}

	cfg, err := config.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("cloudinary config: %w", err)
	}
	cfg.API.Timeout = 120
	cfg.API.UploadTimeout = 120

	cld, err := cloudinary.NewFromConfiguration(*cfg)
	if err != nil {
		return nil, fmt.Errorf("cloudinary init: %w", err)
	}

	return &ImageUploader{
		cld:           cld,
		cloudName:     cloudName,
		defaultFolder: "sports_images",
	}, nil
}

func sanitizePublicID(name string) string {
	s := strings.TrimSpace(name)
	s = strings.ReplaceAll(s, " ", "_")
	s = safePublicIDRe.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	if s == "" {
		return "image"
	}
	return s
}

func (u *ImageUploader) makePublicID(filename string) string {
	base := sanitizePublicID(strings.TrimSuffix(filename, filepath.Ext(filename)))
	return fmt.Sprintf("%s_%d", base, time.Now().Unix())
}

// urlFromResult returns the secure URL from an upload result, or an error.
func (u *ImageUploader) urlFromResult(result *uploader.UploadResult) (string, error) {
	if result == nil {
		return "", fmt.Errorf("cloudinary returned empty response")
	}
	if result.Error.Message != "" {
		return "", fmt.Errorf("cloudinary: %s", result.Error.Message)
	}

	url := result.SecureURL
	if url == "" && result.PublicID != "" {
		rt := result.ResourceType
		if rt == "" {
			rt = "image"
		}
		format := result.Format
		if format == "" {
			format = "jpg"
		}
		url = fmt.Sprintf("https://res.cloudinary.com/%s/%s/upload/v%d/%s.%s",
			u.cloudName, rt, result.Version, result.PublicID, format)
	}
	if url == "" {
		return "", fmt.Errorf("cloudinary returned no URL")
	}
	return url, nil
}

func (u *ImageUploader) upload(ctx context.Context, reader interface{}, folder, publicID, transformation string) (string, error) {
	result, err := u.cld.Upload.Upload(ctx, reader, uploader.UploadParams{
		Folder:         folder,
		PublicID:       publicID,
		ResourceType:   "image",
		Transformation: transformation,
	})
	if err != nil {
		return "", fmt.Errorf("cloudinary: %w", err)
	}
	return u.urlFromResult(result)
}

func (u *ImageUploader) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size > maxFileSize {
		return "", fmt.Errorf("file too large: max 5MB allowed")
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedImageExts[ext] {
		return "", fmt.Errorf("icon must be an image (jpg, jpeg, png, webp); got %s", ext)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	publicID := u.makePublicID(fileHeader.Filename)
	url, err := u.upload(ctx, file, u.defaultFolder, publicID, "q_auto,f_auto,w_1000")
	if err != nil {
		log.Printf("[Cloudinary] %s: %v", fileHeader.Filename, err)
		return "", err
	}
	return url, nil
}

func (u *ImageUploader) UploadImageWithFolder(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (string, error) {
	orig := u.defaultFolder
	u.defaultFolder = folder
	defer func() { u.defaultFolder = orig }()
	return u.UploadImage(ctx, fileHeader)
}

func (u *ImageUploader) UploadFromBytes(ctx context.Context, data []byte, filename, folder string) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty file data")
	}
	if len(data) > maxFileSize {
		return "", fmt.Errorf("file too large: max 5MB allowed")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedImageExts[ext] {
		return "", fmt.Errorf("icon must be an image (jpg, jpeg, png, webp); got %s", ext)
	}

	publicID := u.makePublicID(filename)
	url, err := u.upload(ctx, bytes.NewReader(data), folder, publicID, "q_auto,f_auto")
	if err != nil {
		log.Printf("[Cloudinary] %s: %v", filename, err)
		return "", err
	}
	return url, nil
}
