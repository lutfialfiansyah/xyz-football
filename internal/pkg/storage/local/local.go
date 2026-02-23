package local

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	AppURL string
}

// NewLocalStorage creates a new LocalStorage provider with the given base AppURL.
func NewLocalStorage(appURL string) *LocalStorage {
	// Ensure AppURL does not end with a slash for consistent formatting
	return &LocalStorage{
		AppURL: strings.TrimSuffix(appURL, "/"),
	}
}

// Upload saves a file to the local disk and returns its relative path with a leading slash.
func (s *LocalStorage) Upload(file *multipart.FileHeader, directory string, filename string) (string, error) {
	// Ensure the target directory exists
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	dstPath := filepath.Join(directory, filename)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	out, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", fmt.Errorf("failed to write to destination file: %w", err)
	}

	// Always return path with forward slashes for URLs
	return "/" + strings.ReplaceAll(dstPath, "\\", "/"), nil
}

// GetURL converts a relative path to a full public URL using the configured AppURL.
func (s *LocalStorage) GetURL(path string) string {
	if path == "" {
		return ""
	}

	// If the path is already an absolute HTTP URL, return it as is.
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	// Ensure the path has a leading slash before appending
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return s.AppURL + path
}
