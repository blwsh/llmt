package assistant

import (
	"context"

	"github.com/blwsh/llmt/lib/logger"
	"github.com/blwsh/llmt/pkg/analyzer"
)

type projectAnalyzer struct {
	logger logger.Logger
}

type AnalyzerOption func(*projectAnalyzer)

func New(options ...AnalyzerOption) analyzer.ProjectAnalyzer {
	p := &projectAnalyzer{
		logger: logger.New(false),
	}

	for _, option := range options {
		option(p)
	}

	return p
}

func WithLogger(l logger.Logger) AnalyzerOption {
	return func(p *projectAnalyzer) {
		p.logger = l
	}
}

func (s *projectAnalyzer) AnalyzeProject(ctx context.Context, projectPath string, destinationPath string, analyzers []analyzer.FileAnalyzerConfig) error {
	return nil
}
