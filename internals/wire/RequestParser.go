package wire

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/AhmedAshraf780/wire/internals/utils"
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

//func parseQuery(path string) map[string]string {
//	params := make(map[string]string)
//
//	idx := strings.Index(path, "?")
//	if idx == -1 {
//		return params
//	}
//
//	query := path[idx+1:]
//
//	for _, pair := range strings.Split(query, "&") {
//		if pair == "" {
//			continue
//		}
//
//		kv := strings.SplitN(pair, "=", 2)
//
//		if len(kv) == 1 {
//			params[kv[0]] = ""
//			continue
//		}
//
//		params[kv[0]] = kv[1]
//	}
//
//	return params
//}

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
func ReadAndParseRequest(app *Application, reader *bufio.Reader, conn net.Conn) (*Request[[]byte], string, bool) {
	line, err := reader.ReadString('\n')
	request := &Request[[]byte]{
		Headers: map[string]string{},
	}
	if err != nil {
		if !errors.Is(err, io.EOF) {
			println("END OF FILE(CONN ENDED): ", err)
		}
		return nil, "", false
	}

	tokens, err := ParseRequestLine(line)
	if err != nil {
		// TODO: send http response
		err = utils.WriteResponse(conn, http.StatusBadRequest, "Something went wrong", []byte("Something went wrong"), map[string]string{}, request.Version)
		if err != nil {
			println("FAILED TO WRITE RESPONSE:", err)
		}
		return nil, "", false
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
		err := utils.WriteResponse(conn, http.StatusNotFound, "No method exists", []byte("Method or path not found"), map[string]string{}, request.Version)
		if err != nil {
			println("FAILED TO WRITE RESPONSE:", err)
		}
		return nil, "", false
	}
label1:
	// request headers
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				println("END OF FILE(CONN ENDED): ", err)
			}
			return nil, "", false
		}
		if line == "\r\n" {
			break
		}
		key, value, err := ParseHeader(line)
		if err != nil {
			log.Println(err)
			// TODO: send http response
			err = utils.WriteResponse(conn, http.StatusBadRequest, "Something went wrong", []byte("Something went wrong"), map[string]string{}, request.Version)
			if err != nil {
				println("FAILED TO WRITE RESPONSE:", err)
			}
			return nil, "", false
		}
		request.Headers[key] = value
	}
	// request body
	lengthStr, _ := request.Headers["Content-Length"]
	length, _ := strconv.Atoi(lengthStr)
	body := make([]byte, length)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		log.Println(err)
		err = utils.WriteResponse(conn, http.StatusBadRequest, "Something went wrong", []byte("Something went wrong"), map[string]string{}, request.Version)
		if err != nil {
			println("FAILED TO WRITE RESPONSE:", err)
		}
		return nil, "", false
	}
	request.Body = body
	request.Context = make(map[string]interface{})
	return request, orgpath, true
}
