package akismet

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestVerifyKey_Integrated(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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

func TestCheckComment_Ham_Integrated(t *testing.T) {
	t.Parallel()
	key := os.Getenv("AKISMET_KEY")
	if key == "" {
		t.Skip("AKISMET_KEY is not set")
	}

	c := &Client{
		APIKey: key,
	}
	result, err := c.CheckComment(context.Background(), &Comment{
		Blog:                "https://example.com",
		UserIP:              "192.0.2.1",
		UserRole:            "administrator",
		CommentDate:         time.Unix(1234567890, 0),
		CommentPostModified: time.Unix(1234567890, 0),
		IsTest:              true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.Spam {
		t.Error("got spam, want ham")
	}
}

func TestCheckComment_Spam_Integrated(t *testing.T) {
	t.Parallel()
	key := os.Getenv("AKISMET_KEY")
	if key == "" {
		t.Skip("AKISMET_KEY is not set")
	}

	c := &Client{
		APIKey: key,
	}
	result, err := c.CheckComment(context.Background(), &Comment{
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

	if !result.Spam {
		t.Error("got ham, want spam")
	}
}

func TestSubmitHam_Integrated(t *testing.T) {
	t.Parallel()
	key := os.Getenv("AKISMET_KEY")
	if key == "" {
		t.Skip("AKISMET_KEY is not set")
	}

	c := &Client{
		APIKey: key,
	}
	err := c.SubmitHam(context.Background(), &Comment{
		Blog:   "https://example.com",
		UserIP: "192.0.2.1",
		IsTest: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSubmitSpam_Integrated(t *testing.T) {
	t.Parallel()
	key := os.Getenv("AKISMET_KEY")
	if key == "" {
		t.Skip("AKISMET_KEY is not set")
	}

	c := &Client{
		APIKey: key,
	}
	err := c.SubmitSpam(context.Background(), &Comment{
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
}
