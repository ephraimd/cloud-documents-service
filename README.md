# Cloud Documents Service

A production-ready cloud-based file upload service supporting multiple cloud storage providers including AWS S3, DigitalOcean Spaces, and Cloudinary. Built following best practices.

## Features

- **Multi-provider Support**: Upload files to AWS S3, DigitalOcean Spaces, and Cloudinary
- **RESTful API**: Simple and easy-to-use REST endpoints
- **File Management**: Upload, retrieve, and delete files
- **Folder Organization**: Optional folder/bucket organization
- **File Validation**: Comprehensive file type, size, and MIME type validation
- **Swagger Documentation**: Complete API documentation with interactive testing
- **Production Ready**: Following Go best practices and patterns
- **Configurable**: Environment-based configuration
- **Security**: File sanitization, size limits, and type restrictions
- **Error Handling**: Clear error messages and status codes

## Supported Cloud Providers

### AWS S3
- Upload files to Amazon S3 buckets
- Public URL generation
- File management operations

### DigitalOcean Spaces
- Upload files to DigitalOcean Spaces
- S3-compatible API
- CDN-ready URLs

### Cloudinary
- Image and video optimization
- Automatic format conversion
- Advanced media management

## API Endpoints

### File Upload & Management

#### Generic Upload (Recommended)
- **POST** `/v1/upload/` - Upload to any provider
- **GET** `/v1/upload/providers` - Get available providers  
- **GET** `/v1/upload/validation` - Get file validation settings

#### Provider-Specific Endpoints

**AWS S3**
- **POST** `/v1/aws/upload` - Upload file to AWS S3
- **GET** `/v1/aws/files/{filename}` - Get file URL from AWS S3
- **DELETE** `/v1/aws/files/{filename}` - Delete file from AWS S3

**DigitalOcean Spaces**
- **POST** `/v1/spaces/upload` - Upload file to Spaces
- **GET** `/v1/spaces/files/{filename}` - Get file URL from Spaces
- **DELETE** `/v1/spaces/files/{filename}` - Delete file from Spaces

**Cloudinary**
- **POST** `/v1/cloudinary/upload` - Upload file to Cloudinary
- **GET** `/v1/cloudinary/files/{filename}` - Get file URL from Cloudinary
- **DELETE** `/v1/cloudinary/files/{filename}` - Delete file from Cloudinary

#### Documentation & Health
- **GET** `/v1/swagger/` - Interactive Swagger API documentation
- **GET** `/v1/health` - Service health check
- **GET** `/` - Service information

## API Usage Guide

### 1. File Upload Examples

#### Upload to Any Provider (Generic Endpoint)

**Request:**
```bash
curl -X POST http://localhost:8081/v1/upload/ \
  -F "file=@document.pdf" \
  -F "folder=documents" \
  -F "provider=aws"
```

**Request Parameters:**
- `file` (required): The file to upload
- `folder` (optional): Folder/directory name (default: "uploads")
- `provider` (required): Provider name (`aws`, `spaces`, or `cloudinary`)

**Success Response (201 Created):**
```json
{
  "data": {
    "upload": {
      "url": "https://onboard-test.s3.us-east-1.amazonaws.com/documents/document_1695456789.pdf",
      "filename": "document_1695456789.pdf",
      "size": 2048576,
      "provider": "aws-s3"
    }
  },
  "message": "File uploaded successfully"
}
```

**Error Response - File Validation Failed (400 Bad Request):**
```json
{
  "data": {
    "error": "file type 'exe' is not allowed. Allowed types: jpg, jpeg, png, gif, pdf, doc, docx, txt, csv, zip, mp4, mov, avi"
  },
  "message": "File validation failed"
}
```

**Error Response - File Too Large (400 Bad Request):**
```json
{
  "data": {
    "error": "file size 15728640 bytes exceeds maximum allowed size of 10485760 bytes (10.00 MB)"
  },
  "message": "File validation failed"
}
```

**Error Response - Invalid Provider (400 Bad Request):**
```json
{
  "data": {
    "error": "Unsupported provider: invalid_provider",
    "allowed_providers": ["aws", "spaces", "cloudinary"]
  },
  "message": "Invalid provider"
}
```

#### Upload to AWS S3 (Provider-Specific)

**Request:**
```bash
curl -X POST http://localhost:8081/v1/aws/upload \
  -F "file=@image.jpg" \
  -F "folder=images"
```

**Success Response (201 Created):**
```json
{
  "data": {
    "upload": {
      "url": "https://onboard-test.s3.us-east-1.amazonaws.com/images/image_1695456789.jpg",
      "filename": "image_1695456789.jpg",
      "size": 1024567,
      "provider": "aws-s3"
    }
  },
  "message": "File uploaded successfully"
}
```

#### Upload to DigitalOcean Spaces

