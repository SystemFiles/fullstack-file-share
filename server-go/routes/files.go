package routes

import (
	"fmt"
	"net/http"
	"sync"

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
		return err
	}
	files := form.File["files"]

	wg.Add(len(files))
	idChan, errChan := make(chan string, len(files)), make(chan error)
	for _, file := range files {
		go provider.UploadFile(&wg, file, idChan, errChan)
	}

	go func() {
		for fid := range idChan {
			fileIds = append(fileIds, fid)
		}
		
		err =<- errChan
	}()

	wg.Wait()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSONPretty(http.StatusCreated, FileUploadResponse{
		Message: fmt.Sprintf("Uploaded %d files successfully", len(fileIds)),
		Count: len(fileIds),
		Files: fileIds,
	}, " ")
}