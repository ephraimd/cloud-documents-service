package storage

import (
	"fmt"
)

type StorageFactory struct{}

func NewStorageFactory() *StorageFactory {
	return &StorageFactory{}
}

func (f *StorageFactory) GetProvider(providerName string) (CloudStorageProvider, error) {
	switch providerName {
	case "aws", "aws-s3", "s3":
		return NewAWSS3Provider()
	case "spaces", "digitalocean-spaces", "digitalocean":
		return NewSpacesProvider()
	case "cloudinary":
		return NewCloudinaryProvider()
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", providerName)
	}
}

func (f *StorageFactory) GetAllProviders() map[string]CloudStorageProvider {
	providers := make(map[string]CloudStorageProvider)

	if aws, err := NewAWSS3Provider(); err == nil {
		providers["aws"] = aws
	}

	if spaces, err := NewSpacesProvider(); err == nil {
		providers["spaces"] = spaces
	}

	if cloudinary, err := NewCloudinaryProvider(); err == nil {
		providers["cloudinary"] = cloudinary
	}

	return providers
}
