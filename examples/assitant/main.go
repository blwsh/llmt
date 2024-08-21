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
	"github.com/blwsh/llmt/pkg/analyzer/project_analyzer/chat"
)

const (
	EnvOpenAIToken = "OPENAI_TOKEN"
)

var l = logger.New(false)

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
		AnalyzeProject(ctx, cwd+"/examples/examplePhpProject", here+"/exampleGoProject", []analyzer.FileAnalyzer{
			{
				Prompt:        prompt,
				ItemAnalyzer:  chat2.New(openAIToken, "gpt-3.5-turbo"),
				Condition:     myFancyConditionFunc,
				ResultHandler: myDocsWriterFunc,
			},
		}, nil)
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

	outputPath := filepath.Join(destFilepath + ".go")

	err := file.WriteTo(outputPath, result)
	if err != nil {
		l.Error(err)
	}

	l.Info("converted file: ", filepath.Base(destFilepath), " -> ", outputPath)

	return nil
}
