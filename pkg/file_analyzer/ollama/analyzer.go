package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"cdb/pkg/file_analyzer"
)

func New(llamaHost, model string) file_analyzer.Analyzer {
	return &llama3Analyzer{llamaHost: llamaHost, model: model}
}

type llama3Analyzer struct {
	llamaHost string
	model     string
}

type llamaReq struct {
	Stream bool
	Model  string
	Prompt string
}

type llamaRes struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	DoneReason         string    `json:"done_reason"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int       `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int       `json:"eval_duration"`
	Error              *string   `json:"error"`
}

func (l llama3Analyzer) Analyze(ctx context.Context, _ string, contents string) (string, error) {
	req, err := json.Marshal(llamaReq{
		Stream: false,
		Model:  l.model,
		Prompt: contents,
	})
	if err != nil {
		return "", err
	}

	generateUrl, err := url.Parse(l.llamaHost)
	if err != nil {
		return "", err
	}

	generateUrl.Path = "/api/generate"

	request, err := http.NewRequestWithContext(ctx, "POST", generateUrl.String(), bytes.NewReader(req))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res llamaRes
	if err := json.Unmarshal(responseData, &res); err != nil {
		return "", err
	}

	if res.Error != nil {
		return "", errors.New(*res.Error)
	}

	return res.Response, nil
}
