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
	"github.com/blwsh/llmt/pkg/analyzer/item_analyzer/openai"
	"github.com/blwsh/llmt/pkg/analyzer/project_analyzer/chat"
)

const (
	EnvOpenAIToken = "OPENAI_TOKEN"
)

var l = logger.New(false)

//go:embed prompt.txt
var prompt string

var (
	cwd, _ = os.Getwd()
	here   = filepath.Join(cwd, "examples/overview")
)

func main() {
	ctx := context.Background()

	var openAIToken string
	if openAIToken = os.Getenv(EnvOpenAIToken); os.Getenv(EnvOpenAIToken) == "" {
		l.Fatal(EnvOpenAIToken + " environment variable not set")
	}

	err := chat.New(chat.WithLogger(l)).
		AnalyzeProject(ctx, cwd+"/examples/examplePhpProject", here+"/docs", []analyzer.FileAnalyzerConfig{
			{
				Prompt:        prompt,                                 // you may want to just use empty string if your model has a system prompt already
				Analyzer:      openai.New(openAIToken, "gpt-4o-mini"), // you can also use ollama: ollama.New("http://localhost:11434", "overview"),
				Condition:     myFancyConditionFunc,
				ResultHandler: myDocsWriterFunc,
			},
		})
	if err != nil {
		l.Fatal(err)
	}
}

func myFancyConditionFunc(filePath string) bool {
	return strings.HasSuffix(filePath, ".php") &&
		// exclude tests
		!strings.Contains(filePath, "test") &&
		// exclude composer files and directories
		!strings.Contains(filePath, "vendor")
}

func myDocsWriterFunc(destFilepath string, result string) error {
	if _, err := os.Stat(filepath.Dir(destFilepath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(destFilepath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	outputPath := filepath.Join(destFilepath + ".md")

	err := file.WriteTo(outputPath, result)
	if err != nil {
		l.Error(err)
	}

	l.Info("analyzed file: ", filepath.Base(destFilepath), " -> ", outputPath)

	return nil
}
