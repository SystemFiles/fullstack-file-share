package routes

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/golang/glog"
	"github.com/labstack/echo/v4"
	"sykesdev.ca/file-server/provider"
)

type FileUploadResponse struct {
	Message string
	Count int
	Files []string
}

func UploadFiles(c echo.Context) error {
	var wg sync.WaitGroup
	var fileIds []string

	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	files := form.File["files"]

	wg.Add(len(files))
	idChan, errChan := make(chan string, len(files)), make(chan error)
	for _, file := range files {
		glog.Infof("uploading %s", file.Filename)
		go provider.UploadFile(&wg, file, idChan, errChan)
	}

	go func() {
		for fid := range idChan {
			glog.Infof("processing file upload for %s", fid)
			fileIds = append(fileIds, fid)
		}
	}()

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

	glog.Infof("Uploaded %d files successfully", len(fileIds))
	return c.JSONPretty(http.StatusCreated, FileUploadResponse{
		Message: fmt.Sprintf("Uploaded %d files successfully", len(fileIds)),
		Count: len(fileIds),
		Files: fileIds,
	}, " ")
}