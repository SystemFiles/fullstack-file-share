package provider

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"sykesdev.ca/file-server/config"
)

func UploadFile(file *multipart.FileHeader, idChan chan string, errChan chan error) {
	cfg := config.Get()
	fileId, err := uuid.NewUUID()
	if err != nil {
		errChan <- err
	}

	src, err := file.Open()
	if err != nil {
		errChan <- err
	}
	defer src.Close()

	fileParts := strings.Split(file.Filename, ".")
	dstPath, err := filepath.Abs(path.Join(cfg.StorageDir, fileId.String() + "." + fileParts[len(fileParts) -1]))
	if err != nil {
		errChan <- err
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		errChan <- err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		errChan <- err
	}

	idChan <- fileId.String()
}