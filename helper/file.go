package helper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func UploadFile(r *http.Request, fieldName string, uploadDir string) (string, string, string, error) {
	file, handler, err := r.FormFile(fieldName)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", "", "", nil
		}

		return "", "", "", fmt.Errorf("error Retrieving File From : %w", err)
	}

	defer file.Close()

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to create upload directory: %w", err)
		}
	}

	ext := filepath.Ext(handler.Filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, newFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create file on server: %w", err)
	}

	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", "", "", fmt.Errorf("failed to copy file content: %w", err)
	}

	mimetype, err := GetMimeType(filePath)
	if err != nil {
		mimetype = handler.Header.Get("Content-Type")
		if mimetype == "" {
			mimetype = "application/octet-stream"
		}
	}

	fileURL := fmt.Sprintf("/uploads/%s", newFileName)

	return fileURL, handler.Filename, mimetype, nil
}

func GetMimeType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file to detect mime: %w", err)
	}

	defer file.Close()

	buffer := make([]byte, 512)

	_, err = file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read content file to detect mime type: %w", err)
	}

	contentType := http.DetectContentType(buffer)
	return contentType, nil
}
