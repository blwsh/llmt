package file_analyzer

import (
	"context"
	"errors"
)

type Analyzer interface {
	Analyze(ctx context.Context, prompt string, contents string) (string, error)
}

var (
	ErrUnexpectedResponse = errors.New("unexpected response")
)

type RateLimitError struct {
	Err error
}

func (r *RateLimitError) Error() string {
	return r.Err.Error()
}
