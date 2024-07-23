package chat

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/blwsh/llmt/lib/file"
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
	dir, err := os.ReadDir(projectPath)
	if err != nil {
		s.logger.Fatalf("failed to read directory \"%s\": %v", projectPath, err)
	}

	var wg sync.WaitGroup

	for _, projectAnalyzer := range analyzers {
		for _, f := range dir {
			wg.Add(1)

			go func() {
				defer wg.Done()

				filePath := filepath.Join(projectPath, f.Name())

				if f.IsDir() {
					err = s.AnalyzeProject(ctx, filePath, filepath.Join(destinationPath, f.Name()), analyzers)
					if err != nil {
						s.logger.Errorf("failed to analyze project: %v", err)
					}
				}

				if !projectAnalyzer.Condition(filePath) {
					return
				}

				contents, err := file.GetContents(filePath, f)
				if err != nil {
					s.logger.Errorf("failed to get file contents for \"%s\": %v", filePath, err)
					return
				}

				analyzed, err := projectAnalyzer.Analyzer.Analyze(ctx, projectAnalyzer.Prompt, contents)
				if err != nil {
					s.logger.Error(err)
					return
				}

				if projectAnalyzer.ResultHandler != nil {
					err = projectAnalyzer.ResultHandler(filepath.Join(destinationPath, f.Name()), analyzed)
					if err != nil {
						s.logger.Errorf("failed to handle result for \"%s\": %v", filePath, err)
					}
				}
			}()
		}
	}

	wg.Wait()

	return nil
}
