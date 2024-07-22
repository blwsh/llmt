package main

import (
	"context"
	"os"

	"github.com/blwsh/llmt/pkg/project_analyzer"
)

func analyze(ctx context.Context, source, target string, c config) error {
	var (
		writer   = markdownWriter{logger: log}
		resolver = analyzerResolver{
			OpenAITokenResolver: func() string {
				var openAIToken string
				if openAIToken = os.Getenv(envOpenAIToken); os.Getenv(envOpenAIToken) == "" {
					log.Fatalf("env var %s must be set when using OpenAI analyzers", envOpenAIToken)
				}

				return openAIToken
			},
			OllamaHostResolver: func() string {
				var ollamaHost string
				if ollamaHost = os.Getenv(envOllamaHost); os.Getenv(envOllamaHost) == "" {
					ollamaHost = "http://localhost:11434"
				}

				return ollamaHost
			},
		}
		compiledRegexes = compileRegexesFromConfig(c)
		analyzers       []project_analyzer.FileAnalyzer
	)

	for _, a := range c.Analyzers {
		analyzer, err := resolver.resolve(a.Analyzer, a.Model)
		if err != nil {
			log.Warnf("failed to resolve analyzer %s: %v", a.Analyzer, err)
		}

		analyzers = append(analyzers, project_analyzer.FileAnalyzer{
			Prompt:        a.Prompt,
			Analyzer:      analyzer,
			Condition:     condition(a, compiledRegexes),
			ResultHandler: writer.write,
		})
	}

	err := project_analyzer.New(project_analyzer.WithLogger(log)).AnalyzeProject(ctx, source, target, analyzers)
	if err != nil {
		return err
	}

	return nil
}
