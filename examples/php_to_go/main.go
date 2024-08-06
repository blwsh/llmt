package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/blwsh/llmt/lib/file"
	"github.com/blwsh/llmt/lib/logger"
	"github.com/blwsh/llmt/pkg/analyzer"
	chat2 "github.com/blwsh/llmt/pkg/analyzer/item_analyzer/openai/chat"
	"github.com/blwsh/llmt/pkg/analyzer/project_analyzer/assistant"
)

const (
	EnvOpenAIToken = "OPENAI_TOKEN"
)

var l = logger.New(false)

//go:embed prompt.txt
var prompt string

var (
	cwd, _ = os.Getwd()
	here   = filepath.Join(cwd, "examples/php_to_go")
)

func main() {
	ctx := context.Background()

	var openAIToken string
	if openAIToken = os.Getenv(EnvOpenAIToken); os.Getenv(EnvOpenAIToken) == "" {
		l.Fatal(EnvOpenAIToken + " environment variable not set")
	}

	err := assistant.New(openAIToken, assistant.WithLogger(l)).
		AnalyzeProject(ctx, cwd+"/examples/examplePhpProject", here+"/exampleGoProject", []analyzer.FileAnalyzer{
			{
				Prompt:       prompt,                                // you may want to just use empty string if your model has a system prompt already
				ItemAnalyzer: chat2.New(openAIToken, "gpt-4o-mini"), // you can also use ollama: ollama.New("http://localhost:11434", "php_to_go"),
				Condition: func(filePath string) bool {
					return strings.HasSuffix(filePath, ".php") && !strings.Contains(filePath, "test") && !strings.Contains(filePath, "vendor")
				},
				ResultHandler: func(destFilepath string, result string) error {
					if _, err := os.Stat(filepath.Dir(destFilepath)); os.IsNotExist(err) {
						err = os.MkdirAll(filepath.Dir(destFilepath), os.ModePerm)
						if err != nil {
							return fmt.Errorf("failed to create directory: %w", err)
						}
					}

					outputPath := filepath.Join(destFilepath + ".go")

					err := file.WriteTo(outputPath, result)
					if err != nil {
						l.Error(err)
					}

					l.Info("wrote to " + outputPath)

					return nil
				},
			},
		}, nil)
	if err != nil {
		l.Fatal(err)
	}
}
