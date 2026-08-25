package main

import (
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// buildMultipartBody assembles a multipart/form-data body out of a list of
// form fields, the way curl's -F/--form does. A field is either name=value for
// a plain value or name=@path to attach a file. It returns the encoded body
// together with the Content-Type header value, which carries the boundary that
// multipart.Writer generates for us.
func buildMultipartBody(fields []string) (body, contentType string, err error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	for _, f := range fields {
		name, value, ok := strings.Cut(f, "=")
		if !ok {
			return "", "", errInvalidFormField
		}
		if name == "" {
			return "", "", errEmptyFormFieldName
		}
		if path, isFile := strings.CutPrefix(value, "@"); isFile {
			contents, rerr := os.ReadFile(path)
			if rerr != nil {
				return "", "", rerr
			}
			part, cerr := w.CreateFormFile(name, filepath.Base(path))
			if cerr != nil {
				return "", "", cerr
			}
			if _, werr := part.Write(contents); werr != nil {
				return "", "", werr
			}
			continue
		}
		if werr := w.WriteField(name, value); werr != nil {
			return "", "", werr
		}
	}
	if cerr := w.Close(); cerr != nil {
		return "", "", cerr
	}
	return buf.String(), w.FormDataContentType(), nil
}
