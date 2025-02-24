package db

import (
	"context"
	"errors"
	"time"

	"github.com/diegofsousa/explicAI/internal/application"
	"github.com/diegofsousa/explicAI/internal/gateway/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Summary struct {
	database *PgConnection
}

func NewSummary(databaseUrl string) *Summary {
	return &Summary{
		database: NewPgConnection(databaseUrl),
	}
}

func (s *Summary) CreateSummary(ctx context.Context, status repository.Status) (*repository.SummaryCreateOutput, error) {
	conn, err := s.database.Connect(ctx)
	if err != nil {
		return nil, err
	}

	defer s.database.Close(ctx, conn)

	externalId := uuid.New()
	now := time.Now()

	query := `
		insert into summaries (external_id, created_at, updated_at, status, progress)
		values ($1, $2, $3, $4, $5)
		returning external_id, created_at, status, progress;
	`

	var output repository.SummaryCreateOutput

	err = conn.QueryRow(ctx, query, externalId, now, now,
		repository.StatusToString[status].Status, repository.StatusToString[status].Percentage).
		Scan(&output.ExternalID, &output.CreatedAt, &output.Status, &output.Progress)

	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (s *Summary) UpdateSummaryTranscribed(ctx context.Context, input repository.SummaryUpdateTranscribedInput) error {
	conn, err := s.database.Connect(ctx)
	if err != nil {
		return err
	}

	defer s.database.Close(ctx, conn)

	now := time.Now()

	query := `update summaries set progress = $2, status = $3, updated_at = $4 where external_id = $1;`

	command, err := conn.Exec(
		ctx,
		query,
		input.ExternalID,
		repository.StatusToString[input.Status].Percentage,
		repository.StatusToString[input.Status].Status,
		now,
	)

	if err != nil {
		return err
	}

	if command.RowsAffected() == 0 {
		return errors.New("register not found")
	}

	return nil
}

func (s *Summary) UpdateSummarySummarized(ctx context.Context, input repository.SummaryUpdateSummarizedInput) error {
	conn, err := s.database.Connect(ctx)
	if err != nil {
		return err
	}

	defer s.database.Close(ctx, conn)

	now := time.Now()

	query := `update summaries set progress = $2, status = $3, updated_at = $4, title = $5, description = $6, brief_resume = $7, medium_resume = $8, fulltext = $9 where external_id = $1;`

	command, err := conn.Exec(
		ctx,
		query,
		input.ExternalID,
		repository.StatusToString[input.Status].Percentage,
		repository.StatusToString[input.Status].Status,
		now,
		input.Title,
		input.Description,
		input.BriefResume,
		input.MediumResume,
		input.FullText,
	)

	if err != nil {
		return err
	}

	if command.RowsAffected() == 0 {
		return errors.New("register not found")
	}

	return nil
}

func (s *Summary) GetSummaries(ctx context.Context) ([]repository.SummaryOutput, error) {
	conn, err := s.database.Connect(ctx)
	if err != nil {
		return nil, err
	}

	defer s.database.Close(ctx, conn)

	var summaries []repository.SummaryOutput

	query := `
			select 
				s.external_id,
				s.created_at,
				s.updated_at,
				s.status,
				s.title,
				s.description,
				s.brief_resume,
				s.medium_resume,
				s.fulltext,
				s.progress
			from summaries s
			order by s.updated_at desc;
	`

	rows, err := conn.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var summary repository.SummaryOutput

		err = rows.Scan(
			&summary.ExternalID,
			&summary.CreatedAt,
			&summary.UpdatedAt,
			&summary.Status,
			&summary.Title,
			&summary.Description,
			&summary.BriefResume,
			&summary.MediumResume,
			&summary.FullText,
			&summary.Progress,
		)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func (s *Summary) GetSummaryByExternalID(ctx context.Context, externalID uuid.UUID) (*repository.SummaryOutput, error) {
	conn, err := s.database.Connect(ctx)
	if err != nil {
		return nil, err
	}

	defer s.database.Close(ctx, conn)

	var summary repository.SummaryOutput

	query := `
		select
			s.external_id,
				s.created_at,
				s.updated_at,
				s.status,
				s.title,
				s.description,
				s.brief_resume,
				s.medium_resume,
				s.fulltext,
				s.progress
			from summaries s
			where s.external_id = $1;
	`

	row := conn.QueryRow(ctx, query, externalID)

	err = row.Scan(
		&summary.ExternalID,
		&summary.CreatedAt,
		&summary.UpdatedAt,
		&summary.Status,
		&summary.Title,
		&summary.Description,
		&summary.BriefResume,
		&summary.MediumResume,
		&summary.FullText,
		&summary.Progress,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.SummaryNotFound
		}
		return nil, err
	}

	return &summary, nil
}

func (s *Summary) DeleteSummaryByExternalID(ctx context.Context, externalID uuid.UUID) error {
	conn, err := s.database.Connect(ctx)
	if err != nil {
		return err
	}

	defer s.database.Close(ctx, conn)

	query := `
		delete from summaries
		where external_id = $1
	`

	row, err := conn.Exec(ctx, query, externalID)

	if err != nil {
		return err
	}

	if row.RowsAffected() == 0 {
		return application.SummaryNotFound
	}

	return nil
}
