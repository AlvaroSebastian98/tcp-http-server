package app

import (
	r "main/router"
	"main/utils"
	"os"
	"strings"
)

func HandleRequest(router r.Router) {

	router.Get("/", func(req r.Request, res r.Response) {
		res.Send(200)
	})

	router.Get("/user-agent", func(req r.Request, res r.Response) {
		res.Headers["Content-Type"] = "text/plain"
		res.Send(200, req.UserAgent)
	})

	// Prints in body
	router.Get("/echo/:txt", func(req r.Request, res r.Response) {
		text := req.Params["txt"]
		text, contentEncoding := utils.Compress(req.HTTPRequest, text)

		if len(contentEncoding) > 0 {
			res.Headers["Content-Encoding"] = contentEncoding
		}

		res.Headers["Content-Type"] = "text/plain"
		res.Send(200, text)
	})

	// Get a file
	router.Get("/files/:filename", func(req r.Request, res r.Response) {

		if errCode, message := validateFilesParams(req.Params["filename"]); errCode >= 400 {
			res.Headers["Content-Type"] = "text/plain"
			res.Send(errCode, message)
			return
		}

		filePath := os.Args[2]
		file, err := os.ReadFile(filePath + req.Params["filename"])
		if err != nil {
			res.Send(404)
			return
		}

		contentType := "application/octet-stream"

		if (strings.HasSuffix(req.Params["filename"], ".html")) {
			contentType = "text/html"
		}

		res.Headers["Content-Type"] = contentType
		res.Send(200, string(file))
	})

	// Create a file
	router.Post("/files/:filename", func(req r.Request, res r.Response) {

		if errCode, message := validateFilesParams(req.Params["filename"]); errCode >= 400 {
			res.Headers["Content-Type"] = "text/plain"
			res.Send(errCode, message)
			return
		}

		filePath := os.Args[2]
		file, err := os.Create(filePath + req.Params["filename"])
		if err != nil {
			res.Headers["Content-Type"] = "text/plain"
			res.Send(400, "Error creating file")
			return
		}

		defer file.Close()
		file.Write([]byte(req.Body))

		res.Send(201)
	})

	router.Use("*", func (req r.Request, res r.Response)  {
		// ctx.Status(404)
		res.Send(404)
	})

}


func validateFilesParams(path string) (int, string) {

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