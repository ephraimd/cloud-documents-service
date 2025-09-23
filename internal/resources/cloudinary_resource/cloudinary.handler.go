package cloudinaryresource

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/ephraimd/cloud-documents-service/internal/helpers"
	"github.com/ephraimd/cloud-documents-service/pkg/storage"
	"github.com/ephraimd/cloud-documents-service/pkg/validation"

	"github.com/gin-gonic/gin"
)

type CloudinaryHandlerImpl struct {
	provider  storage.CloudStorageProvider
	validator *validation.FileValidator
}

func NewCloudinaryHandler() (*CloudinaryHandlerImpl, error) {
	provider, err := storage.NewCloudinaryProvider()
	if err != nil {
		return nil, err
	}

	return &CloudinaryHandlerImpl{
		provider:  provider,
		validator: validation.NewFileValidator(),
	}, nil
}

// UploadFile godoc
// @Summary Upload file to Cloudinary
// @Description Upload a file to Cloudinary storage
// @Tags cloudinary
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param folder formData string false "Folder name (optional)"
// @Success 201 {object} contracts.UploadResponse
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /cloudinary/upload [post]
func (h *CloudinaryHandlerImpl) UploadFile(c *gin.Context) {
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
// @Summary Get file from Cloudinary
// @Description Retrieve a file from Cloudinary storage
// @Tags cloudinary
// @Produce json
// @Param filename path string true "Filename"
// @Param folder query string false "Folder name (optional)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Router /cloudinary/files/{filename} [get]
func (h *CloudinaryHandlerImpl) GetFile(c *gin.Context) {
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
// @Summary Delete file from Cloudinary
// @Description Delete a file from Cloudinary storage
// @Tags cloudinary
// @Param filename path string true "Filename"
// @Param folder query string false "Folder name (optional)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /cloudinary/files/{filename} [delete]
func (h *CloudinaryHandlerImpl) DeleteFile(c *gin.Context) {
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

func (h *CloudinaryHandlerImpl) validateFile(header *multipart.FileHeader) error {
	return h.validator.ValidateFile(header)
}

func (h *CloudinaryHandlerImpl) generateUniqueFilename(originalName string) string {
	// Generate unique filename to avoid conflicts
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}
