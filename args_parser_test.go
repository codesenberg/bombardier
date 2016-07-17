package main

import (
	"bytes"
	"testing"
)

func TestArgsParsing(t *testing.T) {
	expectations := []struct {
		p   *defaultParser
		in  []string
		out error
	}{
		{
			newDefaultParser(),
			[]string{programName},
			errNoURL,
		},
		{
			newDefaultParser(),
			[]string{programName, "http://google.com", "http://yahoo.com"},
			errTooManyArgs,
		},
		{
			newDefaultParser(),
			[]string{programName, "http://google.com"},
			nil,
		},
	}
	for _, e := range expectations {
		if _, err := e.p.parse(e.in); err != e.out {
			t.Log(err, e.out)
			t.Fail()
		}
	}
}

func TestUnspecifiedArgParsing(t *testing.T) {
	p := newDefaultParser()
	args := []string{programName, "--someunspecifiedflag"}
	_, err := p.parse(args)
	if err == nil {
		t.Fail()
	}
}

func TestUsagePrinting(t *testing.T) {
	b := new(bytes.Buffer)
	p := newDefaultParser()
	p.usage(b)
	if b.Len() == 0 {
		t.Fail()
	}
}
