package analyzer

import (
	"context"
	"errors"
)

type FileAnalyzerConfig struct {
	Prompt        string
	Analyzer      ItemAnalyzer
	Condition     func(filePath string) bool
	ResultHandler func(destFilepath string, result string) error
}

type ProjectAnalyzer interface {
	AnalyzeProject(ctx context.Context, projectPath string, destinationPath string, analyzers []FileAnalyzerConfig) error
}

type ItemAnalyzer interface {
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
