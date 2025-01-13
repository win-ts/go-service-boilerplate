package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const (
	outgoing = "OUTGOING"
)

// CustomTransport is a custom transport for the HTTP client with log trace
type CustomTransport struct {
	transport       http.RoundTripper
	disableLogTrace bool
}

// RoundTrip executes a single HTTP transaction
func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.disableLogTrace {
		return t.transport.RoundTrip(req)
	}

	var requestBody string
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		requestBody = string(bodyBytes)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	slog.Info(fmt.Sprintf("HTTP_CLIENT_REQUEST | %s | %s | %s", req.Method, outgoing, req.URL.Path),
		slog.String("method", req.Method),
		slog.String("path", req.URL.Path),
		slog.String("query", req.URL.RawQuery),
		slog.Any("reqHeader", req.Header),
		slog.String("reqBody", requestBody),
		slog.String("scope", outgoing),
	)

	var responseBody string

	now := time.Now()
	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	latency := time.Since(now).Milliseconds()

	if resp.Body != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}
		responseBody = string(bodyBytes)
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	logMsg := fmt.Sprintf("HTTP_CLIENT_RESPONSE | %s | %d | %d ms | %s | %s", req.Method, resp.StatusCode, latency, outgoing, req.URL.Path)
	logFields := []slog.Attr{
		slog.Int("status", resp.StatusCode),
		slog.String("method", req.Method),
		slog.String("path", req.URL.Path),
		slog.String("query", req.URL.RawQuery),
		slog.Int64("latency", latency),
		slog.Any("respHeaders", resp.Header),
		slog.String("respBody", responseBody),
		slog.String("scope", outgoing),
	}

	switch {
	case resp.StatusCode >= 500:
		slog.Error(logMsg, convertToAny(logFields)...)
	case resp.StatusCode >= 400:
		slog.Warn(logMsg, convertToAny(logFields)...)
	default:
		slog.Info(logMsg, convertToAny(logFields)...)
	}

	return resp, nil
}

func convertToAny(attrs []slog.Attr) []any {
	anyAttrs := make([]any, len(attrs))
	for i, attr := range attrs {
		anyAttrs[i] = attr
	}
	return anyAttrs
}
