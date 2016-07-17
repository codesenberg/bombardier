package main

import (
	"math"
	"math/big"
	"strconv"
	"testing"
	"time"
)

func TestNullableUint64ConversionToString(t *testing.T) {
	nilint := &nullableUint64{val: nil}
	if s := nilint.String(); s != "nil" {
		t.Errorf("Expected \"nil\", but got %v", s)
	}
	v := uint64(42)
	nonnilint := &nullableUint64{val: &v}
	if s, e := nonnilint.String(), strconv.FormatUint(v, 10); s != e {
		t.Errorf("Expected %v, but got %v", e, s)
	}
}

func TestNullableUint64Parsing(t *testing.T) {
	n := &nullableUint64{}
	if err := n.Set("-1"); err == nil {
		t.Error("Should fail on negative values")
	}
	if err := n.Set(""); err == nil {
		t.Error("Should fail on empty string")
	}
	b := big.NewInt(0)
	b.SetUint64(math.MaxUint64)
	b.Add(b, big.NewInt(1))
	if err := n.Set(b.String()); err == nil {
		t.Error("Should fail on large values")
	}
	max := strconv.FormatUint(math.MaxUint64, 10)
	if err := n.Set(max); err != nil || *n.val != uint64(18446744073709551615) {
		t.Error("Shouldn't fail on max value")
	}
}

func TestNullableDurationConversionToString(t *testing.T) {
	nildur := &nullableDuration{val: nil}
	if s := nildur.String(); s != "nil" {
		t.Errorf("Expected \"nil\", but got %v", s)
	}
	d := time.Second
	nonnildir := &nullableDuration{val: &d}
	if s := nonnildir.String(); s != "1s" {
		t.Errorf("Expected 1s, but got %v", s)
	}
}

func TestNullableDurationParsing(t *testing.T) {
	d := &nullableDuration{}
	if err := d.Set(""); err == nil {
		t.Error("Should fail on empty string")
	}
	if err := d.Set("Wubba lubba dub dub!"); err == nil {
		t.Error("Should fail on incorrect values")
	}
	if err := d.Set("1s"); err != nil || *d.val != time.Second {
		t.Error("Shouldn't fail on correct values")
	}
}
