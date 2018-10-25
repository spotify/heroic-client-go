package heroic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	// DefaultURL is the default host to use when none is given.
	DefaultURL = "http://localhost"

	// DefaultPort is the default host to use when none is given.
	DefaultPort = 8080
)

// Client is a heroic client
type Client struct {
	// Base URL for API requests. Must include a trailing slash.
	BaseURL *url.URL

	// HTTP client used to communicate with the API.
	client *http.Client
}

// NewClient returns a new heroic client with the given baseURL and httpClient.
// If no baseURL is provided, the DefaultURL and DefaultPort will be used.
// If no httpClient is provided, the http.DefaultClient will be used.
func NewClient(baseURL *url.URL, httpClient *http.Client) *Client {
	if baseURL == nil {
		baseURL, _ = url.Parse(fmt.Sprintf("%s:%d/", DefaultURL, DefaultPort))
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		BaseURL: baseURL,
		client:  httpClient,
	}
}

func (c *Client) NewRequest(method, urlstr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlstr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("content-type", "application/json")
	}

	req.Header.Set("accept", "application/json")

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}
	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}
