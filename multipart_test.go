package main

import (
	"mime"
	"mime/multipart"
	"strings"
	"testing"
)

func TestBuildMultipartBody(t *testing.T) {
	body, contentType, err := buildMultipartBody(
		[]string{"greeting=hello", "attachment=@testbody.txt"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		t.Fatalf("Failed to parse Content-Type %q: %v", contentType, err)
	}
	if mediaType != "multipart/form-data" {
		t.Errorf("Expected multipart/form-data, but got %q", mediaType)
	}
	boundary := params["boundary"]
	if boundary == "" {
		t.Fatal("Expected an auto generated boundary, but got none")
	}

	r := multipart.NewReader(strings.NewReader(body), boundary)
	form, err := r.ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("Failed to read back the form: %v", err)
	}
	if got := form.Value["greeting"]; len(got) != 1 || got[0] != "hello" {
		t.Errorf("Expected greeting=hello, but got %v", got)
	}
	files := form.File["attachment"]
	if len(files) != 1 {
		t.Fatalf("Expected one attachment, but got %d", len(files))
	}
	if files[0].Filename != "testbody.txt" {
		t.Errorf("Expected filename testbody.txt, but got %q", files[0].Filename)
	}
}

func TestBuildMultipartBodyErrors(t *testing.T) {
	expectations := []struct {
		fields []string
		out    error
	}{
		{[]string{"nosign"}, errInvalidFormField},
		{[]string{"=value"}, errEmptyFormFieldName},
	}
	for _, e := range expectations {
		if _, _, err := buildMultipartBody(e.fields); err != e.out {
			t.Errorf("For %v, expected %v, but got %v", e.fields, e.out, err)
		}
	}
}
