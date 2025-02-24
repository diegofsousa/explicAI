package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diegofsousa/explicAI/internal/application"
	"github.com/diegofsousa/explicAI/internal/application/service"
	servicemocks "github.com/diegofsousa/explicAI/internal/application/service/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	timeLayout            = "2006-01-02 15:04:05"
	summaryExternalIDStr  = "9156fe73-c692-4834-bf58-7474b878a634"
	summaryExternalIDUUID = uuid.MustParse(summaryExternalIDStr)
	createdAt, _          = time.Parse(timeLayout, "2025-01-25 15:04:05")
)

type (
	ControllerTestSuite struct {
		suite.Suite
		ctx     context.Context
		summary *servicemocks.SummaryUseCase
	}
)

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}

func (s *ControllerTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *ControllerTestSuite) TearDownTest() {
	mock.AssertExpectationsForObjects(s.T())
}

func (s *ControllerTestSuite) TestUpload() {
	s.Run("successful create summary", func() {
		e := echo.New()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test.mp3")
		s.Require().NoError(err)
		_, err = part.Write([]byte("test file content"))
		s.Require().NoError(err)
		writer.Close()

		request := httptest.NewRequest(http.MethodPost, "/upload", body)
		request.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/upload")

		s.summary = new(servicemocks.SummaryUseCase)
		s.summary.EXPECT().
			CreateSummaryAndTriggerAIProccess(mock.Anything, mock.Anything).
			Return(&service.SummarySimpleOutput{
				ExternalID: summaryExternalIDUUID,
				Status:     "RECEIVED_FILE",
				CreatedAt:  createdAt,
				UpdatedAt:  createdAt,
				Progress:   33,
			}, nil)

		handler := NewExplicaServer(s.summary)
		handler.Register(ctx.Echo())
		e.ServeHTTP(recorder, request)

		var response service.SummarySimpleOutput
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		s.Require().NoError(err)

		s.Equal(http.StatusCreated, recorder.Code)
		s.Equal(summaryExternalIDUUID, response.ExternalID)
		s.Equal("RECEIVED_FILE", response.Status)
		s.Equal(createdAt, response.CreatedAt)
		s.Equal(33, response.Progress)
	})

	s.Run("failed create summary", func() {
		e := echo.New()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test.mp3")
		s.Require().NoError(err)
		_, err = part.Write([]byte("test file content"))
		s.Require().NoError(err)
		writer.Close()

		request := httptest.NewRequest(http.MethodPost, "/upload", body)
		request.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/upload")

		s.summary = new(servicemocks.SummaryUseCase)
		s.summary.EXPECT().
			CreateSummaryAndTriggerAIProccess(mock.Anything, mock.Anything).
			Return(nil, errors.New("some error"))

		handler := NewExplicaServer(s.summary)
		handler.Register(ctx.Echo())
		e.ServeHTTP(recorder, request)
		handler.Upload(ctx)

		var response service.SummarySimpleOutput
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		s.Require().NoError(err)

		s.Equal(http.StatusInternalServerError, recorder.Code)
	})

	s.Run("invalid file", func() {
		e := echo.New()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test.txt")
		s.Require().NoError(err)
		_, err = part.Write([]byte("test file content"))
		s.Require().NoError(err)
		writer.Close()

		request := httptest.NewRequest(http.MethodPost, "/upload", body)
		request.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/upload")

		handler := NewExplicaServer(s.summary)
		handler.Register(ctx.Echo())
		e.ServeHTTP(recorder, request)
		handler.Upload(ctx)

		var response service.SummarySimpleOutput
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		s.Require().NoError(err)

		s.Equal(http.StatusBadRequest, recorder.Code)
	})

	s.Run("failed missing file", func() {
		e := echo.New()

		request := httptest.NewRequest(http.MethodPost, "/upload", nil)
		request.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/upload")

		handler := NewExplicaServer(s.summary)
		handler.Register(e)
		e.ServeHTTP(recorder, request)
		handler.Upload(ctx)

		var response service.SummarySimpleOutput
		json.Unmarshal(recorder.Body.Bytes(), &response)

		s.Equal(http.StatusBadRequest, recorder.Code)
	})
}

