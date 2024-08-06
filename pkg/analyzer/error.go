package analyzer

import "errors"

var (
	ErrUnexpectedResponse = errors.New("unexpected response")
)

type RateLimitError struct {
	Err error
}

func (r *RateLimitError) Error() string {
	return r.Err.Error()
}
