package assistant

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	oai "github.com/sashabaranov/go-openai"

	"github.com/blwsh/llmt/lib/logger"
	"github.com/blwsh/llmt/pkg/analyzer"
)

type Analyzer struct {
	analyzer.BaseProjectAnalyzer
	gpt    *oai.Client
	logger logger.Logger
}

type AnalyzerOption func(*Analyzer)

func New(authToken string, options ...AnalyzerOption) analyzer.ProjectAnalyzer {
	p := &Analyzer{
		gpt:    oai.NewClient(authToken),
		logger: logger.New(false),
	}

	for _, option := range options {
		option(p)
	}

	return p
}

var homeDir, _ = os.UserHomeDir()

func (s *Analyzer) AnalyzeProject(ctx context.Context, projectPath string, destinationPath string, analyzerConfigs []analyzer.FileAnalyzer, _ *string) error {
	for _, analyzerConfig := range analyzerConfigs {
		if analyzerConfig.Assistant != nil {
			s.logger.Infof("creating assistant %s", *analyzerConfig.Assistant.Name)

			newAssistant, err := s.gpt.CreateAssistant(ctx, oai.AssistantRequest{
				Model:        analyzerConfig.Model,
				Name:         analyzerConfig.Assistant.Name,
				Description:  analyzerConfig.Assistant.Description,
				Instructions: &analyzerConfig.Prompt,
				Tools: []oai.AssistantTool{
					{Type: oai.AssistantToolTypeFileSearch},
				}})
			if err != nil {
				s.logger.Errorf("failed to create assistant %v: %v", analyzerConfig.Assistant.Name, err)
			}

			for _, file := range *analyzerConfig.Assistant.Files {
				createFile, err := s.gpt.CreateFile(ctx, oai.FileRequest{
					FileName: filepath.Base(file),
					FilePath: strings.Replace(file, "~", homeDir, 1),
					Purpose:  "assistants",
				})
				if err != nil {
					return err
				}

				assistantFile, err := s.gpt.CreateAssistantFile(ctx, newAssistant.ID, oai.AssistantFileRequest{
					FileID: createFile.ID,
				})
				if err != nil {
					s.logger.Errorf("failed to create assistant file %v: %v", file, err)
				}

				s.logger.Infof("created assistant file %v", assistantFile.ID)
			}

			defer func() {
				s.logger.Infof("deleting assistant %v", newAssistant.ID)

				_, err := s.gpt.DeleteAssistant(ctx, newAssistant.ID)
				if err != nil {
					s.logger.Errorf("failed to delete assistant %v: %v", newAssistant.ID, err)
				}
			}()

			err = s.BaseProjectAnalyzer.AnalyzeProject(ctx, projectPath, destinationPath, analyzerConfigs, &newAssistant.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func WithLogger(l logger.Logger) AnalyzerOption {
	return func(p *Analyzer) {
		p.logger = l
	}
}
