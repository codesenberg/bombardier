package main

import (
	"bytes"
	"testing"
)

func TestHeadersToStringConversion(t *testing.T) {
	expectations := []struct {
		in  headersList
		out string
	}{
		{
			[]header{},
			"[]",
		},
		{
			[]header{
				{"Key1", "Value1"},
				{"Key2", "Value2"}},
			"[{Key1 Value1} {Key2 Value2}]",
		},
	}
	for _, e := range expectations {
		actual := e.in.String()
		expected := e.out
		if expected != actual {
			t.Errorf("Expected \"%v\", but got \"%v\"", expected, actual)
		}
	}
}

func TestShouldErrorOnInvalidFormat(t *testing.T) {
	h := new(headersList)
	if err := h.Set("Yaba daba do"); err == nil {
		t.Error("Should fail on strings without colon")
	}
}

func TestShouldProperlyAddValidHeaders(t *testing.T) {
	h := new(headersList)
	for _, hs := range []string{"Key1: Value1", "Key2: Value2"} {
		if err := h.Set(hs); err != nil {
			t.Log(err)
			t.Fail()
		}
	}
	e := []header{{"Key1", "Value1"}, {"Key2", "Value2"}}
	for i, v := range *h {
		if e[i] != v {
			t.Fail()
		}
	}
}

func TestShouldTrimHeaderValues(t *testing.T) {
	h := new(headersList)
	if err := h.Set("Key:   Value   "); err != nil {
		t.Log(err)
		t.Fail()
	}
	if (*h)[0].key != "Key" || (*h)[0].value != "Value" {
		t.Fail()
	}
}

func TestShouldProperlyConvertToFastHttpHeaders(t *testing.T) {
	h := new(headersList)
	for _, hs := range []string{
		"Content-Type: application/json", "Custom-Header: xxx42xxx",
	} {
		if err := h.Set(hs); err != nil {
			t.Log(err)
			t.Fail()
		}
	}
	fh := h.toRequestHeader()
	if e, a := []byte("application/json"), fh.Peek("Content-Type"); !bytes.Equal(e, a) {
		t.Errorf("Expected %v, but got %v", e, a)
	}
	if e, a := []byte("xxx42xxx"), fh.Peek("Custom-Header"); !bytes.Equal(e, a) {
		t.Errorf("Expected %v, but got %v", e, a)
	}
}

func TestShouldReturnNilIfNoHeadersWhereSet(t *testing.T) {
	h := new(headersList)
	if h.toRequestHeader() != nil {
		t.Fail()
	}
}
