package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WebLogger func(msg string, args ...interface{})

type DoRequestInput struct {
	Input       any
	Output      any
	URL         string
	Method      string
	RequestFunc func(*http.Request) error
	Headers     map[string][]string
	Client      *http.Client
	Logger      WebLogger
}

func DoRequest(ctx context.Context, d DoRequestInput) error {
	if d.Client == nil {
		d.Client = http.DefaultClient
	}
	var rdr io.Reader
	if d.Input != nil {
		d, err := json.Marshal(d.Input)
		if err != nil {
			return fmt.Errorf("marshal input: %w", err)
		}
		rdr = bytes.NewReader(d)
	}
	if d.Method == "" {
		d.Method = http.MethodGet
	}
	req, err := http.NewRequest(d.Method, d.URL, rdr)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)
	if d.Input != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if d.RequestFunc != nil {
		if err := d.RequestFunc(req); err != nil {
			return fmt.Errorf("request func: %w", err)
		}
		if d.Logger != nil {
			d.Logger("with request func headers %v", req.Header)
		}
	}
	if d.Headers != nil {
		for k, v := range d.Headers {
			req.Header[k] = v
		}
		if d.Logger != nil {
			d.Logger("with headers %v", req.Header)
		}
	}
	resp, err := d.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrorFromResponse(ctx, resp)
	}
	if d.Output != nil {
		if err := json.NewDecoder(resp.Body).Decode(d.Output); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func ErrorFromResponse(ctx context.Context, resp *http.Response) error {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return fmt.Errorf("copy response body: %w", err)
	}
	jerr := &StatusCodeError{}
	if err := json.Unmarshal(buf.Bytes(), jerr); err != nil {
		// not a valid json response!
		return &StatusCodeError{
			StatusCode: resp.StatusCode,
			Raw:        buf.String(),
		}
	}
	jerr.StatusCode = resp.StatusCode
	return jerr
}

type StatusCodeError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code,omitempty"`
	Raw        string `json:"raw,omitempty"`
}

func (e *StatusCodeError) Error() string {
	if e.StatusCode > 0 && e.Message != "" {
		return fmt.Sprintf("status: %d message: %s", e.StatusCode, e.Message)
	}
	if e.Message != "" {
		return e.Message
	}
	if e.StatusCode != 0 && e.Raw != "" {
		return fmt.Sprintf("status: %d raw: %s", e.StatusCode, e.Raw)
	}
	if e.Raw != "" {
		return "raw response: " + e.Raw
	}
	return "unknown error"
}
