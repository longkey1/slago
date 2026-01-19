package output

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileWriter writes to a file
type FileWriter struct {
	*JSONWriter
	file *os.File
	path string
}

// NewFileWriter creates a new file writer
func NewFileWriter(path string) (*FileWriter, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return &FileWriter{
		JSONWriter: NewJSONWriter(file, true),
		file:       file,
		path:       path,
	}, nil
}

// Close closes the file
func (fw *FileWriter) Close() error {
	return fw.file.Close()
}

// Path returns the file path
func (fw *FileWriter) Path() string {
	return fw.path
}
