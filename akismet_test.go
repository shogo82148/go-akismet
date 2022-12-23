package akismet

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyKey_Valid(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/verify-key"; got != want {
			t.Errorf("unexpected path: got %s, want %s", got, want)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			t.Error(err)
			return
		}
		if got, want := string(body), "api_key=very-secret&blog=http%3A%2F%2Fexample.com"; want != got {
			t.Errorf("got %s, want %s", got, want)
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "valid")
	}))
	defer ts.Close()

	c := &Client{
		HTTPClient: ts.Client(),
		BaseURL:    ts.URL,
		APIKey:     "very-secret",
	}
	if err := c.VerifyKey(context.Background(), "http://example.com"); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyKey_Invalid(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "invalid")
	}))
	defer ts.Close()

	c := &Client{
		HTTPClient: ts.Client(),
		BaseURL:    ts.URL,
		APIKey:     "very-secret",
	}
	if err := c.VerifyKey(context.Background(), "http://example.com"); err == nil {
		t.Error("want error, but not")
	}
}

func TestCheck_HAM(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/comment-check"; got != want {
			t.Errorf("unexpected path: got %s, want %s", got, want)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			t.Error(err)
			return
		}
		if got, want := string(body), "api_key=very-secret&blog=https%3A%2F%2Fexample.com&is_test=1&user_ip=192.0.2.1&user_role=administrator"; want != got {
			t.Errorf("got %s, want %s", got, want)
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "false")
	}))
	defer ts.Close()

	c := &Client{
		HTTPClient: ts.Client(),
		BaseURL:    ts.URL,
		APIKey:     "very-secret",
	}
	ret, err := c.CheckComment(context.Background(), &Comment{
		Blog:     "https://example.com",
		UserIP:   "192.0.2.1",
		UserRole: "administrator",
		IsTest:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ret.Spam {
		t.Error("want Ham, but got spam")
	}
}

func TestCheck_SPAM(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/comment-check"; got != want {
			t.Errorf("unexpected path: got %s, want %s", got, want)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			t.Error(err)
			return
		}
		if got, want := string(body), "api_key=very-secret&blog=https%3A%2F%2Fexample.com&comment_author=viagra-test-123&"+
			"comment_author_email=akismet-guaranteed-spam%40example.com&is_test=1&user_ip=192.0.2.1"; want != got {
			t.Errorf("got %s, want %s", got, want)
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "true")
	}))
	defer ts.Close()

	c := &Client{
		HTTPClient: ts.Client(),
		BaseURL:    ts.URL,
		APIKey:     "very-secret",
	}
	ret, err := c.CheckComment(context.Background(), &Comment{
		Blog:   "https://example.com",
		UserIP: "192.0.2.1",
		IsTest: true,

		// known-spammer. ref. https://akismet.com/development/api/#detailed-docs
		CommentAuthor:      "viagra-test-123",
		CommentAuthorEmail: "akismet-guaranteed-spam@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !ret.Spam {
		t.Error("want SPAM, but got HAM")
	}
}

func TestCheck_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "some error")
	}))
	defer ts.Close()

	c := &Client{
		HTTPClient: ts.Client(),
		BaseURL:    ts.URL,
		APIKey:     "very-secret",
	}
	_, err := c.CheckComment(context.Background(), &Comment{
		Blog:   "https://example.com",
		UserIP: "192.0.2.1",
		IsTest: true,
	})
	if err == nil {
		t.Error("want some error, but not")
	}
}

func TestSubmitHam(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/submit-ham"; got != want {
			t.Errorf("unexpected path: got %s, want %s", got, want)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			t.Error(err)
			return
		}
		if got, want := string(body), "api_key=very-secret&blog=https%3A%2F%2Fexample.com&is_test=1&user_ip=192.0.2.1&user_role=administrator"; want != got {
			t.Errorf("got %s, want %s", got, want)
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "false")
	}))
	defer ts.Close()

	c := &Client{
		HTTPClient: ts.Client(),
		BaseURL:    ts.URL,
		APIKey:     "very-secret",
	}
	err := c.SubmitHam(context.Background(), &Comment{
		Blog:     "https://example.com",
		UserIP:   "192.0.2.1",
		UserRole: "administrator",
		IsTest:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSubmitSpam(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/submit-spam"; got != want {
			t.Errorf("unexpected path: got %s, want %s", got, want)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			t.Error(err)
			return
		}
		if got, want := string(body), "api_key=very-secret&blog=https%3A%2F%2Fexample.com&is_test=1&user_ip=192.0.2.1&user_role=administrator"; want != got {
			t.Errorf("got %s, want %s", got, want)
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "false")
	}))
	defer ts.Close()

	c := &Client{
		HTTPClient: ts.Client(),
		BaseURL:    ts.URL,
		APIKey:     "very-secret",
	}
	err := c.SubmitSpam(context.Background(), &Comment{
		Blog:     "https://example.com",
		UserIP:   "192.0.2.1",
		UserRole: "administrator",
		IsTest:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
}
