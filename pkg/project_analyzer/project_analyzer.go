package project_analyzer

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"cdb/lib/file"
	"cdb/lib/logger"
	"cdb/pkg/analyzer"
)

type FileAnalyzer struct {
	Prompt    string
	Analyzer  analyzer.Analyzer
	Condition func(filePath string) bool
}

type ProjectAnalyzer interface {
	AnalyzeProject(ctx context.Context, projectPath string, destinationPath string, analyzers []FileAnalyzer) error
}

type projectAnalyzer struct {
	logger logger.Logger
}

func New() ProjectAnalyzer {
	return &projectAnalyzer{
		logger: logger.New(false),
	}
}

func (s *projectAnalyzer) AnalyzeProject(ctx context.Context, projectPath string, destinationPath string, analyzers []FileAnalyzer) error {
	dir, err := os.ReadDir(projectPath)
	if err != nil {
		s.logger.Fatal(err)
	}

	if _, err := os.Stat(destinationPath); os.IsNotExist(err) {
		err = os.Mkdir(destinationPath, os.ModePerm)
		if err != nil {
			return err
		}
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
						s.logger.Error(err)
					}
				}

				if !projectAnalyzer.Condition(filePath) {
					return
				}

				contents, err := file.GetContents(filePath, f)
				if err != nil {
					s.logger.Error(err)
					return
				}

				analyzed, err := projectAnalyzer.Analyzer.Analyze(ctx, projectAnalyzer.Prompt, contents)
				if err != nil {
					s.logger.Error(err)
					return
				}

				outputPath := filepath.Join(destinationPath, f.Name()+".md")

				err = file.WriteTo(outputPath, analyzed)
				if err != nil {
					s.logger.Error(err)
					return
				}

				s.logger.Info("analyzed file: ", f)
			}()
		}
	}

	wg.Wait()

	destDir, err := os.ReadDir(destinationPath)

	for _, entry := range destDir {
		if !entry.IsDir() {
			return nil
		}
	}

	_ = os.Remove(destinationPath)

	return nil
}
