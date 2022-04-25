package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"sykesdev.ca/file-server/provider"
)

type FileUploadResponse struct {
	Message string
	Count int
	Files []string
}

func UploadFiles(c echo.Context) error {
	var fileIds []string

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	idChan, errChan := make(chan string, len(files)), make(chan error, 1)
	for _, file := range files {
		go provider.UploadFile(file, idChan, errChan)
	}

	select {
	case id := <- idChan:
		fileIds = append(fileIds, id)
	case err := <- errChan:
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSONPretty(http.StatusCreated, FileUploadResponse{
		Message: fmt.Sprintf("Uploaded %d files successfully", len(fileIds)),
		Count: len(fileIds),
		Files: fileIds,
	}, " ")
}