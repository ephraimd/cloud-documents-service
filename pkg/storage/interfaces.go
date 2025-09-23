package storage

import (
	"github.com/ephraimd/cloud-documents-service/internal/contracts"
	"io"
)

type CloudStorageProvider interface {
	Upload(file io.Reader, filename, folder string, fileSize int64) (*contracts.UploadResponse, error)
	Download(filename, folder string) (io.Reader, error)
	Delete(filename, folder string) error
	GetFileURL(filename, folder string) (string, error)
	GetProviderName() string
}
