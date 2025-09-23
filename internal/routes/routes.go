package routes

import (
	"github.com/ephraimd/cloud-documents-service/internal/helpers"
	awsresource "github.com/ephraimd/cloud-documents-service/internal/resources/aws_resource"
	cloudinaryresource "github.com/ephraimd/cloud-documents-service/internal/resources/cloudinary_resource"
	spacesresource "github.com/ephraimd/cloud-documents-service/internal/resources/spaces_resource"
	uploadresource "github.com/ephraimd/cloud-documents-service/internal/resources/upload_resource"
	"github.com/ephraimd/cloud-documents-service/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var defaultRoute = func(ctx *gin.Context) {
	ctx.JSON(helpers.RespondOk("Cloud Documents Service Active", &gin.H{"status": "active", "service": "cloud-documents-service"}))
}

func SetupRouter(router *gin.Engine) {
	setDefaultRoutes(router)

	v1Group := router.Group("/v1")
	setDefaultGroupRoutes(v1Group)

	setupSwagger(v1Group)
	setupResourceRoutes(v1Group)

	logger.Logger.Printf("✓ All routes configured successfully")
}

func setDefaultRoutes(router *gin.Engine) {
	router.GET("/", defaultRoute)
	router.GET("/health", defaultRoute)
}

func setDefaultGroupRoutes(v1Group *gin.RouterGroup) {
	v1Group.GET("/", defaultRoute)
	v1Group.GET("/health", defaultRoute)
}

func setupSwagger(v1Group *gin.RouterGroup) {
	v1Group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.Logger.Printf("   ✓ Swagger documentation available at /v1/swagger/")
}

func setupResourceRoutes(v1Group *gin.RouterGroup) {
	uploadHandler, err := uploadresource.NewUploadHandler()
	if err != nil {
		logger.Logger.Printf("   ⚠ Failed to initialize upload handler: %v", err)
	} else {
		uploadGroup := v1Group.Group("/upload")
		{
			uploadGroup.POST("/", uploadHandler.UploadFile)
			uploadGroup.GET("/providers", uploadHandler.GetProviders)
			uploadGroup.GET("/validation", uploadHandler.GetValidationSettings)
		}
		logger.Logger.Printf("   ✓ Generic upload routes configured")
	}

	awsHandler, err := awsresource.NewAWSHandler()
	if err != nil {
		logger.Logger.Printf("   ⚠ Failed to initialize AWS handler: %v", err)
	} else {
		awsGroup := v1Group.Group("/aws")
		{
			awsGroup.POST("/upload", awsHandler.UploadFile)
			awsGroup.GET("/files/:filename", awsHandler.GetFile)
			awsGroup.DELETE("/files/:filename", awsHandler.DeleteFile)
		}
		logger.Logger.Printf("   ✓ AWS S3 routes configured")
	}

	spacesHandler, err := spacesresource.NewSpacesHandler()
	if err != nil {
		logger.Logger.Printf("   ⚠ Failed to initialize Spaces handler: %v", err)
	} else {
		spacesGroup := v1Group.Group("/spaces")
		{
			spacesGroup.POST("/upload", spacesHandler.UploadFile)
			spacesGroup.GET("/files/:filename", spacesHandler.GetFile)
			spacesGroup.DELETE("/files/:filename", spacesHandler.DeleteFile)
		}
		logger.Logger.Printf("   ✓ DigitalOcean Spaces routes configured")
	}

	cloudinaryHandler, err := cloudinaryresource.NewCloudinaryHandler()
	if err != nil {
		logger.Logger.Printf("   ⚠ Failed to initialize Cloudinary handler: %v", err)
	} else {
		cloudinaryGroup := v1Group.Group("/cloudinary")
		{
			cloudinaryGroup.POST("/upload", cloudinaryHandler.UploadFile)
			cloudinaryGroup.GET("/files/:filename", cloudinaryHandler.GetFile)
			cloudinaryGroup.DELETE("/files/:filename", cloudinaryHandler.DeleteFile)
		}
		logger.Logger.Printf("   ✓ Cloudinary routes configured")
	}
}
