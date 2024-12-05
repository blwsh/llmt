package analyze

import (
	"errors"

	"github.com/blwsh/llmt/pkg/analyzer"
	"github.com/blwsh/llmt/pkg/analyzer/item_analyzer/ollama"
	openai "github.com/blwsh/llmt/pkg/analyzer/item_analyzer/openai/assitant"
)

var ErrUnknownAnalyzer = errors.New("unknown analyzer")

type itemAnalyzerResolver struct {
	OpenAITokenResolver func() string
	OllamaHostResolver  func() string
}

func (a itemAnalyzerResolver) resolve(analyzer, model string) (analyzer.ItemAnalyzer, error) {
	switch analyzer {
	case "openai":
		return openai.New(a.OpenAITokenResolver(), model), nil
	case "ollama":
		return ollama.New(a.OllamaHostResolver(), model), nil
	default:
		return nil, ErrUnknownAnalyzer
	}
}
