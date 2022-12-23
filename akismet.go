package akismet

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	version          = "0.1.0"
	defaultUserAgent = "shogo82148-go-akismet"
	defaultBaseURL   = "https://rest.akismet.com/1.1/"
)

// HTTPClient is a interface for http client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	// HTTPClient is underlying http client.
	// If it is nil, http.DefaultClient is used.
	HTTPClient HTTPClient

	// UserAgent is User-Agent header of requests.
	UserAgent string

	// APIKey is an API Key for using Akismet API.
	APIKey string

	// BaseURL is the endpoint of Akismet API.
	// If is is empty, https://rest.akismet.com/1.1/ is used.
	BaseURL string
}

type Comment struct {
}

type Result struct {
	OK bool
}

func (c *Client) VerifyKey(ctx context.Context, blog string) error {
	// build the request.
	u, err := c.resolvePath("verify-key")
	if err != nil {
		return err
	}
	form := url.Values{}
	form.Set("api_key", c.APIKey)
	form.Set("blog", blog)
	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// send the request.
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// parse the response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("akismet: unexpected status code: %d", resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respBody = bytes.TrimSpace(respBody)
	if !bytes.Equal(respBody, []byte("valid")) {
		return fmt.Errorf("akismet: your api key is %s", respBody)
	}

	return nil
}

func (c *Client) CheckComment(ctx context.Context, comment *Comment) (*Result, error) {
	// build the request.
	u, err := c.resolvePath("comment-check")
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader([]byte{}) // TODO: fill me
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return nil, err
	}

	// send the request.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &Result{}, nil
}

func (c *Client) resolvePath(path string) (*url.URL, error) {
	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return base.JoinPath(path), nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	} else {
		req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", defaultUserAgent, version))
	}
	if c.HTTPClient == nil {
		return http.DefaultClient.Do(req)
	}
	return c.HTTPClient.Do(req)
}
