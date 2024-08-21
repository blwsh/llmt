package openai

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

func (f *openai) Analyze(ctx context.Context, prompt string, contents string, assistantID *string) (string, error) {
	err := f.rpmLimiter.Wait(ctx)
	if err != nil {
		return "", err
	}

	resp, err := f.gpt.CreateThreadAndRun(ctx, oai.CreateThreadAndRunRequest{
		RunRequest: oai.RunRequest{
			AssistantID:  *assistantID,
			Model:        f.model,
			Instructions: prompt,
		},
		Thread: oai.ThreadRequest{
			Messages: []oai.ThreadMessage{
				{Role: oai.ChatMessageRoleUser, Content: contents},
			},
		},
	})
	if err != nil {
		var apiErr *oai.APIError
		if errors.As(err, &apiErr) {
			if apiErr.HTTPStatusCode == 429 {
				return "", &analyzer.RateLimitError{Err: apiErr}
			}
		}

		return "", err
	}

	message, err := f.gpt.ListMessage(ctx, resp.ThreadID, nil, nil, nil, nil)
	if err != nil {
		return "", err
	}

	if len(message.Messages) != 1 {
		return "", fmt.Errorf("%w: expected 1 choice, got %d", analyzer.ErrUnexpectedResponse, len(message.Messages))
	}

	return message.Messages[0].Content[0].Text.Value, nil
}
