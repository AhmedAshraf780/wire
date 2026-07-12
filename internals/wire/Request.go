package wire

import (
	"strings"
	"time"
)

type Cookie struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	HttpOnly bool
	Secure   bool
	SameSite SameSite
	MaxAge   int
	Expires  time.Time
}
type SameSite string

const (
	SameSiteLax    SameSite = "Lax"
	SameSiteStrict SameSite = "Strict"
	SameSiteNone   SameSite = "None"
)

type Request[T any] struct {
	Method  string
	Path    string
	Params  map[string]string
	Query   map[string]string
	Version string
	Headers map[string][]string
	Body    T
	Cookies map[string]Cookie
	Context map[string]interface{}
}

type EmptyBody struct{}

func checkDynamicPath(routes []route, path string) (map[string]string, string, int, string) {
	for i, r := range routes {
		params, ok := match(r.Segments, path)
		if ok {
			return params, r.Pattern, i, r.Method
		}
	}

	return nil, "", -1, ""
}

func match(segments []string, path string) (map[string]string, bool) {
	params := make(map[string]string)
	req := strings.Split(strings.Trim(path, "/"), "/")

	if len(req) != len(segments)-1 {
		return nil, false
	}

	for i, j := 0, 1; i < len(req); i, j = i+1, j+1 {
		s := segments[j]

		if strings.HasPrefix(s, ":") {
			params[s[1:]] = req[i]
			continue
		}

		if s != req[i] {
			return nil, false
		}
	}

	return params, true
}

func parseQuery(path string) map[string]string {
	params := make(map[string]string)

	query := strings.TrimPrefix(path, "?")
	if query == "" {
		return params
	}

	for _, part := range strings.Split(query, "&") {
		key, value, found := strings.Cut(part, "=")
		if found {
			params[key] = value
		} else {
			params[key] = ""
		}
	}

	return params
}

func parseCookies(header string) map[string]Cookie {
	cookies := make(map[string]Cookie)

	parts := strings.Split(header, ";")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		kv := strings.SplitN(part, "=", 2)

		if len(kv) != 2 {
			continue
		}

		cookies[kv[0]] = Cookie{
			Name:  kv[0],
			Value: kv[1],
		}
	}

	return cookies
}
