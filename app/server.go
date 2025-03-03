package app

import (
	// "fmt"
	"main/utils"
	"net"
	"os"
	"strings"
	// "time"
)

func HandleRequest(conn net.Conn, req utils.HTTPRequest, res utils.HTTPResponse) {

	defer conn.Close()

	if req.Path == "/" {
		res.Send(200, "")
		return
	}

	if req.Path == "/user-agent" {
		res.Headers["Content-Type"] = "text/plain"
		res.Send(200, req.UserAgent)
		return
	}

	// Prints in body
	if strings.HasPrefix(req.Path , "/echo") {
		path := strings.Split(strings.Replace(req.Path, "/echo", "", 1), "/")[1:]

		if len(path) == 1 {
			text := path[0]
			text, contentEncoding := utils.Compress(req, text)

			if len(contentEncoding) > 0 {
				res.Headers["Content-Encoding"] = contentEncoding
			}

			res.Headers["Content-Type"] = "text/plain"
			res.Send(200, text)
			return
		}

	}


	// Get a file
	if req.Method == "GET" && strings.HasPrefix(req.Path , "/files") {
		path := strings.Split(strings.Replace(req.Path, "/files", "", 1), "/")[1:]

		if errCode, message := validateFilesParams(path); errCode >= 400 {
			res.Headers["Content-Type"] = "text/plain"
			res.Send(errCode, message)
			return
		}

		filename := path[0]
		filePath := os.Args[2]
		file, err := os.ReadFile(filePath + filename)
		if err != nil {
			res.Send(404, "")
			return
		}

		contentType := "application/octet-stream"

		if (strings.HasSuffix(filename, ".html")) {
			contentType = "text/html"
		}

		res.Headers["Content-Type"] = contentType
		res.Send(200, string(file))
		return
	}


	// Create a file
	if req.Method == "POST" && strings.HasPrefix(req.Path , "/files") {
		path := strings.Split(strings.Replace(req.Path, "/files", "", 1), "/")[1:]

		if errCode, message := validateFilesParams(path); errCode >= 400 {
			res.Headers["Content-Type"] = "text/plain"
			res.Send(errCode, message)
		}

		filename := path[0]
		filePath := os.Args[2]
		file, err := os.Create(filePath + filename)
		if err != nil {
			res.Headers["Content-Type"] = "text/plain"
			res.Send(400, "Error creating file")
			return
		}

		defer file.Close()
		file.Write([]byte(req.Body))

		res.Send(201, "")
		return
	}

	res.Send(404, "")

}


func validateFilesParams(path []string) (int, string) {

	if len(path) == 0 {
		code := 400
		message := "filename is required"
		return code, message
	}

	if len(os.Args) < 3 || os.Args[1] != "--directory" {
		code := 400
		message := "folder path after --directory is required"
		return code, message
	}

	return 200, "OK"

}