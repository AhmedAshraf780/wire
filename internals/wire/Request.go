package wire

import (
	"strings"
)

type Request[T any] struct {
	Method  string
	Path    string
	Params  map[string]string
	Query   map[string]string
	Version string
	Headers map[string]string
	Body    T
}

type EmptyBody struct{}
type EmptyParams struct{}
type EmptyQuery struct{}

func CheckDynamicPath(routes []Route, path string) (map[string]string, map[string]string, int, string) {
	for i, r := range routes {
		params, queries, ok := match(r.Segments, path)
		if ok {
			return params, queries, i, r.Method
		}
	}

	return nil, nil, -1, ""
}

func match(segments []string, reqPath string) (map[string]string, map[string]string, bool) {
	params := make(map[string]string)
	queries := make(map[string]string)

	// Split path and query
	path, query, _ := strings.Cut(reqPath, "?")

	// Parse query parameters
	if query != "" {
		for _, pair := range strings.Split(query, "&") {
			if pair == "" {
				continue
			}

			key, value, _ := strings.Cut(pair, "=")
			queries[key] = value
		}
	}

	req := strings.Split(strings.Trim(path, "/"), "/")

	if len(req) != len(segments)-1 {
		return nil, nil, false
	}

	for i, j := 0, 1; i < len(req); i, j = i+1, j+1 {
		s := segments[j]

		if strings.HasPrefix(s, ":") {
			params[s[1:]] = req[i]
			continue
		}

		if s != req[i] {
			return nil, nil, false
		}
	}

	return params, queries, true
}

func ParseQuery(path string) map[string]string {
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
