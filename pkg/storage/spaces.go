package storage

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/ephraimd/cloud-documents-service/internal/config"
	"github.com/ephraimd/cloud-documents-service/internal/contracts"
	"github.com/ephraimd/cloud-documents-service/pkg/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// SpacesProvider implements CloudStorageProvider for DigitalOcean Spaces
type SpacesProvider struct {
	client   *s3.S3
	uploader *s3manager.Uploader
	bucket   string
	endpoint string
	region   string
}

// NewSpacesProvider creates a new DigitalOcean Spaces provider
func NewSpacesProvider() (*SpacesProvider, error) {
	cfg := config.GlobalConfig

	if cfg.SpacesAccessKeyID == "" || cfg.SpacesSecretAccessKey == "" {
		return nil, fmt.Errorf("Spaces credentials not provided")
	}

	if cfg.SpacesEndpoint == "" {
		return nil, fmt.Errorf("Spaces endpoint not configured")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(cfg.SpacesRegion),
		Endpoint: aws.String(cfg.SpacesEndpoint),
		Credentials: credentials.NewStaticCredentials(
			cfg.SpacesAccessKeyID,
			cfg.SpacesSecretAccessKey,
			"",
		),
		S3ForcePathStyle: aws.Bool(false),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Spaces session: %v", err)
	}

	client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	bucket := cfg.SpacesBucket
	if bucket == "" {
		return nil, fmt.Errorf("Spaces bucket not configured")
	}

	logger.Logger.Printf("✓ DigitalOcean Spaces provider initialized with bucket: %s", bucket)

	return &SpacesProvider{
		client:   client,
		uploader: uploader,
		bucket:   bucket,
		endpoint: cfg.SpacesEndpoint,
		region:   cfg.SpacesRegion,
	}, nil
}

// Upload uploads a file to DigitalOcean Spaces
func (p *SpacesProvider) Upload(file io.Reader, filename, folder string, fileSize int64) (*contracts.UploadResponse, error) {
	key := p.buildKey(filename, folder)

	// Read the file content
	buf := new(bytes.Buffer)
	size, err := buf.ReadFrom(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	_, err = p.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
		Body:   buf,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to Spaces: %v", err)
	}

	url := fmt.Sprintf("https://%s.%s.digitaloceanspaces.com/%s", p.bucket, p.region, key)

	return &contracts.UploadResponse{
		URL:      url,
		Filename: filename,
		Size:     size,
		Provider: "digitalocean-spaces",
	}, nil
}

// Download downloads a file from DigitalOcean Spaces
func (p *SpacesProvider) Download(filename, folder string) (io.Reader, error) {
	key := p.buildKey(filename, folder)

	result, err := p.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from Spaces: %v", err)
	}

	return result.Body, nil
}

// Delete deletes a file from DigitalOcean Spaces
func (p *SpacesProvider) Delete(filename, folder string) error {
	key := p.buildKey(filename, folder)

	_, err := p.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from Spaces: %v", err)
	}

	return nil
}

// GetFileURL returns the public URL of a file
func (p *SpacesProvider) GetFileURL(filename, folder string) (string, error) {
	key := p.buildKey(filename, folder)
	url := fmt.Sprintf("https://%s.%s.digitaloceanspaces.com/%s", p.bucket, p.region, key)
	return url, nil
}

// GetProviderName returns the provider name
func (p *SpacesProvider) GetProviderName() string {
	return "digitalocean-spaces"
}

// buildKey builds the Spaces key from filename and folder
func (p *SpacesProvider) buildKey(filename, folder string) string {
	if folder == "" {
		return filename
	}
	return filepath.Join(strings.Trim(folder, "/"), filename)
}
