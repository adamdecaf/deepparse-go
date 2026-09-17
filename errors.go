package deepparsego

import (
	"encoding/json"
	"fmt"
	"strings"
)

// APIError is returned when the deepparse HTTP API responds with a non-200 status.
type APIError struct {
	StatusCode int
	Detail     string
}

func (e *APIError) Error() string {
	if e.Detail == "" {
		return fmt.Sprintf("deepparse API error: status %d", e.StatusCode)
	}
	return fmt.Sprintf("deepparse API error: status %d: %s", e.StatusCode, e.Detail)
}

func newAPIError(status int, body []byte) *APIError {
	err := &APIError{StatusCode: status, Detail: strings.TrimSpace(string(body))}

	var payload struct {
		Detail json.RawMessage `json:"detail"`
	}
	if json.Unmarshal(body, &payload) != nil || len(payload.Detail) == 0 {
		return err
	}

	var asString string
	if json.Unmarshal(payload.Detail, &asString) == nil {
		err.Detail = asString
		return err
	}

	err.Detail = string(payload.Detail)
	return err
}
