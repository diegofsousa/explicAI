package summarize

import "context"

type Summarize interface {
	Resume(ctx context.Context, transcription string) (*ResumeOutput, error)
	FullTextOrganize(ctx context.Context, transcription string) (*string, error)
}
