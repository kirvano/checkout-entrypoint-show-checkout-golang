package aws

import (
	"path/filepath"
	"strings"

	appConfig "checkout-go/internal/config"
)

// S3FileDriver implements the FileDriver interface for S3
type S3FileDriver struct {
	basePath string
	bucket   string
}

// NewS3FileDriver creates a new S3-based file driver using the provided configuration
func NewS3FileDriver(cfg *appConfig.Config) *S3FileDriver {
	bucket := cfg.AWSS3Bucket
	if bucket == "" {
		bucket = "default-bucket"
	}

	basePath := cfg.GetS3BasePath()
	
	// Ensure base path ends with slash
	if basePath != "" && !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	return &S3FileDriver{
		basePath: basePath,
		bucket:   bucket,
	}
}

// GetBasePath returns the base path for file URLs
func (d *S3FileDriver) GetBasePath() string {
	return d.basePath
}

// GetFullPath returns the full URL for a given relative path
func (d *S3FileDriver) GetFullPath(relativePath string) string {
	if relativePath == "" {
		return ""
	}

	// Clean the relative path
	relativePath = strings.TrimPrefix(relativePath, "/")
	relativePath = filepath.Clean(relativePath)

	return d.basePath + relativePath
}

// GetBucket returns the S3 bucket name
func (d *S3FileDriver) GetBucket() string {
	return d.bucket
}
