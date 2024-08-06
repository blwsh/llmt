package chat

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	oai "github.com/sashabaranov/go-openai"
	"golang.org/x/time/rate"

	"github.com/blwsh/llmt/pkg/analyzer"
)

func New(llmAuthToken, model string) analyzer.ItemAnalyzer {
	return &openai{
		gpt:        oai.NewClient(llmAuthToken),
		model:      model,
		rpmLimiter: rate.NewLimiter(rate.Every(time.Minute), 500),
		tpmLimiter: rate.NewLimiter(rate.Every(time.Minute), 200_000),
	}
}

type openai struct {
	gpt        *oai.Client
	model      string
	rpmLimiter *rate.Limiter
	tpmLimiter *rate.Limiter
}

func (f *openai) Analyze(ctx context.Context, prompt string, contents string, _ *string) (string, error) {
	err := f.rpmLimiter.Wait(ctx)
	if err != nil {
		return "", err
	}

	resp, err := f.gpt.CreateChatCompletion(
		context.Background(),
		oai.ChatCompletionRequest{
			Model: f.model,
			Messages: []oai.ChatCompletionMessage{
				{Role: oai.ChatMessageRoleSystem, Content: prompt},
				{Role: oai.ChatMessageRoleUser, Content: contents},
			},
		},
	)
	if err != nil {
		var apiErr *oai.APIError
		if errors.As(err, &apiErr) {
			if apiErr.HTTPStatusCode == 429 {
				return "", &analyzer.RateLimitError{Err: apiErr}
			}
		}

		return "", err
	}

	if len(resp.Choices) != 1 {
		return "", fmt.Errorf("%w: expected 1 choice, got %d", analyzer.ErrUnexpectedResponse, len(resp.Choices))
	}

	return resp.Choices[0].Message.Content, nil
}
