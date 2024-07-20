package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cdb/lib/file"
	"cdb/lib/logger"
	"cdb/pkg/analyzer/openai"
	"cdb/pkg/project_analyzer"
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

	err := project_analyzer.New(project_analyzer.WithLogger(l)).
		AnalyzeProject(ctx, here+"/myPhpProject", here+"/docs", []project_analyzer.FileAnalyzer{
			{
				Prompt:        prompt,
				Analyzer:      openai.New(openAIToken, "gpt-4o-mini"), // or: ollama.New("http://localhost:11434", "overview"),
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
		!strings.Contains(filePath, "tests") &&
		!strings.Contains(filePath, "test.php") &&
		// exclude composer files and directories
		!strings.Contains(filePath, "vendor") &&
		!strings.Contains(filePath, "autoload") &&
		!strings.Contains(filePath, "bootstrap") &&
		// exclude node_modules
		!strings.Contains(filePath, "node_modules")
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
