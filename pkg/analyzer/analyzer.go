package analyzer

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/blwsh/llmt/cmd/llmt/config"
	"github.com/blwsh/llmt/lib/file"
	"github.com/blwsh/llmt/lib/logger"
)

type ProjectAnalyzer interface {
	AnalyzeProject(ctx context.Context, projectPath string, destinationPath string, analyzerConfigs []FileAnalyzer, assistant *string) error
}

type ItemAnalyzer interface {
	Analyze(ctx context.Context, prompt string, contents string, threadID *string) (string, error)
}

type FileAnalyzer struct {
	Prompt        string
	Model         string
	ItemAnalyzer  ItemAnalyzer
	Condition     func(filePath string) bool
	ResultHandler func(destFilepath string, result string) error
	Assistant     *config.Assistant
}

type BaseProjectAnalyzer struct {
	logger logger.Logger
}

func (s *BaseProjectAnalyzer) AnalyzeProject(ctx context.Context, projectPath string, destinationPath string, analyzerConfigs []FileAnalyzer, assistantID *string) error {
	dir, err := os.ReadDir(projectPath)
	if err != nil {
		s.logger.Fatalf("failed to read directory \"%s\": %v", projectPath, err)
	}

	var wg sync.WaitGroup

	for _, analyzerConfig := range analyzerConfigs {
		for _, f := range dir {
			wg.Add(1)

			f := f
			go func() {
				defer wg.Done()

				filePath := filepath.Join(projectPath, f.Name())

				if f.IsDir() {
					err = s.AnalyzeProject(ctx, filePath, filepath.Join(destinationPath, f.Name()), analyzerConfigs, assistantID)
					if err != nil {
						s.logger.Errorf("failed to analyze project: %v", err)
					}
				}

				if !analyzerConfig.Condition(filePath) {
					return
				}

				contents, err := file.GetContents(filePath, f)
				if err != nil {
					s.logger.Errorf("failed to get file contents for \"%s\": %v", filePath, err)
					return
				}

				analyzed, err := analyzerConfig.ItemAnalyzer.Analyze(ctx, analyzerConfig.Prompt, contents, assistantID)
				if err != nil {
					s.logger.Error(err)
					return
				}

				if analyzerConfig.ResultHandler != nil {
					err = analyzerConfig.ResultHandler(filepath.Join(destinationPath, f.Name()), analyzed)
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
