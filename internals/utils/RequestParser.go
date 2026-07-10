package utils

import (
	"errors"
	"fmt"
	"strings"
)

func ParseRequestLine(line string) ([]string, error) {
	line = strings.TrimRight(line, "\r\n")
	tokens := strings.Fields(line)
	if len(tokens) < 3 {
		return nil, errors.New(fmt.Sprintf("Invalid Request Line: %s", line))
	}
	return tokens, nil
}

func ParseHeader(line string) (string, string, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New(fmt.Sprintf("Invalid Header: %s", line))
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	return key, value, nil
}

func ExtractParams(path, reqPath string) (map[string]string, error) {
	patternParts := strings.Split(strings.Trim(path, "/"), "/")
	pathParts := strings.Split(strings.Trim(reqPath, "/"), "/")

	params := make(map[string]string)

	if len(patternParts) != len(pathParts) {
		return nil, errors.New(fmt.Sprintf("Invalid Request Path: %s", path))
	}

	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], ":") {
			key := patternParts[i][1:] // Remove ':'
			params[key] = pathParts[i]
		} else if patternParts[i] != pathParts[i] {
			return nil, errors.New(fmt.Sprintf("Invalid Request Path: %s", path))
		}
	}
	return params, nil
}

func ParseQuery(path string) map[string]string {
	params := make(map[string]string)

	idx := strings.Index(path, "?")
	if idx == -1 {
		return params
	}

	query := path[idx+1:]

	for _, pair := range strings.Split(query, "&") {
		if pair == "" {
			continue
		}

		kv := strings.SplitN(pair, "=", 2)

		if len(kv) == 1 {
			params[kv[0]] = ""
			continue
		}

		params[kv[0]] = kv[1]
	}

	return params
}

func GenerateDynamicPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[i] = "*"
		}
	}
	return strings.Join(parts, "/")
}

func ValidParamsPath(path, reqPath string) bool {
	pathParts := strings.Split(reqPath, "/")
	reqPathParts := strings.Split(reqPath, "/")
	if len(pathParts) != len(reqPathParts) {
		return false
	}

	for i := range pathParts {
		if pathParts[i] != reqPathParts[i] {
			if pathParts[i] != "*" {
				return false
			}
		}
	}
	return true
}

func StaticPath(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			return false
		}
	}
	return true
}
