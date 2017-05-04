package main

import (
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
			t.Error(err)
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
		t.Error(err)
	}
	if (*h)[0].key != "Key" || (*h)[0].value != "Value" {
		t.Fail()
	}
}
