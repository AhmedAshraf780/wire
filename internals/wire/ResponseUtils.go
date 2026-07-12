package wire

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func MakeResponse(status int, statusMsg string, body []byte, headers map[string][]string, version string) []byte {
	resp := fmt.Sprintf("%s %d %s\r\n",
		version,
		status,
		statusMsg,
	)
	bodyLength := len(body)
	headers[HeaderConnection] = append(headers[HeaderConnection], "keep-alive")
	headers[HeaderContentLength] = append(headers[HeaderContentLength], strconv.Itoa(bodyLength))
	headers[HeaderContentType] = append(headers[HeaderContentType], "application/json")
	headers[HeaderServer] = append(headers[HeaderServer], "Wire/1.0")
	headers["Date"] = append(headers["Date"], HTTPDate())
	for key, values := range headers {
		for _, v := range values {
			resp += fmt.Sprintf("%s: %s\r\n", key, v)
		}
	}
	resp += "\r\n"
	resp += string(body)
	return []byte(resp)
}

func MakeClosingResponse(status int, statusMsg string, body []byte, headers map[string][]string, version string) []byte {
	resp := fmt.Sprintf("%s %d %s\r\n",
		version,
		status,
		statusMsg,
	)
	bodyLength := len(body)
	headers["Connection"] = append(headers["Connection"], "closing")
	headers["Content-Length"] = append(headers["Content-Length"], strconv.Itoa(bodyLength))
	headers["Content-Type"] = append(headers["Content-Type"], "application/json")
	headers["Server"] = append(headers["Server"], "Wire/1.0")
	headers["Date"] = append(headers["Date"], HTTPDate())
	for key, value := range headers {
		for _, v := range value {
			resp += fmt.Sprintf("%s: %s\r\n", key, v)
		}
	}
	resp += "\r\n"
	resp += string(body)
	return []byte(resp)
}

func WriteResponse(
	client Client,
	status int,
	statusMsg string,
	body []byte,
	headers map[string]string,
	version string,
) error {
	if client.closed {
		return errors.New("client already closed")
	}
	bw := bufio.NewWriter(client.conn)

	_, ok := headers["Connection"]
	if !ok {
		headers["Connection"] = "keep-alive"
	}
	headers["Content-Length"] = strconv.Itoa(len(body))
	headers["Content-Type"] = "application/json"
	headers["Server"] = "Wire/1.0"
	headers["Date"] = HTTPDate()

	// Status line
	if _, err := bw.WriteString(version); err != nil {
		return err
	}
	if _, err := bw.WriteString(" "); err != nil {
		return err
	}
	if _, err := bw.WriteString(strconv.Itoa(status)); err != nil {
		return err
	}
	if _, err := bw.WriteString(" "); err != nil {
		return err
	}
	if _, err := bw.WriteString(statusMsg); err != nil {
		return err
	}
	if _, err := bw.WriteString("\r\n"); err != nil {
		return err
	}

	// Headers
	for k, v := range headers {
		if _, err := bw.WriteString(k); err != nil {
			return err
		}
		if _, err := bw.WriteString(": "); err != nil {
			return err
		}
		if _, err := bw.WriteString(v); err != nil {
			return err
		}
		if _, err := bw.WriteString("\r\n"); err != nil {
			return err
		}
	}

	// Empty line
	if _, err := bw.WriteString("\r\n"); err != nil {
		return err
	}

	// Body
	if _, err := bw.Write(body); err != nil {
		return err
	}

	return bw.Flush()
}

const HTTPDateFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

func HTTPDate() string {
	return time.Now().UTC().Format(HTTPDateFormat)
}
