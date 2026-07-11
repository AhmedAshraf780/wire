package utils

import (
	"fmt"
	"strconv"
	"time"
)

func MakeResponse(status int, statusMsg string, body []byte, headers map[string]string, version string) []byte {
	resp := fmt.Sprintf("%s %d %s\r\n",
		version,
		status,
		statusMsg,
	)
	bodyLength := len(body)
	headers["Connection"] = "keep-alive"
	headers["Content-Length"] = strconv.Itoa(bodyLength)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	headers["Accept-Charset"] = "utf-8"
	headers["Accept-Encoding"] = "gzip"
	headers["server"] = "Wire/1.0"
	headers["Date"] = HTTPDate()
	for key, value := range headers {
		resp += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	resp += "\r\n"
	resp += string(body)
	return []byte(resp)
}

const HTTPDateFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

func HTTPDate() string {
	return time.Now().UTC().Format(HTTPDateFormat)
}
