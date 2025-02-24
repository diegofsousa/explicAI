package application

import "errors"

var (
	MissingFile           = errors.New("missing file to upload")
	InvalidFile           = errors.New("invalid file")
	FailedReadFile        = errors.New("fail read to upload")
	SummaryNotFound       = errors.New("summary not found")
	ExternalIDIsInvalid   = errors.New("extenalId is invalid")
	InternalDatabaseError = errors.New("internal database error")
	TranscriptFailed      = errors.New("fail to transcript audio")
	ResumeTextFailed      = errors.New("fail to resume audio")
	UnexpectedErrorList   = errors.New("error on list summaries")
)
