package utils

import "fmt"

func MakeResponse(status int, statusMsg string, body []byte, headers map[string]string, version string) []byte {
	resp := fmt.Sprintf("%s %d %s\r\n",
		version,
		status,
		statusMsg,
	)
	for key, value := range headers {
		resp += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	resp += "\r\n"
	resp += string(body)
	return []byte(resp)
}

func DefaultHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "text/plain; charset=utf-8"
	headers["Content-Length"] = "0"
	return headers
}
