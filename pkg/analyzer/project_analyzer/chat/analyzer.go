package chat

import (
	"context"

	"github.com/blwsh/llmt/lib/logger"
	"github.com/blwsh/llmt/pkg/analyzer"
)

type Analyzer struct {
	analyzer.BaseProjectAnalyzer
	logger logger.Logger
}

type AnalyzerOption func(*Analyzer)

func New(options ...AnalyzerOption) analyzer.ProjectAnalyzer {
	p := &Analyzer{
		logger: logger.New(false),
	}

	for _, option := range options {
		option(p)
	}

	return p
}

func (a *Analyzer) AnalyzeProject(ctx context.Context, projectPath string, destinationPath string, analyzerConfigs []analyzer.FileAnalyzer, _ *string) error {
	err := a.BaseProjectAnalyzer.AnalyzeProject(ctx, projectPath, destinationPath, analyzerConfigs, nil)
	if err != nil {
		return err
	}

	return nil
}

func WithLogger(l logger.Logger) AnalyzerOption {
	return func(p *Analyzer) {
		p.logger = l
	}
}
