package storage

import "mime/multipart"

// Provider defines the interface for File Storage operations.
type Provider interface {
	// Upload saves the file to the specified directory and returns the saved path (or object key).
	Upload(file *multipart.FileHeader, directory string, filename string) (string, error)

	// GetURL transforms a relative path or object key into an absolute, public-facing URL.
	GetURL(path string) string
}