**Request:**
```bash
curl -X POST http://localhost:8081/v1/spaces/upload \
  -F "file=@video.mp4" \
  -F "folder=videos"
```

**Success Response (201 Created):**
```json
{
  "data": {
    "upload": {
      "url": "https://your-bucket.nyc3.digitaloceanspaces.com/videos/video_1695456789.mp4",
      "filename": "video_1695456789.mp4",
      "size": 5242880,
      "provider": "digitalocean-spaces"
    }
  },
  "message": "File uploaded successfully"
}
```

#### Upload to Cloudinary

**Request:**
```bash
curl -X POST http://localhost:8081/v1/cloudinary/upload \
  -F "file=@photo.png" \
  -F "folder=photos"
```

**Success Response (201 Created):**
```json
{
  "data": {
    "upload": {
      "url": "https://res.cloudinary.com/your-cloud/image/upload/v1695456789/photos/photo_1695456789.png",
      "filename": "photo_1695456789.png",
      "size": 512000,
      "provider": "cloudinary"
    }
  },
  "message": "File uploaded successfully"
}
```

### 2. File Retrieval Examples

#### Get File URL

**Request:**
```bash
curl http://localhost:8081/v1/aws/files/document_1695456789.pdf?folder=documents
```

**Success Response (200 OK):**
```json
{
  "data": {
    "file": {
      "url": "https://onboard-test.s3.us-east-1.amazonaws.com/documents/document_1695456789.pdf",
      "filename": "document_1695456789.pdf",
      "provider": "aws-s3"
    }
  },
  "message": "File URL retrieved successfully"
}
```

**Error Response - File Not Found (404 Not Found):**
```json
{
  "data": {
    "error": "File not found"
  },
  "message": "File retrieval failed"
}
```

### 3. File Deletion Examples

#### Delete File

**Request:**
```bash
curl -X DELETE http://localhost:8081/v1/aws/files/document_1695456789.pdf?folder=documents
```

**Success Response (200 OK):**
```json
{
  "data": {
    "filename": "document_1695456789.pdf",
    "provider": "aws-s3"
  },
  "message": "File deleted successfully"
}
```

**Error Response - File Not Found (404 Not Found):**
```json
{
  "data": {
    "error": "File not found"
  },
  "message": "File deletion failed"
}
```

### 4. Service Information Examples

#### Get Available Providers

**Request:**
```bash
curl http://localhost:8081/v1/upload/providers
```

**Response (200 OK):**
```json
{
  "data": {
    "providers": [
      {
        "name": "aws",
        "display_name": "Amazon S3",
        "available": true
      }
    ],
    "count": 1
  },
  "message": "Available providers retrieved successfully"
}
```

#### Get File Validation Settings

**Request:**
```bash
curl http://localhost:8081/v1/upload/validation
```

**Response (200 OK):**
```json
{
  "data": {
    "validation": {
      "max_file_size_bytes": 10485760,
      "max_file_size_mb": 10,
      "allowed_file_types": [
        "jpg", "jpeg", "png", "gif", "pdf", "doc", "docx", 
        "txt", "csv", "zip", "mp4", "mov", "avi"
      ],
      "allowed_mime_types": [
        "image/jpeg", "image/png", "image/gif", "application/pdf",
        "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
        "text/plain", "text/csv", "application/zip", 
        "video/mp4", "video/quicktime", "video/x-msvideo"
      ],
      "max_filename_length": 255,
      "file_type_validation": true,
      "file_size_validation": true
    }
  },
  "message": "Validation settings retrieved successfully"
}
```

#### Health Check

**Request:**
```bash
curl http://localhost:8081/v1/health
```

**Response (200 OK):**
```json
{
  "data": {
    "status": "healthy",
    "timestamp": "2025-09-23T12:00:00Z",
    "version": "1.0.0"
  },
  "message": "Service is healthy"
}
```

### 5. Programming Language Examples

#### JavaScript/Node.js Example

```javascript
const FormData = require('form-data');
const fs = require('fs');
const axios = require('axios');

async function uploadFile() {
  const formData = new FormData();
  formData.append('file', fs.createReadStream('./document.pdf'));
  formData.append('folder', 'documents');
  formData.append('provider', 'aws');

  try {
    const response = await axios.post('http://localhost:8081/v1/upload/', formData, {
      headers: {
        ...formData.getHeaders()
      }
    });
    
    console.log('Upload successful:', response.data);
    return response.data.data.upload.url;
  } catch (error) {
    console.error('Upload failed:', error.response.data);
    throw error;
  }
}

uploadFile();
```

#### Python Example

