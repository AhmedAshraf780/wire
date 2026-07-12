package wire

import (
	"fmt"
	"net/http"
	"strings"
)

type Response[T any] struct {
	StatusCode int
	Body       T
	Version    string
	/* headers now don't support duplicate headers like
	/* set-Cookie: a=1
	/* set-Cookie: b=2
	*/
	Headers map[string][]string
}

func (r *Response[T]) Header(key string) string {
	return r.Headers[key][0]
}
func (r *Response[T]) SetHeader(key string, value string) {
	r.Headers[key] = append(r.Headers[key], value)
}
func (r *Response[T]) Write(statusCode int, body T) error {
	r.StatusCode = statusCode
	r.Body = body
	return nil
}

func (r *Response[T]) SetCookie(c Cookie) {
	parts := []string{
		fmt.Sprintf("%s=%s", c.Name, c.Value),
	}

	if c.Path != "" {
		parts = append(parts, "Path="+c.Path)
	}

	if c.Domain != "" {
		parts = append(parts, "Domain="+c.Domain)
	}

	if c.MaxAge > 0 {
		parts = append(parts, fmt.Sprintf("Max-Age=%d", c.MaxAge))
	}

	if !c.Expires.IsZero() {
		parts = append(parts, "Expires="+c.Expires.UTC().Format(http.TimeFormat))
	}

	if c.HttpOnly {
		parts = append(parts, "HttpOnly")
	}

	if c.Secure {
		parts = append(parts, "Secure")
	}

	if c.SameSite != "" {
		parts = append(parts, "SameSite="+string(c.SameSite))
	}

	r.Headers["Set-Cookie"] = append(
		r.Headers["Set-Cookie"],
		strings.Join(parts, "; "),
	)
}
