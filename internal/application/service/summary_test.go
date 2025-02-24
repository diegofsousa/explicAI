package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/diegofsousa/explicAI/internal/application"
	gatewaymocks "github.com/diegofsousa/explicAI/internal/gateway/mocks"
	"github.com/diegofsousa/explicAI/internal/gateway/repository"
	"github.com/diegofsousa/explicAI/internal/gateway/summarize"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	timeLayout            = "2006-01-02 15:04:05"
	summaryExternalIDStr  = "9156fe73-c692-4834-bf58-7474b878a634"
	summaryExternalIDUUID = uuid.MustParse(summaryExternalIDStr)
	createdAt, _          = time.Parse(timeLayout, "2025-01-25 15:04:05")
	textTranscribed       = "result text transcribed"
	title                 = "title"
	description           = "description"
	briefResume           = "brief resume"
	mediumResume          = "medium resume"
	fulltext              = "full text"
)

type (
	SummaryTestSuite struct {
		suite.Suite

		ctx             context.Context
		audioTranscript *gatewaymocks.AudioTranscript
		summarize       *gatewaymocks.Summarize
		repository      *gatewaymocks.Repository
	}
)

func TestSummaryTestSuite(t *testing.T) {
	suite.Run(t, new(SummaryTestSuite))
}

func (s *SummaryTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *SummaryTestSuite) TearDownTest() {
	mock.AssertExpectationsForObjects(s.T())
}

func (s *SummaryTestSuite) TestSummaryCreate() {
	s.Run("successful create summary", func() {
		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().
			CreateSummary(mock.Anything, repository.ReceivedFile).
			Return(&repository.SummaryCreateOutput{
				ExternalID: summaryExternalIDUUID,
				Status:     repository.StatusToString[repository.ReceivedFile].Status,
				Progress: sql.NullInt32{
					Int32: 33,
					Valid: true,
				},
				CreatedAt: createdAt,
			}, nil)

		s.audioTranscript = new(gatewaymocks.AudioTranscript)
		s.audioTranscript.EXPECT().
			Transcribe(mock.Anything, mock.Anything).
			Return(&textTranscribed, nil)

		ustInput := repository.SummaryUpdateTranscribedInput{
			ExternalID: summaryExternalIDUUID,
			Status:     repository.Trancribed,
		}

		s.repository.EXPECT().
			UpdateSummaryTranscribed(mock.Anything, ustInput).
			Return(nil)

		s.summarize = new(gatewaymocks.Summarize)
		s.summarize.EXPECT().
			Resume(mock.Anything, textTranscribed).
			Return(&summarize.ResumeOutput{
				Title:        title,
				Description:  description,
				BriefResume:  briefResume,
				MediumResume: mediumResume,
			}, nil)

		s.summarize.EXPECT().
			FullTextOrganize(mock.Anything, textTranscribed).
			Return(&fulltext, nil)

		susInput := repository.SummaryUpdateSummarizedInput{
			ExternalID:   summaryExternalIDUUID,
			Status:       repository.Summarized,
			Title:        title,
			Description:  description,
			BriefResume:  briefResume,
			MediumResume: mediumResume,
			FullText:     fulltext,
		}

		s.repository.EXPECT().
			UpdateSummarySummarized(mock.Anything, susInput).
			Return(nil)

		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		output, err := service.CreateSummaryAndTriggerAIProccess(s.ctx, []byte{})
		s.Require().NoError(err)
		s.Equal("RECEIVED_FILE", output.Status)
		s.Equal(createdAt, output.CreatedAt)
	})

	s.Run("fail db create summary", func() {
		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().
			CreateSummary(mock.Anything, repository.ReceivedFile).
			Return(nil, errors.New("some error"))

		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		_, err := service.CreateSummaryAndTriggerAIProccess(s.ctx, []byte{})
		s.Require().ErrorIs(err, application.InternalDatabaseError)
	})
}

