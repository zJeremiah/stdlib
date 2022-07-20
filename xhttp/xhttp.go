// Package xhttp is an extension to the golang http library.
// It includes nice things like mockable http clients
// and http clients with built in stats reporting.
package xhttp

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is an interface that mimics the golang http.Client type.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Head(url string) (*http.Response, error)
	Post(url, bodyType string, body io.Reader) (*http.Response, error)
	PostForm(url string, data url.Values) (*http.Response, error)
}

// NewDefaultClient returns a new implementation of HTTPClient.
func NewDefaultClient() Client {
	return new(http.Client)
}

// NewClient returns a new implementation of HTTPClient.
func NewClient(transport http.RoundTripper, jar http.CookieJar, timeout time.Duration) Client {
	return &http.Client{
		Transport: transport,
		Jar:       jar,
		Timeout:   timeout,
	}
}
