package gateway

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tj/assert"
)

func TestNewRequest_path(t *testing.T) {
	e := events.LambdaFunctionURLRequest{
		RawPath: "/pets/luna",
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, "GET", r.Method)
	assert.Equal(t, `/pets/luna`, r.URL.Path)
	assert.Equal(t, `/pets/luna`, r.URL.String())
	assert.Equal(t, `/pets/luna`, r.RequestURI)
}

func TestNewRequest_method(t *testing.T) {
	e := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "DELETE",
			},
		},
		RawPath: "/pets/luna",
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, "DELETE", r.Method)
}

func TestNewRequest_queryString(t *testing.T) {
	e := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "GET",
			},
		},
		RawPath:        "/pets",
		RawQueryString: "fields=name,species&order=desc",
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, []string{"name,species"}, r.URL.Query()["fields"])
	assert.Equal(t, `desc`, r.URL.Query().Get("order"))
}

func TestNewRequest_multiValueQueryString(t *testing.T) {
	e := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "GET",
			},
		},
		RawPath:        "/pets",
		RawQueryString: "fields=name&fields=species",
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, []string{"name", "species"}, r.URL.Query()["fields"])
}

func TestNewRequest_remoteAddr(t *testing.T) {
	e := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method:   "GET",
				SourceIP: "1.2.3.4",
			},
		},
		RawPath: "/pets",
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `1.2.3.4`, r.RemoteAddr)
}

func TestNewRequest_header(t *testing.T) {
	e := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			RequestID: "1234",
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "POST",
			},
		},
		RawPath: "/pets",
		Body:    `{ "name": "Tobi" }`,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Foo":        "bar",
			"Host":         "example.com",
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `example.com`, r.Host)
	assert.Equal(t, `1234`, r.Header.Get("X-Request-Id"))
	assert.Equal(t, `18`, r.Header.Get("Content-Length"))
	assert.Equal(t, `application/json`, r.Header.Get("Content-Type"))
	assert.Equal(t, `bar`, r.Header.Get("X-Foo"))
}

func TestNewRequest_multiHeader(t *testing.T) {
	e := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			RequestID: "1234",
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "POST",
			},
		},
		RawPath:        "/pets",
		RawQueryString: "",
		Body:           `{ "name": "Tobi" }`,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Foo":        "bar",
			"Host":         "example.com",
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `example.com`, r.Host)
	assert.Equal(t, `1234`, r.Header.Get("X-Request-Id"))
	assert.Equal(t, `18`, r.Header.Get("Content-Length"))
	assert.Equal(t, `application/json`, r.Header.Get("Content-Type"))
	assert.Equal(t, `bar`, r.Header.Get("X-Foo"))
}

func TestNewRequest_body(t *testing.T) {
	e := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "POST",
			},
		},
		RawPath: "/pets",
		Body:    `{ "name": "Tobi" }`,
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	assert.Equal(t, `{ "name": "Tobi" }`, string(b))
}

func TestNewRequest_bodyBinary(t *testing.T) {
	e := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: "POST",
			},
		},
		RawPath:         "/pets",
		Body:            `aGVsbG8gd29ybGQK`,
		IsBase64Encoded: true,
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	assert.Equal(t, "hello world\n", string(b))
}

func TestNewRequest_context(t *testing.T) {
	e := events.LambdaFunctionURLRequest{}
	ctx := context.WithValue(context.Background(), "key", "value")
	r, err := NewRequest(ctx, e)
	assert.NoError(t, err)
	v := r.Context().Value("key")
	assert.Equal(t, "value", v)
}