```python
import requests

def upload_file():
    url = 'http://localhost:8081/v1/upload/'
    
    files = {'file': open('document.pdf', 'rb')}
    data = {
        'folder': 'documents',
        'provider': 'aws'
    }
    
    try:
        response = requests.post(url, files=files, data=data)
        response.raise_for_status()
        
        result = response.json()
        print('Upload successful:', result)
        return result['data']['upload']['url']
    except requests.exceptions.RequestException as e:
        print('Upload failed:', e)
        raise
    finally:
        files['file'].close()

upload_file()
```

#### Go Example

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
)

func uploadFile() error {
    url := "http://localhost:8081/v1/upload/"
    
    // Create multipart form
    var b bytes.Buffer
    w := multipart.NewWriter(&b)
    
    // Add file
    f, err := os.Open("document.pdf")
    if err != nil {
        return err
    }
    defer f.Close()
    
    fw, err := w.CreateFormFile("file", "document.pdf")
    if err != nil {
        return err
    }
    
    if _, err = io.Copy(fw, f); err != nil {
        return err
    }
    
    // Add other fields
    w.WriteField("folder", "documents")
    w.WriteField("provider", "aws")
    w.Close()
    
    // Make request
    req, err := http.NewRequest("POST", url, &b)
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", w.FormDataContentType())
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    fmt.Printf("Upload status: %s\n", resp.Status)
    return nil
}
```

#### PHP Example

```php
<?php
function uploadFile() {
    $url = 'http://localhost:8081/v1/upload/';
    
    $file = new CURLFile('document.pdf', 'application/pdf', 'document.pdf');
    
    $data = [
        'file' => $file,
        'folder' => 'documents',
        'provider' => 'aws'
    ];
    
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    
    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    
    if ($httpCode === 201) {
        $result = json_decode($response, true);
        echo "Upload successful: " . $result['data']['upload']['url'] . "\n";
        return $result['data']['upload']['url'];
    } else {
        echo "Upload failed: " . $response . "\n";
        throw new Exception("Upload failed");
    }
}

uploadFile();
?>
```


## File Validation Rules

### Supported File Types
- **Images**: jpg, jpeg, png, gif
- **Documents**: pdf, doc, docx, txt, csv
- **Archives**: zip
- **Videos**: mp4, mov, avi

### File Size Limits
- **Default Maximum**: 10 MB (10,485,760 bytes)
- **Configurable**: Set via `MAX_FILE_SIZE` environment variable

### Security Features
- **MIME Type Validation**: Files are validated based on actual content, not just extension
- **Filename Sanitization**: Automatic sanitization of uploaded filenames
- **Path Traversal Protection**: Prevention of directory traversal attacks
- **Extension Validation**: Whitelist-based file extension checking
- `GET /v1/cloudinary/files/{filename}` - Get file URL from Cloudinary
- `DELETE /v1/cloudinary/files/{filename}` - Delete file from Cloudinary

### Documentation
- `GET /v1/swagger/` - Swagger API documentation

## Quick Start

### 1. Installation

```bash
# Clone the repository
git clone <repository-url>
cd cloud-documents-service

# Install dependencies
go mod tidy

# Copy environment configuration
cp .env.example .env
```



## Troubleshooting

### Common Issues and Solutions

#### 1. CORS Issues in Browser

**Problem**: Browser blocks requests due to CORS policy.

**Solution**: The service includes CORS middleware. If issues persist:
```bash
# Check CORS headers in response:
curl -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: Content-Type" \
     -X OPTIONS \
     http://localhost:8081/v1/upload/
```


### Health Check Endpoint

Monitor service health:
```bash
# Basic health check
curl http://localhost:8081/v1/health

# Expected response:
{
  "data": {
    "status": "healthy",
    "timestamp": "2025-09-23T12:00:00Z"
  },
  "message": "Service is healthy"
}
```


## Production Deployment

### Environment Setup

```bash
# Production environment variables
ENV=production
GIN_MODE=release
PORT=8081

# Security headers (automatically applied)
# - CORS enabled
# - Content Security Policy
# - X-Frame-Options: DENY
# - X-Content-Type-Options: nosniff
# - X-XSS-Protection: 1; mode=block
```


## Development & Contributing

### Local Development Setup

```bash
# 1. Clone and setup
git clone <repository-url>
cd cloud-documents-service
cp .env.example .env

# 2. Install dependencies
go mod tidy

# 3. Configure at least one provider in .env
# 4. Run the service
go run cmd/main.go

# 5. Test basic functionality
curl http://localhost:8081/v1/health
```


### Testing Guidelines

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/storage/

# Integration tests (require environment setup)
go test -tags=integration ./...
```


## License

This project is licensed under [BSD 2-Clause License] - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions:

1. **Documentation**: Check this README and `/v1/swagger/` endpoint
2. **Issues**: Open an issue in the repository
3. **Debugging**: Enable debug mode and check logs
4. **Health Check**: Use `/v1/health` endpoint to verify service status

---

**Cloud Documents Service** - A production-ready, multi-cloud file upload service built with Go
