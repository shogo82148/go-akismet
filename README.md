# go-akismet

[![Go Reference](https://pkg.go.dev/badge/github.com/shogo82148/go-akismet.svg)](https://pkg.go.dev/github.com/shogo82148/go-akismet)

Go library for accessing the [Akismet API](https://akismet.com/).

## SYNOPSIS

```go
key := os.Getenv("AKISMET_KEY")
c := &Client{
    APIKey: key,
}

// check whether comment is a spam.
result, err := c.CheckComment(context.Background(), &Comment{
    Blog:                "https://example.com",
    UserIP:              "192.0.2.1",
    UserRole:            "administrator",
    CommentDate:         time.Unix(1234567890, 0),
    CommentPostModified: time.Unix(1234567890, 0),
    IsTest:              true,
})
if err != nil {
    fmt.Fatal(err)
}

if result.Spam {
    fmt.Println("it's SPAM")
} else {
    fmt.Printkn("it's HAM")
}
```
