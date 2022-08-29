package mime

import (
	"mime/multipart"
	"net/http"

	"golang.org/x/exp/slices"
)

func Get(file multipart.File) (string, error) {
	// Only allocate 512 bytes since DetectContentType only reads that amount
	buffer := make([]byte, 512)

	if _, err := file.Read(buffer); err != nil {
		return "", err
	}

	_, err := file.Seek(0, 0) // Reset the file position after reading the first 512 bytes
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

func Contains(file multipart.File, contentTypes []string) bool {
	contentType, err := Get(file)
	if err != nil {
		return false
	}

	return slices.Contains(contentTypes, contentType)
}