func (s *SummaryTestSuite) TestAIProccessSumary() {
	s.Run("successful proccess summary", func() {
		ctx, cancel := context.WithCancel(context.Background())

		s.audioTranscript = new(gatewaymocks.AudioTranscript)
		s.audioTranscript.EXPECT().
			Transcribe(mock.Anything, mock.Anything).
			Return(&textTranscribed, nil)

		ustInput := repository.SummaryUpdateTranscribedInput{
			ExternalID: summaryExternalIDUUID,
			Status:     repository.Trancribed,
		}

		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().
			UpdateSummaryTranscribed(mock.Anything, ustInput).
			Return(nil)

		s.summarize = new(gatewaymocks.Summarize)
		s.summarize.EXPECT().
			Resume(mock.Anything, textTranscribed).
			Return(&summarize.ResumeOutput{
				Title:        title,
				Description:  description,
				BriefResume:  briefResume,
				MediumResume: mediumResume,
			}, nil)

		s.summarize.EXPECT().
			FullTextOrganize(mock.Anything, textTranscribed).
			Return(&fulltext, nil)

		susInput := repository.SummaryUpdateSummarizedInput{
			ExternalID:   summaryExternalIDUUID,
			Status:       repository.Summarized,
			Title:        title,
			Description:  description,
			BriefResume:  briefResume,
			MediumResume: mediumResume,
			FullText:     fulltext,
		}

		s.repository.EXPECT().
			UpdateSummarySummarized(mock.Anything, susInput).
			Return(nil)

		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		service.AISummaryProccess(ctx, cancel, []byte{}, summaryExternalIDUUID)
		s.audioTranscript.AssertCalled(s.T(), "Transcribe", mock.Anything, mock.Anything)
		s.repository.AssertCalled(s.T(), "UpdateSummaryTranscribed", mock.Anything, ustInput)
		s.summarize.AssertCalled(s.T(), "Resume", mock.Anything, textTranscribed)
		s.summarize.AssertCalled(s.T(), "FullTextOrganize", mock.Anything, textTranscribed)
		s.repository.AssertCalled(s.T(), "UpdateSummarySummarized", mock.Anything, susInput)
	})

	s.Run("fail full text organize summary", func() {
		ctx, cancel := context.WithCancel(context.Background())

		s.audioTranscript = new(gatewaymocks.AudioTranscript)
		s.audioTranscript.EXPECT().
			Transcribe(mock.Anything, mock.Anything).
			Return(&textTranscribed, nil)

		ustInput := repository.SummaryUpdateTranscribedInput{
			ExternalID: summaryExternalIDUUID,
			Status:     repository.Trancribed,
		}

		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().
			UpdateSummaryTranscribed(mock.Anything, ustInput).
			Return(nil)

		s.summarize = new(gatewaymocks.Summarize)
		s.summarize.EXPECT().
			Resume(mock.Anything, textTranscribed).
			Return(&summarize.ResumeOutput{
				Title:        title,
				Description:  description,
				BriefResume:  briefResume,
				MediumResume: mediumResume,
			}, nil)

		s.summarize.EXPECT().
			FullTextOrganize(mock.Anything, textTranscribed).
			Return(&fulltext, errors.New("some error"))

		susInput := repository.SummaryUpdateSummarizedInput{
			ExternalID: summaryExternalIDUUID,
			Status:     repository.SummarizedFailed,
		}

		s.repository.EXPECT().
			UpdateSummarySummarized(mock.Anything, susInput).
			Return(nil)

		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		service.AISummaryProccess(ctx, cancel, []byte{}, summaryExternalIDUUID)
		s.audioTranscript.AssertCalled(s.T(), "Transcribe", mock.Anything, mock.Anything)
		s.repository.AssertCalled(s.T(), "UpdateSummaryTranscribed", mock.Anything, ustInput)
		s.summarize.AssertCalled(s.T(), "Resume", mock.Anything, textTranscribed)
		s.summarize.AssertCalled(s.T(), "FullTextOrganize", mock.Anything, textTranscribed)
		s.repository.AssertCalled(s.T(), "UpdateSummarySummarized", mock.Anything, susInput)
	})

	s.Run("fail resume summary", func() {
		ctx, cancel := context.WithCancel(context.Background())

		s.audioTranscript = new(gatewaymocks.AudioTranscript)
		s.audioTranscript.EXPECT().
			Transcribe(mock.Anything, mock.Anything).
			Return(&textTranscribed, nil)

		ustInput := repository.SummaryUpdateTranscribedInput{
			ExternalID: summaryExternalIDUUID,
			Status:     repository.Trancribed,
		}

		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().
			UpdateSummaryTranscribed(mock.Anything, ustInput).
			Return(nil)

		s.summarize = new(gatewaymocks.Summarize)
		s.summarize.EXPECT().
			Resume(mock.Anything, textTranscribed).
			Return(nil, errors.New("some error"))

		s.summarize.EXPECT().
			FullTextOrganize(mock.Anything, textTranscribed).
			Return(&fulltext, nil)

		susInput := repository.SummaryUpdateSummarizedInput{
			ExternalID: summaryExternalIDUUID,
			Status:     repository.SummarizedFailed,
		}

		s.repository.EXPECT().
			UpdateSummarySummarized(mock.Anything, susInput).
			Return(nil)

		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		service.AISummaryProccess(ctx, cancel, []byte{}, summaryExternalIDUUID)
		s.audioTranscript.AssertCalled(s.T(), "Transcribe", mock.Anything, mock.Anything)
		s.repository.AssertCalled(s.T(), "UpdateSummaryTranscribed", mock.Anything, ustInput)
		s.summarize.AssertCalled(s.T(), "Resume", mock.Anything, textTranscribed)
		s.summarize.AssertCalled(s.T(), "FullTextOrganize", mock.Anything, textTranscribed)
		s.repository.AssertCalled(s.T(), "UpdateSummarySummarized", mock.Anything, susInput)
	})

	s.Run("fail audio transcribe", func() {
		ctx, cancel := context.WithCancel(context.Background())

		s.audioTranscript = new(gatewaymocks.AudioTranscript)
		s.audioTranscript.EXPECT().
			Transcribe(mock.Anything, mock.Anything).
			Return(nil, errors.New("some error"))

		ustInput := repository.SummaryUpdateTranscribedInput{
			ExternalID: summaryExternalIDUUID,
			Status:     repository.TranscribedFailed,
		}

		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().
			UpdateSummaryTranscribed(mock.Anything, ustInput).
			Return(nil)

		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		service.AISummaryProccess(ctx, cancel, []byte{}, summaryExternalIDUUID)
		s.audioTranscript.AssertCalled(s.T(), "Transcribe", mock.Anything, mock.Anything)
		s.repository.AssertCalled(s.T(), "UpdateSummaryTranscribed", mock.Anything, ustInput)
	})
}

