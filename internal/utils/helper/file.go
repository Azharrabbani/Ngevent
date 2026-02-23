package helper

import (
	"errors"
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

func ValidatePDF(file multipart.FileHeader) error {
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
		path = filepath.Join(filePath, newFileName)

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

	// ================= PDF =================
	newFileName = baseName + ext
	path = filepath.Join(filePath, newFileName)

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

	return dstFile.Name(), nil
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
