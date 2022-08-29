package mime

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"golang.org/x/exp/slices"
)

func Get(file multipart.File) (string, error) {
	// Only allocate 512 bytes since DetectContentType only reads that amount
	buffer := bytes.NewBuffer(make([]byte, 0, 512))

	if _, err := io.Copy(buffer, file); err != nil {
		return "", err
	}

	return http.DetectContentType(buffer.Bytes()), nil
}

func Contains(file multipart.File, contentTypes []string) bool {
	contentType, err := Get(file)
	if err != nil {
		return false
	}

	return slices.Contains(contentTypes, contentType)
}
