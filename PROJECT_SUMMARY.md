# Cloud Documents Service 


A production-ready cloud-based file upload service has been implemented following the exact same code and project structure pattern as the RevHero-user-fe-backend project.


## 🚀 Features Implemented

### ✅ Multi-Cloud Support
- **AWS S3**: Full implementation with upload, download, delete, URL generation
- **DigitalOcean Spaces**: S3-compatible implementation
- **Cloudinary**: Image/media optimization platform

### ✅ RESTful API Endpoints

#### Generic Upload (Multi-provider)
- `POST /v1/upload` - Upload to any provider
- `GET /v1/upload/providers` - List available providers

#### Provider-Specific Endpoints
- `POST /v1/aws/upload` - AWS S3 upload
- `GET /v1/aws/files/{filename}` - Get AWS S3 file URL
- `DELETE /v1/aws/files/{filename}` - Delete from AWS S3

- `POST /v1/spaces/upload` - DigitalOcean Spaces upload  
- `GET /v1/spaces/files/{filename}` - Get Spaces file URL
- `DELETE /v1/spaces/files/{filename}` - Delete from Spaces

- `POST /v1/cloudinary/upload` - Cloudinary upload
- `GET /v1/cloudinary/files/{filename}` - Get Cloudinary file URL
- `DELETE /v1/cloudinary/files/{filename}` - Delete from Cloudinary

#### Service Endpoints
- `GET /health` - Health check
- `GET /v1/swagger/*any` - API documentation


## 📋 Usage Examples

### Upload File to Any Provider
```bash
curl -X POST http://localhost:8080/v1/upload \
  -F "file=@example.jpg" \
  -F "folder=images" \
  -F "provider=aws"
```

### Upload to Specific Provider
```bash
curl -X POST http://localhost:8080/v1/aws/upload \
  -F "file=@example.jpg" \
  -F "folder=images"
```

### Get Available Providers
```bash
curl http://localhost:8080/v1/upload/providers
```

### Response Format (Example)
```json
{
  "data": {
    "upload": {
      "url": "https://bucket.s3.region.amazonaws.com/folder/filename.jpg",
      "filename": "filename.jpg", 
      "size": 12345,
      "provider": "aws-s3"
    }
  },
  "message": "File uploaded successfully"
}
```

## 🔧 Configuration

Set up environment variables in `.env`:

```bash
# Server
PORT=8080
ENV=local

# AWS S3
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret
AWS_REGION=us-east-1
AWS_BUCKET=your_bucket

# DigitalOcean Spaces  
SPACES_ACCESS_KEY_ID=your_key
SPACES_SECRET_ACCESS_KEY=your_secret
SPACES_REGION=nyc3
SPACES_BUCKET=your_bucket
SPACES_ENDPOINT=https://nyc3.digitaloceanspaces.com

# Cloudinary
CLOUDINARY_CLOUD_NAME=your_cloud
CLOUDINARY_API_KEY=your_key
CLOUDINARY_API_SECRET=your_secret
```
