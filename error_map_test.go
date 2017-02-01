package main

import (
	"errors"

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
