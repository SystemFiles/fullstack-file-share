package provider

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"sykesdev.ca/file-server/config"
)

func UploadFile(wg *sync.WaitGroup, file *multipart.FileHeader, idChan chan string, errChan chan error) {
	defer wg.Done()

	cfg := config.Get()
	fileId, err := uuid.NewRandom()
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