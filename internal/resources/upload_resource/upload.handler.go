package uploadresource

import (
	"mime/multipart"
	"net/http"

	"github.com/ephraimd/cloud-documents-service/internal/helpers"
	"github.com/ephraimd/cloud-documents-service/pkg/storage"
	"github.com/ephraimd/cloud-documents-service/pkg/validation"

	"github.com/gin-gonic/gin"
)

type UploadHandlerImpl struct {
	factory   *storage.StorageFactory
	validator *validation.FileValidator
}

func NewUploadHandler() (*UploadHandlerImpl, error) {
	factory := storage.NewStorageFactory()

	return &UploadHandlerImpl{
		factory:   factory,
		validator: validation.NewFileValidator(),
	}, nil
}

// UploadFile godoc
// @Summary Upload file to specified provider
// @Description Upload a file to a cloud storage provider (aws, spaces, cloudinary)
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param folder formData string false "Folder/bucket name (optional)"
// @Param provider formData string true "Storage provider (aws, spaces, cloudinary)"
// @Success 201 {object} contracts.UploadResponse
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /upload [post]
func (h *UploadHandlerImpl) UploadFile(c *gin.Context) {
	providerName := c.PostForm("provider")
	if providerName == "" {
		code, resp := helpers.RespondError("Provider is required", &gin.H{"allowed_providers": []string{"aws", "spaces", "cloudinary"}}, http.StatusBadRequest)
		c.JSON(code, resp)
		return
	}
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

	if err := h.validateFile(header); err != nil {
		code, resp := helpers.RespondError("File validation failed", &gin.H{"error": err.Error()}, http.StatusBadRequest)
		c.JSON(code, resp)
		return
	}

	if err := h.validator.ValidateMimeType(file); err != nil {
		code, resp := helpers.RespondError("File MIME type validation failed", &gin.H{"error": err.Error()}, http.StatusBadRequest)
		c.JSON(code, resp)
		return
	}

	supportedProviders := []string{"aws", "spaces", "cloudinary"}
	isValidProvider := false
	for _, p := range supportedProviders {
		if providerName == p || providerName == p+"-s3" || providerName == "digitalocean-"+p || providerName == "s3" {
			isValidProvider = true
			break
		}
	}

	if !isValidProvider {
		code, resp := helpers.RespondError("Invalid provider", &gin.H{"error": "Unsupported provider: " + providerName, "allowed_providers": supportedProviders}, http.StatusBadRequest)
		c.JSON(code, resp)
		return
	}

	provider, err := h.factory.GetProvider(providerName)
	if err != nil {
		code, resp := helpers.RespondError("Provider configuration error", &gin.H{"error": err.Error(), "allowed_providers": supportedProviders}, http.StatusBadRequest)
		c.JSON(code, resp)
		return
	}

	result, err := provider.Upload(file, header.Filename, folder, header.Size)
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

// GetProviders godoc
// @Summary Get available storage providers
// @Description Get list of available cloud storage providers
// @Tags upload
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /upload/providers [get]
func (h *UploadHandlerImpl) GetProviders(c *gin.Context) {
	providers := h.factory.GetAllProviders()

	providerList := make([]map[string]interface{}, 0)
	for name, provider := range providers {
		providerList = append(providerList, map[string]interface{}{
			"name":         name,
			"display_name": provider.GetProviderName(),
			"available":    true,
		})
	}

	code, resp := helpers.RespondOk("Available providers retrieved successfully", &gin.H{
		"providers": providerList,
		"count":     len(providerList),
	})
	c.JSON(code, resp)
}

// GetValidationSettings godoc
// @Summary Get file validation settings
// @Description Get current file validation configuration including size limits and allowed types
// @Tags upload
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /upload/validation [get]
func (h *UploadHandlerImpl) GetValidationSettings(c *gin.Context) {
	settings := h.validator.GetValidationSummary()

	code, resp := helpers.RespondOk("Validation settings retrieved successfully", &gin.H{
		"validation": settings,
	})
	c.JSON(code, resp)
}

func (h *UploadHandlerImpl) validateFile(header *multipart.FileHeader) error {
	return h.validator.ValidateFile(header)
}