func (s *SummaryTestSuite) TestListSummaries() {
	s.Run("successful list summaries", func() {
		summaryExternalID2Str := "7748608c-e9fc-4dfa-a083-6f4014457b8a"
		summaryExternalID2UUID := uuid.MustParse(summaryExternalID2Str)

		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().GetSummaries(mock.Anything).
			Return([]repository.SummaryOutput{
				{
					ExternalID: summaryExternalIDUUID,
					CreatedAt:  createdAt,
					UpdatedAt:  createdAt,
					Status:     "SUMMARIZED",
					Title: sql.NullString{
						String: title,
						Valid:  true,
					},
					Description: sql.NullString{
						String: description,
						Valid:  true,
					},
					BriefResume: sql.NullString{
						String: briefResume,
						Valid:  true,
					},
					MediumResume: sql.NullString{
						String: mediumResume,
						Valid:  true,
					},
					Progress: sql.NullInt32{
						Int32: 100,
						Valid: true,
					},
					FullText: sql.NullString{
						String: fulltext,
						Valid:  true,
					},
				},
				{
					ExternalID: summaryExternalID2UUID,
					CreatedAt:  createdAt,
					Status:     "RECEIVED_FILE",
				},
			}, nil)

		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		output, err := service.ListSummaries(s.ctx)
		s.Require().NoError(err)
		s.Equal(summaryExternalIDUUID, output.Data[0].ExternalID)
		s.Equal("SUMMARIZED", output.Data[0].Status)
		s.Equal(createdAt, output.Data[0].CreatedAt)
		s.Equal(title, output.Data[0].Title)
		s.Equal(description, output.Data[0].Description)
		s.Equal(createdAt, output.Data[0].CreatedAt)

		s.Equal(summaryExternalID2UUID, output.Data[1].ExternalID)
		s.Equal("RECEIVED_FILE", output.Data[1].Status)
		s.Equal(createdAt, output.Data[1].CreatedAt)
	})

	s.Run("fail list summaries", func() {
		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().GetSummaries(mock.Anything).
			Return(nil, errors.New("some error"))

		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		_, err := service.ListSummaries(s.ctx)
		s.Require().ErrorIs(err, application.UnexpectedErrorList)
	})
}

