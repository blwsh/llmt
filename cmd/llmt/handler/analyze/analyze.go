package analyze

import (
	"context"
	"os"

	"github.com/blwsh/llmt/config"
	"github.com/blwsh/llmt/lib/logger"
	"github.com/blwsh/llmt/pkg/analyzer"
	"github.com/blwsh/llmt/pkg/analyzer/project_analyzer/chat"
)

const (
	envOpenAIToken = "OPENAI_TOKEN"
	envOllamaHost  = "OLLAMA_HOST"
)

var log = logger.NewCMDLogger()

func Analyze(ctx context.Context, source, target string, c config.Config) error {
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
		analyzers       []analyzer.FileAnalyzerConfig
	)

	for _, a := range c.Analyzers {
		an, err := resolver.resolve(a.Analyzer, a.Model)
		if err != nil {
			log.Warnf("failed to resolve analyzer %s: %v", a.Analyzer, err)
		}

		analyzers = append(analyzers, analyzer.FileAnalyzerConfig{
			Prompt:        a.Prompt,
			Analyzer:      an,
			Condition:     condition(a, compiledRegexes),
			ResultHandler: writer.write,
		})
	}

	err := chat.New(chat.WithLogger(log)).AnalyzeProject(ctx, source, target, analyzers)
	if err != nil {
		return err
	}

	return nil
}
