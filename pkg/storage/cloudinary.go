package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/ephraimd/cloud-documents-service/internal/config"
	"github.com/ephraimd/cloud-documents-service/internal/contracts"
	"github.com/ephraimd/cloud-documents-service/pkg/logger"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryProvider implements CloudStorageProvider for Cloudinary
type CloudinaryProvider struct {
	client *cloudinary.Cloudinary
}

// NewCloudinaryProvider creates a new Cloudinary provider
func NewCloudinaryProvider() (*CloudinaryProvider, error) {
	cfg := config.GlobalConfig

	if cfg.CloudinaryCloudName == "" || cfg.CloudinaryAPIKey == "" || cfg.CloudinaryAPISecret == "" {
		return nil, fmt.Errorf("Cloudinary credentials not provided")
	}

	cld, err := cloudinary.NewFromParams(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloudinary client: %v", err)
	}

	logger.Logger.Printf("✓ Cloudinary provider initialized with cloud: %s", cfg.CloudinaryCloudName)

	return &CloudinaryProvider{
		client: cld,
	}, nil
}

// Upload uploads a file to Cloudinary
func (p *CloudinaryProvider) Upload(file io.Reader, filename, folder string, fileSize int64) (*contracts.UploadResponse, error) {
	ctx := context.Background()

	publicID := p.buildPublicID(filename, folder)

	uploadParams := uploader.UploadParams{
		PublicID: publicID,
		Folder:   folder,
	}

	result, err := p.client.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to Cloudinary: %v", err)
	}

	return &contracts.UploadResponse{
		URL:      result.SecureURL,
		Filename: filename,
		Size:     int64(result.Bytes),
		Provider: "cloudinary",
	}, nil
}

// Download downloads a file from Cloudinary (returns redirect URL)
func (p *CloudinaryProvider) Download(filename, folder string) (io.Reader, error) {
	// Cloudinary doesn't support direct download like S3
	// This would typically return a redirect to the file URL
	return nil, fmt.Errorf("direct download not supported by Cloudinary provider - use GetFileURL instead")
}

// Delete deletes a file from Cloudinary
func (p *CloudinaryProvider) Delete(filename, folder string) error {
	ctx := context.Background()

	publicID := p.buildPublicID(filename, folder)

	_, err := p.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from Cloudinary: %v", err)
	}

	return nil
}

// GetFileURL returns the public URL of a file
func (p *CloudinaryProvider) GetFileURL(filename, folder string) (string, error) {
	publicID := p.buildPublicID(filename, folder)

	asset, err := p.client.Image(publicID)
	if err != nil {
		return "", fmt.Errorf("failed to generate Cloudinary URL: %v", err)
	}

	url, err := asset.String()
	if err != nil {
		return "", fmt.Errorf("failed to convert asset to URL: %v", err)
	}

	return url, nil
}

// GetProviderName returns the provider name
func (p *CloudinaryProvider) GetProviderName() string {
	return "cloudinary"
}

// buildPublicID builds the Cloudinary public ID from filename and folder
func (p *CloudinaryProvider) buildPublicID(filename, folder string) string {
	// Remove file extension for Cloudinary public ID
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	if folder == "" {
		return name
	}
	return filepath.Join(strings.Trim(folder, "/"), name)
}
