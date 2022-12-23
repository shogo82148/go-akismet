package akismet

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	if string(respBody) != "valid" {
		return fmt.Errorf("akismet: your api key is %s", respBody)
	}

	return nil
}

type Comment struct {
	// Blog is The front page or home URL of the instance making the request.
	// For a blog or wiki this would be the front page. Note: Must be a full URI, including http://.
	Blog string

	// UserIP is IP address of the comment submitter.
	UserIP string

	// UserAgent is the user agent string of the web browser submitting the comment.
	// Typically the HTTP_USER_AGENT cgi variable. Not to be confused with the user agent of your Akismet library.
	UserAgent string

	// Referrer is the content of the HTTP_REFERER header should be sent here.
	Referrer string

	// Permalink is the full permanent URL of the entry the comment was submitted to.
	Permalink string

	// CommentType is a string that describes the type of content being sent.
	CommentType CommentType

	CommentAuthor string

	CommentAuthorEmail string

	CommentAuthorURL string

	CommentContent string

	CommentDate time.Time

	CommentPostModified time.Time

	BlogLang string

	BlogCharset string

	UserRole string

	IsTest bool

	RecheckReason string

	HoneypotFieldName string
}

type CommentType string

const (
	// CommentTypeComment is a blog comment.
	CommentTypeComment CommentType = "comment"

	// CommentTypeForumPost is a top-level forum post.
	CommentTypeForumPost CommentType = "forum-post"

	// CommentTypeReply is reply to a top-level forum post.
	CommentTypeReply CommentType = "reply"

	// CommentTypeBlogPost is a blog post.
	CommentTypeBlogPost CommentType = "blog-post"

	// CommentTypeContactForm is a contact form or feedback form submission.
	CommentTypeContactForm CommentType = "contact-form"

	// CommentTypeSignUp is new user account.
	CommentTypeSignUp CommentType = "signup"

	// CommentTypeMessage is a message sent between just a few users.
	CommentTypeMessage CommentType = "message"
)

type Result struct {
	Spam bool
}

func (c *Client) CheckComment(ctx context.Context, comment *Comment) (*Result, error) {
	// build the request.
	u, err := c.resolvePath("comment-check")
	if err != nil {
		return nil, err
	}
	form := url.Values{}
	form.Set("api_key", c.APIKey)
	form.Set("blog", comment.Blog)
	form.Set("user_ip", comment.UserIP)
	if comment.UserAgent != "" {
		form.Set("user_agent", comment.UserAgent)
	}
	if comment.Referrer != "" {
		form.Set("referrer", comment.Referrer)
	}
	if comment.Permalink != "" {
		form.Set("permalink", comment.Permalink)
	}
	if comment.CommentType != "" {
		form.Set("comment_type", string(comment.CommentType))
	}
	if comment.CommentAuthor != "" {
		form.Set("comment_author", comment.CommentAuthor)
	}
	if comment.CommentAuthorEmail != "" {
		form.Set("comment_author_email", comment.CommentAuthorEmail)
	}
	if comment.CommentAuthorURL != "" {
		form.Set("comment_author_url", comment.CommentAuthorURL)
	}
	if comment.CommentContent != "" {
		form.Set("comment_content", comment.CommentContent)
	}
	if !comment.CommentDate.IsZero() {
		form.Set("comment_date_gmt", comment.CommentDate.UTC().Format(time.RFC3339))
	}
	if !comment.CommentPostModified.IsZero() {
		form.Set("comment_post_modified_gmt", comment.CommentPostModified.UTC().Format(time.RFC3339))
	}
	if comment.BlogLang != "" {
		form.Set("blog_lang", comment.BlogLang)
	}
	if comment.BlogCharset != "" {
		form.Set("blog_charset", comment.BlogCharset)
	}
	if comment.UserRole != "" {
		form.Set("user_role", comment.UserRole)
	}
	if comment.IsTest {
		form.Set("is_test", "1")
	}
	if comment.RecheckReason != "" {
		form.Set("recheck_reason", comment.RecheckReason)
	}
	if comment.HoneypotFieldName != "" {
		form.Set("honeypot_field_name", comment.HoneypotFieldName)
	}
	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// send the request.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse the response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("akismet: unexpected status code: %d", resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respBody = bytes.TrimSpace(respBody)
	if string(respBody) == "true" {
		return &Result{
			Spam: true,
		}, nil
	}
	if string(respBody) == "false" {
		return &Result{
			Spam: false,
		}, nil
	}

	return nil, fmt.Errorf("akismet: error from the server: %s", respBody)
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
