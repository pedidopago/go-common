package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type ExtractContentTypeFunc func(contentType string)

type ContextVar string

const (
	ContextVarExtractContentType ContextVar = "extract_content_type"
)

// DownloadInto downloads a file from a public URL
func DownloadInto(ctx context.Context, url string, w io.Writer) (int64, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()
	if vi := ctx.Value(ContextVarExtractContentType); vi != nil {
		if vf, ok := vi.(ExtractContentTypeFunc); ok {
			vf(resp.Header.Get("Content-Type"))
		}
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to download: %s", resp.Status)
	}
	sz, err := io.Copy(w, resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to download: %w", err)
	}
	return sz, nil
}
