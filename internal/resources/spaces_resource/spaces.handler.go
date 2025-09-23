package spacesresource

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"ephraimd/cloud-document-service/internal/helpers"
	"ephraimd/cloud-document-service/pkg/storage"
	"ephraimd/cloud-document-service/pkg/validation"

	"github.com/gin-gonic/gin"
)

type SpacesHandlerImpl struct {
	provider  storage.CloudStorageProvider
	validator *validation.FileValidator
}

func NewSpacesHandler() (*SpacesHandlerImpl, error) {
	provider, err := storage.NewSpacesProvider()
	if err != nil {
		return nil, err
	}

	return &SpacesHandlerImpl{
		provider:  provider,
		validator: validation.NewFileValidator(),
	}, nil
}

// UploadFile godoc
// @Summary Upload file to DigitalOcean Spaces
// @Description Upload a file to DigitalOcean Spaces storage
// @Tags spaces
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param folder formData string false "Folder/bucket name (optional)"
// @Success 201 {object} contracts.UploadResponse
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /spaces/upload [post]
func (h *SpacesHandlerImpl) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		code, resp := helpers.RespondError("No file provided", &gin.H{"error": err.Error()}, http.StatusBadRequest)
		c.JSON(code, resp)
		return
	}
	defer file.Close()

	folder := c.PostForm("folder")
	if folder == "" {
		folder = "uploads"
	}

	// Validate file type and size
	if err := h.validateFile(header); err != nil {
		code, resp := helpers.RespondError("File validation failed", &gin.H{"error": err.Error()}, http.StatusBadRequest)
		c.JSON(code, resp)
		return
	}

	// Validate MIME type
	if err := h.validator.ValidateMimeType(file); err != nil {
		code, resp := helpers.RespondError("File MIME type validation failed", &gin.H{"error": err.Error()}, http.StatusBadRequest)
		c.JSON(code, resp)
		return
	}

	filename := h.generateUniqueFilename(header.Filename)

	result, err := h.provider.Upload(file, filename, folder, header.Size)
	if err != nil {
		code, resp := helpers.RespondError("Failed to upload file", &gin.H{"error": err.Error()}, http.StatusInternalServerError)
		c.JSON(code, resp)
		return
	}

	code, resp := helpers.RespondCreated("File uploaded successfully", &gin.H{
		"upload": result,
	})
	c.JSON(code, resp)
}

// GetFile godoc
// @Summary Get file from DigitalOcean Spaces
// @Description Retrieve a file from DigitalOcean Spaces storage
// @Tags spaces
// @Produce json
// @Param filename path string true "Filename"
// @Param folder query string false "Folder name (optional)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Router /spaces/files/{filename} [get]
func (h *SpacesHandlerImpl) GetFile(c *gin.Context) {
	filename := c.Param("filename")
	folder := c.Query("folder")

	if filename == "" {
		code, resp := helpers.RespondError("Filename is required", nil, http.StatusBadRequest)
		c.JSON(code, resp)
		return
	}

	url, err := h.provider.GetFileURL(filename, folder)
	if err != nil {
		code, resp := helpers.RespondError("Failed to get file URL", &gin.H{"error": err.Error()}, http.StatusNotFound)
		c.JSON(code, resp)
		return
	}

	code, resp := helpers.RespondOk("File URL retrieved successfully", &gin.H{
		"url":      url,
		"filename": filename,
		"provider": h.provider.GetProviderName(),
	})
	c.JSON(code, resp)
}

// DeleteFile godoc
// @Summary Delete file from DigitalOcean Spaces
// @Description Delete a file from DigitalOcean Spaces storage
// @Tags spaces
// @Param filename path string true "Filename"
// @Param folder query string false "Folder name (optional)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /spaces/files/{filename} [delete]
func (h *SpacesHandlerImpl) DeleteFile(c *gin.Context) {
	filename := c.Param("filename")
	folder := c.Query("folder")

	if filename == "" {
		code, resp := helpers.RespondError("Filename is required", nil, http.StatusBadRequest)
		c.JSON(code, resp)
		return
	}

	err := h.provider.Delete(filename, folder)
	if err != nil {
		code, resp := helpers.RespondError("Failed to delete file", &gin.H{"error": err.Error()}, http.StatusInternalServerError)
		c.JSON(code, resp)
		return
	}

	code, resp := helpers.RespondOk("File deleted successfully", &gin.H{
		"filename": filename,
		"provider": h.provider.GetProviderName(),
	})
	c.JSON(code, resp)
}

func (h *SpacesHandlerImpl) validateFile(header *multipart.FileHeader) error {
	return h.validator.ValidateFile(header)
}

func (h *SpacesHandlerImpl) generateUniqueFilename(originalName string) string {
	// Generate unique filename to avoid conflicts
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}
