package storage

import "io"

// Storage interface defines the methods that any storage provider must implement
type Storage interface {
	// Upload uploads a file from srcPath to the storage destination
	Upload(srcPath string, destPath string) error

	// Download downloads a file from the storage source to the local destPath
	Download(srcPath string, destPath string) error

	// List lists files in the storage directory
	List(path string) ([]string, error)

	// Delete deletes a file from storage
	Delete(path string) error
	
	// GetReader returns a reader for a file in storage
	GetReader(path string) (io.ReadCloser, error)
}
