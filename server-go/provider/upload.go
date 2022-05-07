package provider

import (
	"io/ioutil"
	"mime/multipart"
	"sync"

	"sykesdev.ca/file-server/archive"
)

func UploadLocal(wg *sync.WaitGroup, fileArchive *archive.Archive, file *multipart.FileHeader, errChan chan error) {
	defer wg.Done()

	srcData, err := file.Open()
	if err != nil {
		errChan <- err
	}
	defer srcData.Close()

	data, err := ioutil.ReadAll(srcData)
	if err != nil {
		errChan <- err
	}
	if err := fileArchive.AddFile(file.Filename, data); err != nil {
		errChan <- err
	}
}