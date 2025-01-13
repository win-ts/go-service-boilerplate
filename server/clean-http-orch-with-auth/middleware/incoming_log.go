package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	healthCheckPath = "/health"
	incoming        = "INCOMING"
)

// IncomingLogTrace logs the incoming request and response
func IncomingLogTrace() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == healthCheckPath {
				return next(c)
			}

			now := time.Now()

			// request
			req := c.Request()
			reqBody, err := io.ReadAll(req.Body)
			if err != nil {
				slog.Error("failed to read incoming request body", slog.Any("error", err))
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
			}
			req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			slog.Info(fmt.Sprintf("REQUEST_INFO | %s | %s | %s", req.Method, incoming, req.URL.Path),
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.String("query", req.URL.RawQuery),
				slog.String("userAgent", req.UserAgent()),
				slog.String("referer", req.Referer()),
				slog.String("protocol", req.Proto),
				slog.Any("reqHeader", req.Header),
				slog.String("reqBody", string(reqBody)),
				slog.String("scope", incoming),
			)

			// response
			res := c.Response()
			bw := newBodyWriter(res.Writer)
			res.Writer = bw

			err = next(c)
			if err != nil {
				return err
			}

			latency := time.Since(now).Milliseconds()
			res.Writer = bw

			logMsg := fmt.Sprintf("RESPONSE_INFO | %s | %d | %d ms | %s | %s", req.Method, res.Status, latency, incoming, req.URL.Path)
			logFields := []slog.Attr{
				slog.String("respBody", bw.body.String()),
				slog.Int("status", res.Status),
				slog.Int64("latency", latency),
				slog.String("scope", incoming),
			}

			switch {
			case res.Status >= 500:
				slog.Error(logMsg, convertToAny(logFields)...)
			case res.Status >= 400:
				slog.Warn(logMsg, convertToAny(logFields)...)
			default:
				slog.Info(logMsg, convertToAny(logFields)...)
			}

			return nil
		}
	}
}

type bodyWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func newBodyWriter(writer http.ResponseWriter) *bodyWriter {
	body := bytes.NewBufferString("")

	return &bodyWriter{
		ResponseWriter: writer,
		body:           body,
	}
}

// implements http.ResponseWriter
func (w *bodyWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}

	return w.ResponseWriter.Write(b)
}

func convertToAny(attrs []slog.Attr) []any {
	anyAttrs := make([]any, len(attrs))
	for i, attr := range attrs {
		anyAttrs[i] = attr
	}
	return anyAttrs
}
