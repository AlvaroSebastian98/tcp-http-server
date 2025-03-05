package app

import (
	"main/utils"
	"os"
	"strings"
)

func HandleRequest(router utils.Router, req utils.HTTPRequest, res utils.HTTPResponse) {

	router.Get("/", func(ctx utils.RouteContext) {
		res.Send(200, "")
	})

	router.Get("/user-agent", func(ctx utils.RouteContext) {
		res.Headers["Content-Type"] = "text/plain"
		res.Send(200, req.UserAgent)
	})

	// Prints in body
	router.Get("/echo/:txt", func(ctx utils.RouteContext) {
		text := ctx.Params["txt"]
		text, contentEncoding := utils.Compress(req, text)

		if len(contentEncoding) > 0 {
			res.Headers["Content-Encoding"] = contentEncoding
		}

		res.Headers["Content-Type"] = "text/plain"
		res.Send(200, text)
	})

	// Get a file
	router.Get("/files/:filename", func(ctx utils.RouteContext) {

		if errCode, message := validateFilesParams(ctx.Params["filename"]); errCode >= 400 {
			res.Headers["Content-Type"] = "text/plain"
			res.Send(errCode, message)
			return
		}

		filePath := os.Args[2]
		file, err := os.ReadFile(filePath + ctx.Params["filename"])
		if err != nil {
			res.Send(404, "")
			return
		}

		contentType := "application/octet-stream"

		if (strings.HasSuffix(ctx.Params["filename"], ".html")) {
			contentType = "text/html"
		}

		res.Headers["Content-Type"] = contentType
		res.Send(200, string(file))
	})

	// Create a file
	router.Post("/files/:filename", func(ctx utils.RouteContext) {

		if errCode, message := validateFilesParams(ctx.Params["filename"]); errCode >= 400 {
			res.Headers["Content-Type"] = "text/plain"
			res.Send(errCode, message)
			return
		}

		filePath := os.Args[2]
		file, err := os.Create(filePath + ctx.Params["filename"])
		if err != nil {
			res.Headers["Content-Type"] = "text/plain"
			res.Send(400, "Error creating file")
			return
		}

		defer file.Close()
		file.Write([]byte(req.Body))

		res.Send(201, "")
	})

	router.Use("*", func (ctx utils.RouteContext)  {
		// ctx.Status(404)
		res.Send(404, "")
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