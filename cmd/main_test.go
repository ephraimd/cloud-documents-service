package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ephraimd/cloud-documents-service/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	routes.SetupRouter(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "message")
}

func TestGetProviders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	routes.SetupRouter(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/upload/providers", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "data")
}

func createTestFile() (*bytes.Buffer, *multipart.Writer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a test file
	part, _ := writer.CreateFormFile("file", "test.txt")
	io.WriteString(part, "This is a test file content")

	// Add provider field
	writer.WriteField("provider", "aws")
	writer.WriteField("folder", "test")

	writer.Close()
	return body, writer
}

func TestUploadWithoutCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	routes.SetupRouter(router)

	body, writer := createTestFile()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Should fail without proper credentials configured
	assert.True(t, w.Code >= 400)
}
