package routes

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"sykesdev.ca/file-server/archive"
	"sykesdev.ca/file-server/provider"
)

type FileUploadResponse struct {
	Message string
	Count int
	UploadId string
}

func UploadFiles(c echo.Context) error {
	var wg sync.WaitGroup

	uploadId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	uploadArchive, err := archive.New(uploadId.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	files := form.File["files"]

	wg.Add(len(files))
	errChan := make(chan error)
	for _, file := range files {
		glog.Infof("uploading %s", file.Filename)
		go provider.UploadLocal(&wg, uploadArchive, file, errChan)
	}

	go func() {
		for e := range errChan {
			err = e
		}
	}()

	wg.Wait()

	if err != nil {
		glog.Errorf("error occurred when uploading files. %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	glog.Infof("Uploaded %d files successfully", len(files))
	return c.JSONPretty(http.StatusCreated, FileUploadResponse{
		Message: fmt.Sprintf("Uploaded %d files successfully", len(files)),
		Count: len(files),
		UploadId: uploadArchive.UploadID(),
	}, " ")
}