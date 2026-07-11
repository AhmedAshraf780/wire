package utils

import (
	"bufio"
	"net"
	"strconv"
	"time"
)

func WriteResponse(
	conn net.Conn,
	status int,
	statusMsg string,
	body []byte,
	headers map[string]string,
	version string,
) error {
	bw := bufio.NewWriter(conn)

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
