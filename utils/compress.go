package utils

import (
	"bytes"
	"compress/gzip"
	"strings"
	"slices"
)

var supportedEncoding []string = []string{"gzip"}

func Compress(req HTTPRequest, content string) (string, string) {
	var contentEncoding string

	// If receive some compress encoding
	if len(req.AcceptEncoding) > 0 {
		acceptedEncodings := []string{}

		for _, s := range strings.Split(req.AcceptEncoding, ",") {
			encoding := strings.TrimSpace(s)
			if slices.Contains(supportedEncoding, encoding) {
				acceptedEncodings = append(acceptedEncodings, encoding)
			}
		}

		if len(acceptedEncodings) > 0 {
			contentEncoding = strings.Join(acceptedEncodings, ", ")

			if acceptedEncodings[0] == "gzip" {
				var b bytes.Buffer
				w := gzip.NewWriter(&b)
				w.Write([]byte(content))
				w.Close()
				content = b.String()
			}
		}
	}

	return content, contentEncoding

}