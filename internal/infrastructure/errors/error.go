package errors

import (
	"github.com/pkg/errors"

	"github.com/diegofsousa/explicAI/internal/application"
	"github.com/labstack/echo/v4"
)

func Handle(c echo.Context, err error) error {
	switch errors.Cause(err) {
	case application.MissingFile, application.InvalidFile, application.ExternalIDIsInvalid:
		return echo.ErrBadRequest
	case application.SummaryNotFound:
		return echo.ErrNotFound
	case application.FailedReadFile:
		return echo.ErrUnprocessableEntity
	default:
		return echo.ErrInternalServerError
	}
}
