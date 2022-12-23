package akismet

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func TestCheck(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
