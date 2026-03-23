package helper

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/ryanbekhen/go-webp"
)

func ValidateImage(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		return nil
	}

	return errors.New("image type must be .png/.jpg/.jpeg")
}

func ValidatePDF(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	if ext == ".pdf" {
		return nil
	}

	return errors.New("file type must be .pdf")
}

func SaveToLocal(file *multipart.FileHeader, filePath string) (string, string, error) {
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		return "", "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Get filename without its extension
	baseName := strings.TrimSuffix(file.Filename, ext)

	fileName, err := TransformFileName(baseName)
	if err != nil {
		return "", "", err
	}

	// ============= IMAGE ================
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		path, newFileName, err := SaveImage(file, fileName, filePath)
		if err != nil {
			return "", "", err
		}

		return path, newFileName, nil
	}

	// ================= PDF =================
	path, newFileName, err := SavePDF(file, fileName, filePath, ext)
	if err != nil {
		return "", "", err
	}

	return path, newFileName, nil
}

func SaveImage(file *multipart.FileHeader, fileName, filePath string) (string, string, error) {
	// Validate the imaage
	if err := ValidateImage(file); err != nil {
		return "", "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return "", "", err
	}

	newFileName := fileName + ".webp"
	path := filepath.Join(filePath, newFileName)

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

func SavePDF(file *multipart.FileHeader, fileName, filePath, ext string) (string, string, error) {
	if err := ValidatePDF(file); err != nil {
		return "", "", err
	}

	newFileName := fileName + ext
	path := filepath.Join(filePath, newFileName)

	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	out, err := os.Create(path)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", "", err
	}

	return path, newFileName, nil
}

func TransformFileName(baseName string) (string, error) {
	if baseName == "" {
		return "", errors.New("file name empyt")
	}

	baseName = fmt.Sprintf("%s%s", baseName, uuid.New().String())

	// remove spacing in the first and the last
	baseName = strings.TrimSpace(baseName)

	// Separation based on any number of spaces
	parts := strings.Fields(baseName)

	// Join with "_"
	newName := strings.Join(parts, "_")

	return newName, nil
}

func CopyFile(srcPath, dstPath string) (string, error) {
	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	// Make sure the destination folder exist
	if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
		return "", err
	}

	// Make destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	// Copy source file
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return "", err
	}

	// Sync to disk
	if err := dstFile.Sync(); err != nil {
		return "", err
	}

	return filepath.Base(dstPath), nil
}

func CompressPDF(input, output string) error {
	cmd := exec.Command(
		"gs",
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dPDFSETTINGS=/ebook",
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		"-sOutputFile="+output,
		input,
	)

	return cmd.Run()
}

func OptimizePDF(path string) error {
	conf := model.NewDefaultConfiguration()
	conf.Optimize = true

	return api.OptimizeFile(path, "", conf)
}