func (s *SummaryTestSuite) TestGetSummaryByExternalID() {
	s.Run("successful get summary by external id", func() {
		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().
			GetSummaryByExternalID(mock.Anything, summaryExternalIDUUID).
			Return(&repository.SummaryOutput{
				ExternalID:   summaryExternalIDUUID,
				Status:       repository.StatusToString[repository.Summarized].Status,
				CreatedAt:    createdAt,
				UpdatedAt:    createdAt,
				Progress:     sql.NullInt32{Int32: 100, Valid: true},
				Title:        sql.NullString{String: title, Valid: true},
				Description:  sql.NullString{String: description, Valid: true},
				BriefResume:  sql.NullString{String: briefResume, Valid: true},
				MediumResume: sql.NullString{String: mediumResume, Valid: true},
				FullText:     sql.NullString{String: fulltext, Valid: true},
			}, nil)

		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		output, err := service.GetSummaryByExternalID(s.ctx, summaryExternalIDUUID)
		s.Require().NoError(err)
		s.Equal(summaryExternalIDUUID, output.ExternalID)
		s.Equal(repository.StatusToString[repository.Summarized].Status, output.Status)
		s.Equal(title, output.Title)
		s.Equal(description, output.Description)
		s.Equal(briefResume, output.BriefResume)
		s.Equal(mediumResume, output.MediumResume)
		s.Equal(fulltext, output.FullText)
	})

	s.Run("error getting summary", func() {
		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().
			GetSummaryByExternalID(mock.Anything, summaryExternalIDUUID).
			Return(nil, errors.New("some error"))

		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		_, err := service.GetSummaryByExternalID(s.ctx, summaryExternalIDUUID)
		s.Require().Error(err)
	})

	s.Run("get summary not found", func() {
		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().
			GetSummaryByExternalID(mock.Anything, summaryExternalIDUUID).
			Return(nil, application.SummaryNotFound)

		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		_, err := service.GetSummaryByExternalID(s.ctx, summaryExternalIDUUID)
		s.Require().ErrorIs(err, application.SummaryNotFound)
	})
}

func (s *SummaryTestSuite) TestDeleteSummaryByExternalID() {
	s.Run("successful delete summary by external id", func() {
		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().
			DeleteSummaryByExternalID(s.ctx, summaryExternalIDUUID).
			Return(nil)
		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		err := service.DeleteSummaryByExternalID(s.ctx, summaryExternalIDUUID)
		s.Require().NoError(err)
	})

	s.Run("delete summary not found", func() {
		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().
			DeleteSummaryByExternalID(s.ctx, summaryExternalIDUUID).
			Return(application.SummaryNotFound)
		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		err := service.DeleteSummaryByExternalID(s.ctx, summaryExternalIDUUID)
		s.Require().ErrorIs(err, application.SummaryNotFound)
	})

	s.Run("error deleting summary", func() {
		s.repository = new(gatewaymocks.Repository)
		s.repository.EXPECT().
			DeleteSummaryByExternalID(s.ctx, summaryExternalIDUUID).
			Return(errors.New("some error"))
		service := NewSummary(s.audioTranscript, s.summarize, s.repository)
		err := service.DeleteSummaryByExternalID(s.ctx, summaryExternalIDUUID)
		s.Require().Error(err)
	})
}
