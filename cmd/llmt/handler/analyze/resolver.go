package analyze

import (
	"errors"

	"github.com/blwsh/llmt/pkg/file_analyzer"
	"github.com/blwsh/llmt/pkg/file_analyzer/ollama"
	"github.com/blwsh/llmt/pkg/file_analyzer/openai"
)

var ErrUnknownAnalyzer = errors.New("unknown analyzer")

type analyzerResolver struct {
	OpenAITokenResolver func() string
	OllamaHostResolver  func() string
}

func (a analyzerResolver) resolve(analyzer, model string) (file_analyzer.Analyzer, error) {
	switch analyzer {
	case "openai":
		return openai.New(a.OpenAITokenResolver(), model), nil
	case "ollama":
		return ollama.New(a.OllamaHostResolver(), model), nil
	default:
		return nil, ErrUnknownAnalyzer
	}
}
