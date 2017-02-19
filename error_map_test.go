package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestErrorMapAdd(t *testing.T) {
	m := newErrorMap()
	err := errors.New("add")
	m.add(err)
	if c := m.get(err); c != 1 {
		t.Error(c)
	}
}

func TestErrorMapGet(t *testing.T) {
	m := newErrorMap()
	err := errors.New("get")
	if c := m.get(err); c != 0 {
		t.Error(c)
	}
}

func TestByFrequency(t *testing.T) {
	m := newErrorMap()
	a := errors.New("A")
	b := errors.New("B")
	c := errors.New("C")
	m.add(a)
	m.add(a)
	m.add(b)
	m.add(b)
	m.add(b)
	m.add(c)
	e := errorsByFrequency{
		{"B", 3},
		{"A", 2},
		{"C", 1},
	}
	if a := m.byFrequency(); !reflect.DeepEqual(a, e) {
		t.Logf("Expected: %+v", e)
		t.Logf("Got: %+v", a)
		t.Fail()
	}
}

func TestErrorWithCountToStringConversion(t *testing.T) {
	ewc := errorWithCount{"A", 1}
	exp := "<A:1>"
	if act := ewc.String(); act != exp {
		t.Logf("Expected: %+v", exp)
		t.Logf("Got: %+v", act)
		t.Fail()
	}
}

func BenchmarkErrorMapAdd(b *testing.B) {
	m := newErrorMap()
	err := errors.New("benchmark")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.add(err)
		}
	})
}

func BenchmarkErrorMapGet(b *testing.B) {
	m := newErrorMap()
	err := errors.New("benchmark")
	m.add(err)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.get(err)
		}
	})
}
