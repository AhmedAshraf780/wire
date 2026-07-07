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
