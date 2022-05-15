package gateway

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// NewRequest returns a new http.Request from the given Lambda event.
func NewRequest(ctx context.Context, e events.LambdaFunctionURLRequest) (*http.Request, error) {

	rCtx := e.RequestContext

	// path
	urlString := e.RawPath
	if len(e.RawQueryString) > 0 {
		urlString += fmt.Sprintf("?%s", e.RawQueryString)
	}

	u, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("parse path and query into url: %w", err)
	}

	// base64 encoded body
	body := e.Body
	if e.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, fmt.Errorf("decoding base64 body: %w", err)
		}
		body = string(b)
	}

	// new request
	req, err := http.NewRequest(rCtx.HTTP.Method, u.String(), strings.NewReader((body)))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// manually set RequestURI because NewRequest is for clients and req.RequestURI is for servers
	req.RequestURI = u.RequestURI()

	// remote addr
	req.RemoteAddr = rCtx.HTTP.SourceIP

	// header fields
	for k, v := range e.Headers {
		req.Header.Set(k, v)
	}

	// content-length
	if req.Header.Get("Content-Length") == "" && body != "" {
		req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	}

	// custom fields
	req.Header.Set("X-Request-Id", rCtx.RequestID)

	// custom context values
	req = req.WithContext(newContext(ctx, e))

	// xray support
	if traceID := ctx.Value("x-amzn-trace-id"); traceID != nil {
		req.Header.Set("X-Amzn-Trace-Id", fmt.Sprintf("%v", traceID))
	}

	// host
	req.URL.Host = req.Header.Get("Host")
	req.Host = req.URL.Host

	return req, nil
}
