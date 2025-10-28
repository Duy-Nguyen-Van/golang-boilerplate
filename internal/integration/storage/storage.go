package storage

import (
	"context"
	"fmt"
	"golang-boilerplate/internal/config"
	"golang-boilerplate/internal/constants"
	"golang-boilerplate/internal/errors"
	"mime/multipart"
	"time"
)

// UploadResult represents the result of a single file upload to a storage backend
type UploadResult struct {
	URL      string `json:"url"`
	Key      string `json:"key"`
	Bucket   string `json:"bucket"`
	Location string `json:"location"`
}

// BatchUploadResult represents the result of a multiple file upload operation
type BatchUploadResult struct {
	Files []UploadResult
}

// StorageAdapter defines the interface for storage operations
type StorageAdapter interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, key string) (*UploadResult, error)
	UploadFiles(ctx context.Context, files []*multipart.FileHeader) (*BatchUploadResult, error)
	GetObjectURL(key string) string
	GetPresignedURL(ctx context.Context, key string, duration ...time.Duration) (string, error)
}

func ProvideStorageAdapter(config *config.Config) (StorageAdapter, error) {
	switch config.StorageProvider {
	case constants.StorageProviderGCS:
		gcsAdapter, err := NewGCSAdapter(config)
		if err != nil {
			return nil, errors.ExternalServiceError("Failed to initialize GCS storage adapter", err).
				WithOperation("initialize_storage_adapter").
				WithResource("storage")
		}
		return gcsAdapter, nil
	default:
		return nil, errors.InternalError("Invalid storage provider", fmt.Errorf("invalid storage provider: %s", config.StorageProvider)).
			WithOperation("initialize_storage_adapter").
			WithResource("storage").
			WithContext("storage_provider", config.StorageProvider)
	}
}
