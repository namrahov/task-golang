package mapper

import (
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"task-golang/model"
	"time"
)

func BuildAttachmentFile(multipartFileHeader *multipart.FileHeader, minioBucket string) (*model.AttachmentFileDto, error) {
	// Validate file name
	originalFilename := sanitizeFileName(multipartFileHeader.Filename)
	if originalFilename == "" {
		return nil, fmt.Errorf("invalid file name")
	}

	// Generate a unique name for the file
	uniqueName := generateUniqueName(originalFilename)

	filePath := path.Join(minioBucket, uniqueName)

	attachmentFile := model.AttachmentFile{
		FileName:  originalFilename,
		FilePath:  filePath,
		FileType:  multipartFileHeader.Header.Get("Content-Type"),
		CreatedAt: time.Now(),
	}

	return &model.AttachmentFileDto{
		AttachmentFile: attachmentFile,
		UniqueName:     uniqueName,
	}, nil
}

// RESERVED_NAMES is a set of reserved file names on Windows
var RESERVED_NAMES = map[string]bool{
	"CON": true, "PRN": true, "AUX": true, "NUL": true,
	"COM1": true, "COM2": true, "COM3": true, "COM4": true,
	"COM5": true, "COM6": true, "COM7": true, "COM8": true, "COM9": true,
	"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true,
	"LPT5": true, "LPT6": true, "LPT7": true, "LPT8": true, "LPT9": true,
}

// sanitizeFileName sanitizes the file name
func sanitizeFileName(fileName string) string {
	if fileName == "" {
		return "default_file"
	}

	// Remove any path components
	fileName = filepath.Base(fileName)

	// Remove any non-alphanumeric characters except for dots, dashes, and underscores
	reg := regexp.MustCompile(`[^a-zA-Z0-9.-_]`)
	fileName = reg.ReplaceAllString(fileName, "_")

	// Limit the length to 255 characters
	maxLength := 255
	if len(fileName) > maxLength {
		extension := filepath.Ext(fileName) // Get the file extension
		baseName := fileName[:len(fileName)-len(extension)]
		if len(baseName) > maxLength-len(extension) {
			baseName = baseName[:maxLength-len(extension)]
		}
		fileName = baseName + extension
	}

	// Check for reserved names
	baseName := strings.ToUpper(strings.TrimSuffix(fileName, filepath.Ext(fileName)))
	if RESERVED_NAMES[baseName] {
		fileName = "_" + fileName
	}

	// Ensure non-empty output
	if fileName == "" {
		fileName = "default_file"
	}

	return fileName
}

// generateUniqueName generates a unique file name
func generateUniqueName(originalFilename string) string {
	// Use UUID and current timestamp to ensure uniqueness
	uuid := uuid.New().String()
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d_%s", uuid, timestamp, originalFilename)
}
