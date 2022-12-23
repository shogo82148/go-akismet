package akismet

import (
	"context"
	"os"
	"testing"
)

func TestVerifyKey_Integrated(t *testing.T) {
	key := os.Getenv("AKISMET_KEY")
	if key == "" {
		t.Skip("AKISMET_KEY is not set")
	}

	c := &Client{
		APIKey: key,
	}
	if err := c.VerifyKey(context.Background(), "http://example.com"); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyKey_Integrated_Fail(t *testing.T) {
	key := os.Getenv("AKISMET_KEY")
	if key == "" {
		t.Skip("AKISMET_KEY is not set")
	}

	c := &Client{
		APIKey: key + "_foo",
	}
	err := c.VerifyKey(context.Background(), "http://example.com")
	if err == nil {
		t.Errorf("want error, but not")
	}
}