func (s *ControllerTestSuite) TestListSummaries() {
	s.Run("successful list summaries", func() {
		e := echo.New()

		request := httptest.NewRequest(http.MethodGet, "/summaries", nil)
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/summaries")

		expected := []service.SummarySimpleOutput{
			{
				ExternalID:  summaryExternalIDUUID,
				Status:      "SUMMARIZED",
				CreatedAt:   createdAt,
				UpdatedAt:   createdAt,
				Progress:    100,
				Title:       "title",
				Description: "description",
			},
		}

		s.summary = new(servicemocks.SummaryUseCase)
		s.summary.EXPECT().
			ListSummaries(mock.Anything).
			Return(&service.SummaryListOutput{Data: expected}, nil)

		handler := NewExplicaServer(s.summary)
		handler.Register(ctx.Echo())
		e.ServeHTTP(recorder, request)

		var response service.SummaryListOutput
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		s.Require().NoError(err)

		s.Equal(http.StatusOK, recorder.Code)
		s.Len(response.Data, 1)
		s.Equal(summaryExternalIDUUID, response.Data[0].ExternalID)
		s.Equal("SUMMARIZED", response.Data[0].Status)
		s.Equal("title", response.Data[0].Title)
		s.Equal("description", response.Data[0].Description)
		s.Equal(createdAt, response.Data[0].CreatedAt)
		s.Equal(createdAt, response.Data[0].UpdatedAt)
		s.Equal(100, response.Data[0].Progress)
	})

	s.Run("failed to list summaries", func() {
		e := echo.New()

		request := httptest.NewRequest(http.MethodGet, "/summaries", nil)
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/summaries")

		s.summary = new(servicemocks.SummaryUseCase)
		s.summary.EXPECT().
			ListSummaries(mock.Anything).
			Return(nil, errors.New("some error"))

		handler := NewExplicaServer(s.summary)
		handler.Register(ctx.Echo())
		e.ServeHTTP(recorder, request)

		s.Equal(http.StatusInternalServerError, recorder.Code)
	})
}

func (s *ControllerTestSuite) TestGetSummaryByExternalID() {
	s.Run("successful get summary", func() {
		e := echo.New()

		request := httptest.NewRequest(http.MethodGet, "/summaries/"+summaryExternalIDStr, nil)
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/summaries/:externalId")
		ctx.SetParamNames("externalId")
		ctx.SetParamValues(summaryExternalIDStr)

		expected := &service.SummaryDetailedOutput{
			ExternalID:   summaryExternalIDUUID,
			Status:       "SUMMARIZED",
			CreatedAt:    createdAt,
			UpdatedAt:    createdAt,
			Progress:     100,
			Title:        "title",
			Description:  "description",
			BriefResume:  "brief",
			MediumResume: "medium",
			FullText:     "full",
		}

		s.summary = new(servicemocks.SummaryUseCase)
		s.summary.EXPECT().
			GetSummaryByExternalID(mock.Anything, summaryExternalIDUUID).
			Return(expected, nil)

		handler := NewExplicaServer(s.summary)
		handler.Register(e)
		e.ServeHTTP(recorder, request)

		var response service.SummaryDetailedOutput
		json.Unmarshal(recorder.Body.Bytes(), &response)

		s.Equal(http.StatusOK, recorder.Code)
		s.Equal(summaryExternalIDUUID, response.ExternalID)
		s.Equal("SUMMARIZED", response.Status)
		s.Equal("title", response.Title)
		s.Equal("description", response.Description)
		s.Equal("brief", response.BriefResume)
		s.Equal("medium", response.MediumResume)
		s.Equal("full", response.FullText)
		s.Equal(createdAt, response.CreatedAt)
		s.Equal(createdAt, response.UpdatedAt)
		s.Equal(100, response.Progress)
	})

	s.Run("invalid external id format", func() {
		e := echo.New()
		invalidID := "invalid-uuid"
		request := httptest.NewRequest(http.MethodGet, "/summaries/"+invalidID, nil)
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/summaries/:externalId")
		ctx.SetParamNames("externalId")
		ctx.SetParamValues(invalidID)

		handler := NewExplicaServer(s.summary)
		handler.Register(e)
		e.ServeHTTP(recorder, request)

		s.Equal(http.StatusBadRequest, recorder.Code)
	})

	s.Run("summary not found", func() {
		e := echo.New()

		request := httptest.NewRequest(http.MethodGet, "/summaries/"+summaryExternalIDStr, nil)
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/summaries/:externalId")
		ctx.SetParamNames("externalId")
		ctx.SetParamValues(summaryExternalIDStr)

		s.summary = new(servicemocks.SummaryUseCase)
		s.summary.EXPECT().
			GetSummaryByExternalID(mock.Anything, summaryExternalIDUUID).
			Return(nil, application.SummaryNotFound)

		handler := NewExplicaServer(s.summary)
		handler.Register(e)
		e.ServeHTTP(recorder, request)

		s.Equal(http.StatusNotFound, recorder.Code)
	})

	s.Run("service error when get summary", func() {
		e := echo.New()

		request := httptest.NewRequest(http.MethodGet, "/summaries/"+summaryExternalIDStr, nil)
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/summaries/:externalId")
		ctx.SetParamNames("externalId")
		ctx.SetParamValues(summaryExternalIDStr)

		s.summary = new(servicemocks.SummaryUseCase)
		s.summary.EXPECT().
			GetSummaryByExternalID(mock.Anything, summaryExternalIDUUID).
			Return(nil, errors.New("service error"))

		handler := NewExplicaServer(s.summary)
		handler.Register(e)
		e.ServeHTTP(recorder, request)

		s.Equal(http.StatusInternalServerError, recorder.Code)
	})
}

