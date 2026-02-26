package domain

// MinioConfig holds MinIO/S3 configuration
type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
	Region          string
}

// NewMinioConfig creates a new MinIO configuration with defaults
func NewMinioConfig(endpoint, accessKey, secretKey, bucket string, useSSL bool) *MinioConfig {
	return &MinioConfig{
		Endpoint:        endpoint,
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
		BucketName:      bucket,
		UseSSL:          useSSL,
		Region:          "us-east-1", // Default region
	}
}
