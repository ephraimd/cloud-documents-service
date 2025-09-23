package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ephraimd/cloud-documents-service/internal/config"
	"github.com/ephraimd/cloud-documents-service/internal/routes"

	"github.com/gin-gonic/gin"

	_ "github.com/ephraimd/cloud-documents-service/docs"
	_ "github.com/ephraimd/cloud-documents-service/internal/resources/aws_resource"
	_ "github.com/ephraimd/cloud-documents-service/internal/resources/cloudinary_resource"
	_ "github.com/ephraimd/cloud-documents-service/internal/resources/spaces_resource"
	_ "github.com/ephraimd/cloud-documents-service/internal/resources/upload_resource"
)

// @title Cloud Documents Service API
// @version 1.0
// @description Production-ready cloud based file upload service supporting AWS S3, DigitalOcean Spaces, and Cloudinary
// @host localhost:8080
// @BasePath /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if config.GlobalConfig.Env != "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := []string{"http://localhost:3000", "http://127.0.0.1:3000"}

		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else if config.GlobalConfig.Env == "local" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Next()
	})

	routes.SetupRouter(router)

	server := &http.Server{
		Addr:           ":" + config.GlobalConfig.Port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Starting Cloud Documents Service on port %s", config.GlobalConfig.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %s", err)
	}
	log.Println("Server stopped gracefully")
}
