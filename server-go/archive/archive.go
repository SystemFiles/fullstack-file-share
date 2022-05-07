package archive

import (
	"archive/tar"
	"os"
	"path"
	"path/filepath"

	"sykesdev.ca/file-server/config"
)

type Archive struct {
	uploadId string
	path string
	size int
}

func New(uploadId string) (*Archive, error) {
	cfg := config.Get()

	dest, err := filepath.Abs(path.Join(cfg.StorageDir, uploadId + ".tar"))
	if err != nil {
		return nil, err
	}

	return &Archive{
		uploadId: uploadId,
		path: dest,
		size: 0,
	}, nil
}

func (a *Archive) Path() string {
	return a.path
}

func (a *Archive) UploadID() string {
	return a.uploadId
}

func (a *Archive) Size() int {
	return a.size
}

func (a *Archive) AddFile(fileName string, data []byte) error {
	var destFile *os.File
	_, err := os.Stat(a.path)
	if err != nil {
		destFile, err = os.Create(a.path)
		if err != nil {
			return err
		}
	} else {
		destFile, err = os.OpenFile(a.path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}
	}
	defer destFile.Close()

	hdr := tar.Header{
		Name: fileName,
		Mode: 0644,
		Size: int64(len(data)),
	}

	if a.size > 0 {
		if _, err = destFile.Seek(-2<<9, os.SEEK_END); err != nil {
			return err
		}
	}

	t := tar.NewWriter(destFile)
	if err := t.WriteHeader(&hdr); err != nil {
		return err
	}
	_, err = t.Write(data)
	if err != nil {
		return err
	}

	if err := t.Close(); err != nil {
		return err
	}
	
	a.size += len(data)

	return nil
}