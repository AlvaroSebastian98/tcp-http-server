package utils

import (
	"strings"
	"strconv"
	"fmt"
	"net"
)

type HTTPRequest struct {
	Method string
	Path string
	Protocol string
	Host string
	UserAgent string
	Accept string
	ContentType string
	ContentLength int
	Body string
	AcceptEncoding string
}

func ParseRequest(request []string) HTTPRequest {

	fmt.Printf("%#v\n", request)

	var httpRequest HTTPRequest

	for _, line := range request {
		if strings.HasPrefix(line, "POST") || strings.HasPrefix(line, "GET") || strings.HasPrefix(line, "PUT") {
			requestLine := strings.Split(request[0], " ")

			httpRequest.Method = requestLine[0]
			httpRequest.Path = requestLine[1]
			httpRequest.Protocol = requestLine[2]
		} else if strings.HasPrefix(line, "Host:") {
			httpRequest.Host = strings.TrimPrefix(line, "Host: ")
		} else if strings.HasPrefix(line, "User-Agent:") {
			httpRequest.UserAgent = strings.TrimPrefix(line, "User-Agent: ")
		} else if strings.HasPrefix(line, "Accept:") {
			httpRequest.Accept = strings.TrimPrefix(line, "Accept: ")
		} else if strings.HasPrefix(line, "Content-Type:") {
			httpRequest.ContentType = strings.TrimPrefix(line, "Content-Type: ")
		} else if strings.HasPrefix(line, "Content-Length:") {
			httpRequest.ContentLength, _ = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:")))

			if httpRequest.ContentLength > 0 {
				httpRequest.Body = request[len(request) - 1]
			}
		} else if strings.HasPrefix(line, "Accept-Encoding:") {
			httpRequest.AcceptEncoding = strings.TrimPrefix(line, "Accept-Encoding: ")
		}
	}

	fmt.Printf("\n%#v\n", httpRequest)

	return httpRequest

}

type HTTPResponse struct {
	Send func(code int, body string)
	Headers map[string]string
}

var httpCode = map[int]string{
	200: "OK",
	201: "Created",
	400: "Bad Request",
	404: "Not Found",
	500: "Internal Server Error",
}

func BuildResponse(conn net.Conn) HTTPResponse {

	var httpResponse HTTPResponse

	httpResponse.Headers = make(map[string]string)

	httpResponse.Send = func(code int, body string) {
		response := fmt.Sprintf("HTTP/1.1 %d %v\r\n", code, httpCode[code])

		if len(body) > 0 {
			httpResponse.Headers["Content-Length"] = strconv.Itoa(len([]byte(body)))
		}

		for k, v := range httpResponse.Headers {
			response += fmt.Sprintf("%v: %v\r\n", k, v)
		}

		response += "\r\n" + body
		conn.Write([]byte(response))
	}

	return httpResponse

}