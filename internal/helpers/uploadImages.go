package helpers

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/config"
)

// ImageUploader handles image uploads to Cloudinary
type ImageUploader struct {
	cld           *cloudinary.Cloudinary
	defaultFolder string
}

// NewImageUploader initializes the uploader from environment variables
func NewImageUploader() (*ImageUploader, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("cloudinary credentials not configured")
	}

	cfg, err := config.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to init cloudinary config: %w", err)
	}
	// SDK default timeout is 60s; use 120s for slow networks / large uploads
	cfg.API.Timeout = 120
	cfg.API.UploadTimeout = 120

	cld, err := cloudinary.NewFromConfiguration(*cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init cloudinary: %w", err)
	}

	return &ImageUploader{
		cld:           cld,
		defaultFolder: "sports_images", // default folder
	}, nil
}

// UploadImage takes a multipart file and returns the Cloudinary URL
func (u *ImageUploader) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	// Validation 1: Check file size (max 5MB)
	if fileHeader.Size > 5*1024*1024 {
		return "", fmt.Errorf("file too large: max 5MB allowed")
	}

	// Validation 2: Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		return "", fmt.Errorf("invalid file type: %s, allowed: jpg, jpeg, png, webp", ext)
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	// Generate unique public ID (filename without ext + timestamp)
	fileName := strings.TrimSuffix(fileHeader.Filename, ext)
	timestamp := time.Now().Unix()
	publicID := fmt.Sprintf("%s_%d", fileName, timestamp)

	// Upload parameters
	uploadParams := uploader.UploadParams{
		Folder:       u.defaultFolder,
		PublicID:     publicID,
		ResourceType: "image",
		// Overwrite:      false,                  // Prevent accidental overwrites
		Transformation: "q_auto,f_auto,w_1000", // Auto-optimize quality, format, max width 1000px
	}

	// Execute upload with context (respects timeout/cancellation)
	result, err := u.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("cloudinary upload failed: %w", err)
	}

	// Return HTTPS URL
	return result.SecureURL, nil
}

// UploadImageWithFolder allows custom folder (e.g., "sports/icons", "users/avatars")
func (u *ImageUploader) UploadImageWithFolder(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Temporarily change folder
	originalFolder := u.defaultFolder
	u.defaultFolder = folder
	defer func() { u.defaultFolder = originalFolder }()

	return u.UploadImage(ctx, fileHeader)
}

// UPLOADER - Add this method
func (u *ImageUploader) UploadFromBytes(ctx context.Context, data []byte, filename, folder string) (string, error) {
	publicID := fmt.Sprintf("%s_%d", strings.TrimSuffix(filename, filepath.Ext(filename)), time.Now().Unix())

	result, err := u.cld.Upload.Upload(ctx, bytes.NewReader(data), uploader.UploadParams{
		Folder:         folder,
		PublicID:       publicID,
		ResourceType:   "image",
		Transformation: "q_auto,f_auto",
	})

	if err != nil {
		return "", err
	}
	return result.SecureURL, nil
}
