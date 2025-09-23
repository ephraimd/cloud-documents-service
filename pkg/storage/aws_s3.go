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

type AWSS3Provider struct {
	client   *s3.S3
	uploader *s3manager.Uploader
	bucket   string
}

func NewAWSS3Provider() (*AWSS3Provider, error) {
	cfg := config.GlobalConfig

	if cfg.AWSAccessKeyID == "" || cfg.AWSSecretAccessKey == "" {
		return nil, fmt.Errorf("AWS credentials not provided")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWSRegion),
		Credentials: credentials.NewStaticCredentials(
			cfg.AWSAccessKeyID,
			cfg.AWSSecretAccessKey,
			"",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	bucket := cfg.AWSBucket
	if bucket == "" {
		return nil, fmt.Errorf("AWS bucket not configured")
	}

	logger.Logger.Printf("✓ AWS S3 provider initialized with bucket: %s", bucket)

	return &AWSS3Provider{
		client:   client,
		uploader: uploader,
		bucket:   bucket,
	}, nil
}

// Upload uploads a file to AWS S3
func (p *AWSS3Provider) Upload(file io.Reader, filename, folder string, fileSize int64) (*contracts.UploadResponse, error) {
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
		return nil, fmt.Errorf("failed to upload file to S3: %v", err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", p.bucket, config.GlobalConfig.AWSRegion, key)

	return &contracts.UploadResponse{
		URL:      url,
		Filename: filename,
		Size:     size,
		Provider: "aws-s3",
	}, nil
}

// Download downloads a file from AWS S3
func (p *AWSS3Provider) Download(filename, folder string) (io.Reader, error) {
	key := p.buildKey(filename, folder)

	result, err := p.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %v", err)
	}

	return result.Body, nil
}

// Delete deletes a file from AWS S3
func (p *AWSS3Provider) Delete(filename, folder string) error {
	key := p.buildKey(filename, folder)

	_, err := p.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %v", err)
	}

	return nil
}

// GetFileURL returns the public URL of a file
func (p *AWSS3Provider) GetFileURL(filename, folder string) (string, error) {
	key := p.buildKey(filename, folder)
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", p.bucket, config.GlobalConfig.AWSRegion, key)
	return url, nil
}

// GetProviderName returns the provider name
func (p *AWSS3Provider) GetProviderName() string {
	return "aws-s3"
}

// buildKey builds the S3 key from filename and folder
func (p *AWSS3Provider) buildKey(filename, folder string) string {
	if folder == "" {
		return filename
	}
	return filepath.Join(strings.Trim(folder, "/"), filename)
}
