package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/time/rate"
)

// HTTPClientOptions represents an option for the HTTP client
type HTTPClientOptions struct {
	MaxConns                 int
	MaxRetry                 int
	Timeout                  time.Duration
	InsecureSkipVerify       bool
	MaxTransactionsPerSecond int
}

// NewHTTPClient creates a new HTTP client with the given options
func NewHTTPClient(opts *HTTPClientOptions) *http.Client {
	httpClient := retryablehttp.NewClient()

	httpClient.RetryMax = opts.MaxRetry
	httpClient.RetryWaitMin = 1 * time.Second
	httpClient.RetryWaitMax = 5 * time.Second
	httpClient.HTTPClient.Timeout = opts.Timeout

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = opts.MaxConns
	t.MaxConnsPerHost = opts.MaxConns
	t.MaxIdleConnsPerHost = opts.MaxConns
	if t.TLSClientConfig == nil {
		t.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	t.TLSClientConfig.InsecureSkipVerify = opts.InsecureSkipVerify
	httpClient.HTTPClient.Transport = t

	if opts.MaxTransactionsPerSecond != 0 {
		limiter := rate.NewLimiter(rate.Limit(opts.MaxTransactionsPerSecond), 1)
		setMaxTPS := func(req *http.Request, via []*http.Request) error {
			ctx := req.Context()
			err := limiter.Wait(ctx)
			if err != nil {
				return err
			}
			return nil
		}
		httpClient.HTTPClient.CheckRedirect = setMaxTPS
	}

	return httpClient.StandardClient()
}

func Get[RES any](ctx context.Context, client *http.Client, url string, header map[string]string) (*Response[RES], error) {
	return Do[bytes.Buffer, RES](ctx, client, http.MethodGet, url, header, bytes.Buffer{})
}

func Post[REQ, RES any](ctx context.Context, client *http.Client, url string, header map[string]string, payload REQ) (*Response[RES], error) {
	return Do[REQ, RES](ctx, client, http.MethodPost, url, header, payload)
}

func Put[REQ, RES any](ctx context.Context, client *http.Client, url string, header map[string]string, payload REQ) (*Response[RES], error) {
	return Do[REQ, RES](ctx, client, http.MethodPut, url, header, payload)
}

func Patch[REQ, RES any](ctx context.Context, client *http.Client, url string, header map[string]string, payload REQ) (*Response[RES], error) {
	return Do[REQ, RES](ctx, client, http.MethodPatch, url, header, payload)
}

func Delete[RES any](ctx context.Context, client *http.Client, url string, header map[string]string) (*Response[RES], error) {
	return Do[bytes.Buffer, RES](ctx, client, http.MethodDelete, url, header, bytes.Buffer{})
}

type Response[T any] struct {
	HTTPStatusCode int
	Response       T
}

func Do[REQ, RES any](ctx context.Context, client *http.Client, method, url string, header map[string]string, payload REQ) (*Response[RES], error) {
	req, err := newRequest(ctx, method, url, header, payload)
	if err != nil {
		return nil, err
	}
	return doRequest[RES](client, req)
}

func newRequest(ctx context.Context, method, url string, header map[string]string, payload any) (*http.Request, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&payload); err != nil {
		return nil, err
	}

	var req *http.Request
	req, err := http.NewRequestWithContext(ctx, method, url, &buf)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	return req, nil
}


func doRequest[RES any](client *http.Client, req *http.Request) (*Response[RES], error) {
	var resp *http.Response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var v RES
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	response := Response[RES]{
		HTTPStatusCode: resp.StatusCode,
		Response:       v,
	}

	return &response, nil
}
