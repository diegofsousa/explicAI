package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/diegofsousa/explicAI/internal/application"
	"github.com/diegofsousa/explicAI/internal/application/service"
	"github.com/diegofsousa/explicAI/internal/infrastructure/errors"
	"github.com/diegofsousa/explicAI/internal/infrastructure/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ExplicaServer struct {
	summary service.SummaryUseCase
}

func NewExplicaServer(summary service.SummaryUseCase) *ExplicaServer {
	return &ExplicaServer{
		summary: summary,
	}
}

func (api *ExplicaServer) Register(server *echo.Echo) {
	server.POST("/upload", api.Upload)
	server.GET("/summaries", api.ListSummaries)
	server.GET("/summaries/:externalId", api.GetSummaryByExternalID)
	server.DELETE("/summaries/:externalId", api.DeleteSummaryByExternalID)
}

func (api *ExplicaServer) Upload(c echo.Context) error {
	ctx := c.Request().Context()
	file, err := api.getFileFromRequest(ctx, c)
	if err != nil {
		return errors.Handle(c, err)
	}

	result, err := api.summary.CreateSummaryAndTriggerAIProccess(ctx, file)
	if err != nil {
		return errors.Handle(c, err)
	}

	return c.JSON(http.StatusCreated, result)
}

func (api *ExplicaServer) ListSummaries(c echo.Context) error {
	ctx := c.Request().Context()

	result, err := api.summary.ListSummaries(ctx)
	if err != nil {
		return errors.Handle(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (api *ExplicaServer) GetSummaryByExternalID(c echo.Context) error {
	ctx := c.Request().Context()
	externalID := c.Param("externalId")

	parsedExternalID, err := uuid.Parse(externalID)
	if err != nil {
		return errors.Handle(c, application.ExternalIDIsInvalid)
	}

	result, err := api.summary.GetSummaryByExternalID(ctx, parsedExternalID)
	if err != nil {
		return errors.Handle(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (api *ExplicaServer) DeleteSummaryByExternalID(c echo.Context) error {
	ctx := c.Request().Context()
	externalID := c.Param("externalId")

	parsedExternalID, err := uuid.Parse(externalID)
	if err != nil {
		return errors.Handle(c, application.ExternalIDIsInvalid)
	}

	err = api.summary.DeleteSummaryByExternalID(ctx, parsedExternalID)
	if err != nil {
		return errors.Handle(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "summary has been removed",
	})
}

func (api *ExplicaServer) getFileFromRequest(ctx context.Context, c echo.Context) ([]byte, error) {
	file, err := c.FormFile("file")
	if err != nil {
		log.LogError(ctx, "missing file", err)
		return nil, application.MissingFile
	}

	allowedExtensions := map[string]bool{
		".mp3":  true,
		".mp4":  true,
		".mpeg": true,
		".mpga": true,
		".m4a":  true,
		".wav":  true,
		".webm": true,
	}

	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[fileExtension] {
		return nil, application.InvalidFile
	}

	src, err := file.Open()
	if err != nil {
		log.LogError(ctx, "fail to open file", err)
		return nil, application.FailedReadFile
	}
	defer src.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, src); err != nil {
		log.LogError(ctx, "fail to read file", err)
		return nil, application.FailedReadFile
	}

	return buf.Bytes(), nil
}
