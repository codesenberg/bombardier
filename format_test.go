package main

import (
	"testing"
)

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024

	K = 1000
	M = K * 1000
)

func TestShouldFormatBinary(t *testing.T) {
	expectations := []struct {
		in  float64
		out string
	}{
		{10.0, "10.00"},
		{10.001, "10.00"},
		{1.0 * KB, "1.00KB"},
		{1.2 * KB, "1.20KB"},
		{1.202 * KB, "1.20KB"},
		{5 * KB, "5.00KB"},
		{1.0 * MB, "1.00MB"},
		{1.3 * MB, "1.30MB"},
		{1.302 * MB, "1.30MB"},
		{6 * MB, "6.00MB"},
		{1.0 * GB, "1.00GB"},
		{1.4 * GB, "1.40GB"},
		{1.402 * GB, "1.40GB"},
		{7 * GB, "7.00GB"},
	}
	for _, e := range expectations {
		actual := formatBinary(e.in)
		expected := e.out
		if expected != actual {
			t.Errorf("Expected \"%v\", but got \"%v\"", expected, actual)
		}
	}
}

func TestShouldFormatUs(t *testing.T) {
	expectations := []struct {
		in  float64
		out string
	}{
		{20, "20.00us"},
		{22.222, "22.22us"},
		{20 * K, "20.00ms"},
		{20 * M, "20.00s"},
		{60 * M, "1.00m"},
		{10 * 60 * M, "10.00m"},
		{90 * 60 * M, "1.50h"},
	}
	for _, e := range expectations {
		actual := formatTimeUs(e.in)
		expected := e.out
		if expected != actual {
			t.Errorf("Expected \"%v\", but got \"%v\"", expected, actual)
		}
	}
}
