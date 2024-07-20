package main

import (
	"context"
	_ "embed"
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

func main() {
	ctx := context.Background()

	if os.Getenv(EnvOpenAIToken) == "" {
		l.Fatal(EnvOpenAIToken + " environment variable not set")
	}

	var projectDir = os.Args[1]
	if projectDir == "" {
		l.Fatal("usage: analyze <project_dir>")
	}

	err := project_analyzer.New().AnalyzeProject(ctx, projectDir, "./out", []project_analyzer.FileAnalyzer{
		{
			Prompt: prompt,
			//Analyzer: ollama.New("http://localhost:11434", "overview"),
			Analyzer:  openai.New(os.Getenv(EnvOpenAIToken), "gpt-4o-mini"),
			Condition: myFancyConditionFunc,
			ResultHandler: func(destFilePath string, result string) error {
				outputPath := filepath.Join(destFilePath + ".md")

				err := file.WriteTo(outputPath, result)
				if err != nil {
					l.Error(err)
				}

				l.Info("analyzed file: ", filepath.Base(destFilePath), " -> ", outputPath)

				return nil
			},
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
