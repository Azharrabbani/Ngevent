package helper

import (
	"errors"
	"image"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/ryanbekhen/go-webp"
)

func ValidateImage(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	if ext == ".png" || ext == ".jpg" || ext == "jpeg" {
		return nil
	}

	return errors.New("image type must be .png/.jpg/.jpeg")
}

func SaveToLocal(file *multipart.FileHeader, filePath string) (string, string, error) {
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		return "", "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	baseName := uuid.NewString()

	var (
		path        string
		newFileName string
	)

	// ============= IMAGE ================
	if ext == ".png" || ext == ".jpg" || ext == "jpeg" {
		src, err := file.Open()
		if err != nil {
			return "", "", err
		}

		img, _, err := image.Decode(src)
		if err != nil {
			return "", "", err
		}

		newFileName = baseName + ".webp"
		path := filepath.Join(path, newFileName)

		out, err := os.Create(path)
		if err != nil {
			return "", "", err
		}
		defer out.Close()

		if file.Size > 1*1024*1024 {
			err = webp.Encode(img, 75, out)
			if err != nil {
				return "", "", err
			}
		} else {
			err = webp.Encode(img, 100, out)
			if err != nil {
				return "", "", err
			}
		}

		return path, newFileName, nil
	}

	return "", "", errors.New("failed to save photo")
}
