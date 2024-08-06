package analyze

import (
	"context"
	"os"

	"github.com/blwsh/llmt/cmd/llmt/config"
	"github.com/blwsh/llmt/lib/logger"
	"github.com/blwsh/llmt/pkg/analyzer"
	"github.com/blwsh/llmt/pkg/analyzer/project_analyzer/assistant"
)

const (
	envOpenAIToken = "OPENAI_TOKEN"
	envOllamaHost  = "OLLAMA_HOST"
)

var log = logger.NewCMDLogger()

func Analyze(ctx context.Context, source, target string, c config.AnalyzeConfig) error {
	var (
		writer   = markdownWriter{logger: log}
		resolver = itemAnalyzerResolver{
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
		condition = checker{compiledRegexes: compileRegexesFromConfig(c)}
		analyzers []analyzer.FileAnalyzer
	)

	for _, a := range c.Analyzers {
		itemAnalyzer, err := resolver.resolve(a.Analyzer, a.Model)
		if err != nil {
			log.Warnf("failed to resolve analyzer %s: %v", a.Analyzer, err)
		}

		condition.a = a

		analyzers = append(analyzers, analyzer.FileAnalyzer{
			Prompt:        a.Prompt,
			Model:         a.Model,
			Assistant:     a.Assistant,
			ItemAnalyzer:  itemAnalyzer,
			Condition:     condition.check,
			ResultHandler: writer.write,
		})
	}

	err := assistant.New(resolver.OpenAITokenResolver(), assistant.WithLogger(log)).
		AnalyzeProject(ctx, source, target, analyzers, nil)
	if err != nil {
		return err
	}

	return nil
}
