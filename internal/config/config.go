package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/ephraimd/cloud-documents-service/internal/helpers"
	"github.com/ephraimd/cloud-documents-service/pkg/logger"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                     string
	Env                      string
	MaxFileSize              int64
	AllowedFileTypes         []string
	AllowedMimeTypes         []string
	MaxFilenameLength        int
	EnableFileTypeValidation bool
	EnableFileSizeValidation bool

	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSBucket          string

	SpacesAccessKeyID     string
	SpacesSecretAccessKey string
	SpacesRegion          string
	SpacesBucket          string
	SpacesEndpoint        string

	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string

	RedisAddr               string
	RedisPassword           string
	RedisDB                 int
	DefaultCacheExpiration  time.Duration
	MaxCachedRequestRetries int
	RetryStreamName         string
	EnableCaching           bool
	EnableRetries           bool
}

var GlobalConfig *Config

func LoadConfig() {
	err := godotenv.Load(".env")

	if err != nil {
		logger.Logger.Info("No .env file found, using environment variables from system")
	}

	GlobalConfig = &Config{
		Port: helpers.GetOsEnvOrDefault("PORT", "8081"),
		Env:  helpers.GetOsEnvOrDefault("ENV", "local"),
		MaxFileSize: func() int64 {
			val, _ := strconv.ParseInt(helpers.GetOsEnvOrDefault("MAX_FILE_SIZE", "10485760"), 10, 64)
			return val
		}(),
		AllowedFileTypes: func() []string {
			types := helpers.GetOsEnvOrDefault("ALLOWED_FILE_TYPES", "jpg,jpeg,png,gif,pdf,doc,docx,txt,csv,zip")
			if types == "" {
				return []string{}
			}
			return strings.Split(strings.ToLower(types), ",")
		}(),
		AllowedMimeTypes: func() []string {
			types := helpers.GetOsEnvOrDefault("ALLOWED_MIME_TYPES", "image/jpeg,image/png,image/gif,application/pdf,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,text/plain,text/csv,application/zip")
			if types == "" {
				return []string{}
			}
			return strings.Split(strings.ToLower(types), ",")
		}(),
		MaxFilenameLength: func() int {
			val, _ := strconv.Atoi(helpers.GetOsEnvOrDefault("MAX_FILENAME_LENGTH", "255"))
			return val
		}(),
		EnableFileTypeValidation: helpers.GetOsEnvOrDefault("ENABLE_FILE_TYPE_VALIDATION", "true") == "true",
		EnableFileSizeValidation: helpers.GetOsEnvOrDefault("ENABLE_FILE_SIZE_VALIDATION", "true") == "true",

		AWSAccessKeyID:     helpers.GetOsEnvOrDefault("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: helpers.GetOsEnvOrDefault("AWS_SECRET_ACCESS_KEY", ""),
		AWSRegion:          helpers.GetOsEnvOrDefault("AWS_REGION", "us-east-1"),
		AWSBucket:          helpers.GetOsEnvOrDefault("AWS_BUCKET", ""),

		SpacesAccessKeyID:     helpers.GetOsEnvOrDefault("SPACES_ACCESS_KEY_ID", ""),
		SpacesSecretAccessKey: helpers.GetOsEnvOrDefault("SPACES_SECRET_ACCESS_KEY", ""),
		SpacesRegion:          helpers.GetOsEnvOrDefault("SPACES_REGION", "nyc3"),
		SpacesBucket:          helpers.GetOsEnvOrDefault("SPACES_BUCKET", "geteco"),
		SpacesEndpoint:        helpers.GetOsEnvOrDefault("SPACES_ENDPOINT", ""),

		CloudinaryCloudName: helpers.GetOsEnvOrDefault("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    helpers.GetOsEnvOrDefault("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: helpers.GetOsEnvOrDefault("CLOUDINARY_API_SECRET", ""),

		RedisAddr:       helpers.GetOsEnvOrDefault("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   helpers.GetOsEnvOrDefault("REDIS_PASSWORD", ""),
		RetryStreamName: helpers.GetOsEnvOrDefault("RETRY_STREAM_NAME", "retry"),
		MaxCachedRequestRetries: func() int {
			val, _ := strconv.Atoi(helpers.GetOsEnvOrDefault("MAX_CACHED_REQUEST_RETRIES", "3"))
			return val
		}(),
		RedisDB: func() int {
			val, _ := strconv.Atoi(helpers.GetOsEnvOrDefault("REDIS_DB", "0"))
			return val
		}(),
		DefaultCacheExpiration: func() time.Duration {
			val, _ := strconv.Atoi(helpers.GetOsEnvOrDefault("DEFAULT_CACHE_EXPIRATION", "5"))
			return time.Duration(val) * time.Minute
		}(),
		EnableCaching: helpers.GetOsEnvOrDefault("ENABLE_CACHING", "false") == "true",
		EnableRetries: helpers.GetOsEnvOrDefault("ENABLE_RETRIES", "false") == "true",
	}
}

func init() {
	LoadConfig()
}
