package amiami

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrItemNotFound = errors.New("item not found")

type StatusError struct {
	StatusCode int
	Status     string
	GCode      string
	Body       string
	RetryAfter time.Duration
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("AmiAmi %s: %s - %s", e.GCode, e.Status, e.Body)
}

func (e *StatusError) Unwrap() error {
	if e.StatusCode == http.StatusBadRequest {
		return ErrItemNotFound
	}
	return nil
}

func (e *StatusError) Retryable() bool {
	return e.StatusCode == 429 || e.StatusCode >= 500
}
