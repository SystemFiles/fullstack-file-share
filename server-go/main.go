package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator"
	"github.com/golang/glog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"sykesdev.ca/file-server/config"
	"sykesdev.ca/file-server/routes"
	"sykesdev.ca/file-server/version"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (c *CustomValidator) Validate(i interface{}) error {
	if err := c.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return nil
}

func root() error {
	e := echo.New()

	cfg := config.Get()

	// echo configuration
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:  []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	// request validator instantiation
	e.Validator = &CustomValidator{validator: validator.New()}

	// route definitions
	api := e.Group(fmt.Sprintf("/api/v%s", strings.Split(version.Version, ".")[0]))
	api.GET("/", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, routes.CustomResponse{Message: fmt.Sprintf("Fileshare API Version: %s", version.Version)}, " ")
	})
	api.POST("/files", routes.UploadFiles)

	e.Logger.Info(e.Start(fmt.Sprintf(":%d", cfg.Port)))
	return nil
}

func main() {
	flag.Parse()
	
	if err := root(); err != nil {
		glog.Errorf("api error. %s", err)
		os.Exit(1)
	}
}