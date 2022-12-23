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
