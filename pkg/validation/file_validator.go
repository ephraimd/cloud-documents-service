package validation

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/ephraimd/cloud-documents-service/internal/config"
	"github.com/h2non/filetype"
)

type FileValidator struct {
	config *config.Config
}

func NewFileValidator() *FileValidator {
	return &FileValidator{
		config: config.GlobalConfig,
	}
}

func (v *FileValidator) ValidateFile(header *multipart.FileHeader) error {
	if err := v.validateFileSize(header.Size); err != nil {
		return err
	}

	if err := v.validateFileType(header.Filename); err != nil {
		return err
	}

	if err := v.validateFilename(header.Filename); err != nil {
		return err
	}

	return nil
}

func (v *FileValidator) validateFileSize(size int64) error {
	if !v.config.EnableFileSizeValidation {
		return nil
	}

	if size <= 0 {
		return fmt.Errorf("file size cannot be zero or negative")
	}

	if size > v.config.MaxFileSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes (%.2f MB)",
			size, v.config.MaxFileSize, float64(v.config.MaxFileSize)/1024/1024)
	}

	return nil
}

func (v *FileValidator) validateFileType(filename string) error {
	if !v.config.EnableFileTypeValidation {
		return nil
	}

	if len(v.config.AllowedFileTypes) == 0 {
		return nil
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	if ext == "" {
		return fmt.Errorf("file must have a valid extension")
	}

	for _, allowedType := range v.config.AllowedFileTypes {
		if ext == strings.TrimSpace(allowedType) {
			return nil
		}
	}

	return fmt.Errorf("file type '%s' is not allowed. Allowed types: %s",
		ext, strings.Join(v.config.AllowedFileTypes, ", "))
}

func (v *FileValidator) validateFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	if len(filename) > v.config.MaxFilenameLength {
		return fmt.Errorf("filename exceeds maximum length of %d characters", v.config.MaxFilenameLength)
	}

	dangerousChars := []string{"../", "..\\", "<", ">", ":", "\"", "|", "?", "*", "\x00"}
	for _, char := range dangerousChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("filename contains invalid character: %s", char)
		}
	}

	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4",
		"COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5",
		"LPT6", "LPT7", "LPT8", "LPT9"}

	baseFilename := strings.ToUpper(strings.TrimSuffix(filename, filepath.Ext(filename)))
	for _, reserved := range reservedNames {
		if baseFilename == reserved {
			return fmt.Errorf("filename '%s' is a reserved system name", filename)
		}
	}

	return nil
}

func (v *FileValidator) ValidateMimeType(file multipart.File) error {
	if !v.config.EnableFileTypeValidation || len(v.config.AllowedMimeTypes) == 0 {
		return nil
	}

	// Read header bytes
	header := make([]byte, 261) // filetype recommends at least 261 bytes
	_, err := file.Read(header)
	if err != nil {
		return fmt.Errorf("failed to read file for MIME type detection: %v", err)
	}

	// reset cursor if possible
	if seeker, ok := file.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	// Detect using filetype
	kind, unknown := filetype.Match(header)
	if unknown != nil {
		// fallback to generic mime sniffing
		return fmt.Errorf("unknown MIME type: %v", unknown)
	}

	detectedType := strings.ToLower(strings.TrimSpace(kind.MIME.Value))

	// Normalize allowed list
	for _, allowed := range v.config.AllowedMimeTypes {
		allowed = strings.ToLower(strings.TrimSpace(allowed))
		if detectedType == allowed {
			return nil
		}
	}

	return fmt.Errorf("MIME type '%s' is not allowed. Allowed types: %s",
		detectedType, strings.Join(v.config.AllowedMimeTypes, ", "))
}

func DetectMimeType(data []byte) string {
	detectedType := http.DetectContentType(data)

	if detectedType == "text/plain; charset=utf-8" || detectedType == "text/plain" {
		return "text/plain"
	}

	if idx := strings.Index(detectedType, ";"); idx != -1 {
		return detectedType[:idx]
	}

	return detectedType
}

func (v *FileValidator) GetValidationSummary() map[string]interface{} {
	return map[string]interface{}{
		"max_file_size_bytes":  v.config.MaxFileSize,
		"max_file_size_mb":     float64(v.config.MaxFileSize) / 1024 / 1024,
		"allowed_file_types":   v.config.AllowedFileTypes,
		"allowed_mime_types":   v.config.AllowedMimeTypes,
		"max_filename_length":  v.config.MaxFilenameLength,
		"file_type_validation": v.config.EnableFileTypeValidation,
		"file_size_validation": v.config.EnableFileSizeValidation,
	}
}
