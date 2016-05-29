package main

import (
	"fmt"
)

type units struct {
	scale uint64
	base  string
	units []string
}

var (
	binaryUnits = &units{
		scale: 1024,
		base:  "",
		units: []string{"KB", "MB", "GB", "TB", "PB"},
	}
	timeUnitsUs = &units{
		scale: 1000,
		base:  "us",
		units: []string{"ms", "s"},
	}
	timeUnitsS = &units{
		scale: 60,
		base:  "s",
		units: []string{"m", "h"},
	}
)

func formatUnits(n float64, m *units, prec int) string {
	amt := n
	unit := m.base

	scale := float64(m.scale) * 0.85

	for i := 0; i+1 < len(m.units) && amt >= scale; i++ {
		amt /= float64(m.scale)
		unit = m.units[i]
	}
	return fmt.Sprintf("%.*f%s", prec, amt, unit)
}

func formatBinary(n float64) string {
	return formatUnits(n, binaryUnits, 2)
}

func formatTimeUs(n float64) string {
	units := timeUnitsUs
	if n >= 1000000.0 {
		n /= 1000000.0
		units = timeUnitsS
	}
	return formatUnits(n, units, 2)
}