func (s *ControllerTestSuite) TestDeleteSummaryByExternalID() {
	s.Run("successful delete summary", func() {
		e := echo.New()

		request := httptest.NewRequest(http.MethodDelete, "/summaries/"+summaryExternalIDStr, nil)
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/summaries/:externalId")
		ctx.SetParamNames("externalId")
		ctx.SetParamValues(summaryExternalIDStr)

		s.summary = new(servicemocks.SummaryUseCase)
		s.summary.EXPECT().
			DeleteSummaryByExternalID(mock.Anything, summaryExternalIDUUID).
			Return(nil)

		handler := NewExplicaServer(s.summary)
		handler.Register(e)
		e.ServeHTTP(recorder, request)

		var response map[string]string
		json.Unmarshal(recorder.Body.Bytes(), &response)
		s.Equal(http.StatusOK, recorder.Code)
	})

	s.Run("invalid external id format", func() {
		e := echo.New()

		invalidID := "invalid-id"

		request := httptest.NewRequest(http.MethodDelete, "/summaries/"+invalidID, nil)
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/summaries/:externalId")
		ctx.SetParamNames("externalId")
		ctx.SetParamValues(summaryExternalIDStr)

		handler := NewExplicaServer(s.summary)
		handler.Register(e)
		e.ServeHTTP(recorder, request)

		s.Equal(http.StatusBadRequest, recorder.Code)
	})

	s.Run("summary not found", func() {
		e := echo.New()

		request := httptest.NewRequest(http.MethodDelete, "/summaries/"+summaryExternalIDStr, nil)
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/summaries/:externalId")
		ctx.SetParamNames("externalId")
		ctx.SetParamValues(summaryExternalIDStr)

		s.summary = new(servicemocks.SummaryUseCase)
		s.summary.EXPECT().
			DeleteSummaryByExternalID(mock.Anything, summaryExternalIDUUID).
			Return(application.SummaryNotFound)

		handler := NewExplicaServer(s.summary)
		handler.Register(e)
		e.ServeHTTP(recorder, request)

		s.Equal(http.StatusNotFound, recorder.Code)
	})

	s.Run("service error when deleting summary", func() {
		e := echo.New()

		request := httptest.NewRequest(http.MethodDelete, "/summaries/"+summaryExternalIDStr, nil)
		recorder := httptest.NewRecorder()

		ctx := e.NewContext(request, recorder)
		ctx.SetPath("/summaries/:externalId")
		ctx.SetParamNames("externalId")
		ctx.SetParamValues(summaryExternalIDStr)

		s.summary = new(servicemocks.SummaryUseCase)
		s.summary.EXPECT().
			DeleteSummaryByExternalID(mock.Anything, summaryExternalIDUUID).
			Return(errors.New("some error"))

		handler := NewExplicaServer(s.summary)
		handler.Register(e)
		e.ServeHTTP(recorder, request)

		s.Equal(http.StatusInternalServerError, recorder.Code)
	})
}
