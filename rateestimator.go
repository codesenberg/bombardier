package main

import (
	"math/big"
	"time"
)

const (
	panicZeroRate         = "rate can't be zero"
	panicNegativeAdjustTo = "adjustTo can't be negative or zero"
)

func estimate(rate uint64, adjustTo time.Duration) (time.Duration, uint64) {
	if rate == 0 {
		panic(panicZeroRate)
	}
	if adjustTo <= 0 {
		panic(panicNegativeAdjustTo)
	}
	br := new(big.Int).SetUint64(rate)
	bd := new(big.Int).SetInt64(oneSecond.Nanoseconds())
	gcd := new(big.Int).GCD(nil, nil, br, bd).Uint64()
	nr, nd := rate/gcd, uint64(oneSecond.Nanoseconds())/gcd
	adjustInt := uint64(adjustTo.Nanoseconds())
	if nd >= adjustInt {
		return time.Duration(nd), nr
	}
	coef := adjustInt / nd
	return time.Duration(coef * nd), coef * nr
}
