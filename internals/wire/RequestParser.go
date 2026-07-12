package wire

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/AhmedAshraf780/wire/internals/utils"
)

func parseRequestLine(line string) ([]string, error) {
	line = strings.TrimRight(line, "\r\n")
	tokens := strings.Fields(line)
	if len(tokens) < 3 {
		return nil, errors.New(fmt.Sprintf("Invalid Request Line: %s", line))
	}
	return tokens, nil
}

func parseHeader(line string) (string, string, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New(fmt.Sprintf("Invalid Header: %s", line))
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	return key, value, nil
}

func staticPath(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			return false
		}
	}
	return true
}
func readAndParseRequest(app *Application, reader *bufio.Reader, client Client) (*Request[[]byte], string, bool) {
	line, err := reader.ReadString('\n')
	request := &Request[[]byte]{
		Headers: map[string][]string{},
	}
	if err != nil {
		if !errors.Is(err, io.EOF) {
			client.conn.Close()
			client.closed = true
			return nil, "", false
		}

		resp := MakeResponse(http.StatusBadRequest, BadRequest, []byte(BadRequest), map[string][]string{}, request.Version)
		client.conn.Write(resp)
		return nil, "", true
	}

	tokens, err := parseRequestLine(line)
	if err != nil {
		resp := MakeResponse(http.StatusBadRequest, BadRequest, []byte(BadRequest), map[string][]string{}, request.Version)
		client.conn.Write(resp)
		return nil, "", true
	}

	request.Method = tokens[0]
	request.Path = tokens[1]
	request.Version = tokens[2]
	// here we can check if we have a handler for this path
	path, query, found := strings.Cut(tokens[1], "?")
	if found {
		request.Query = parseQuery(query)
	}
	key := utils.GenerateHandlerKey(request.Method, path)
	_, ok := app.staticRoutes[key]
	orgpath := request.Path
	if !ok {
		// check in dynamic routes
		params, org, idx, meth := checkDynamicPath(app.dynamicRoutes, path)
		if idx != -1 && meth == request.Method {
			request.Params = params
			orgpath = org
			goto label1
		}

		resp := MakeResponse(http.StatusNotFound, NotFound, []byte(NotFound), map[string][]string{}, request.Version)
		client.conn.Write(resp)

		return nil, "", true
	}
label1:
	// request headers
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				client.conn.Close()
				client.closed = true
				return nil, "", false
			}
			resp := MakeResponse(http.StatusBadRequest, BadRequest, []byte(BadRequest), map[string][]string{}, request.Version)
			client.conn.Write(resp)
			return nil, "", true
		}
		if line == "\r\n" {
			break
		}
		key, value, err := parseHeader(line)
		if err != nil {
			log.Println(err)
			resp := MakeResponse(http.StatusBadRequest, BadRequest, []byte(BadRequest), map[string][]string{}, request.Version)
			client.conn.Write(resp)
			return nil, "", true
		}
		request.Headers[key] = append(request.Headers[key], value)
		if key == HeaderCookie {
			cookies := parseCookies(value)
			request.Cookies = cookies
		}
	}
	// request body
	lengthStr := ""
	if len(request.Headers[HeaderContentLength]) > 0 {
		lengthStr = request.Headers[HeaderContentLength][0]
	}
	length, _ := strconv.Atoi(lengthStr)
	body := make([]byte, length)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		log.Println(err)
		resp := MakeResponse(http.StatusBadRequest, BadRequest, []byte(BadRequest), map[string][]string{}, request.Version)
		client.conn.Write(resp)
		return nil, "", true
	}
	request.Body = body
	request.Context = make(map[string]interface{})
	// check if connection is alive
	if len(request.Headers[HeaderConnection]) > 0 && request.Headers[HeaderConnection][0] == "close" {
		client.closed = true
		return nil, "", true
	}
	return request, orgpath, true
}
