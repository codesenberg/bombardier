package main

import (
	"strconv"
	"time"
)

type nullableUint64 struct {
	val *uint64
}

func (n *nullableUint64) String() string {
	if n.val == nil {
		return "nil"
	}
	return strconv.FormatUint(*n.val, 10)
}

func (n *nullableUint64) Set(value string) error {
	res, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	n.val = new(uint64)
	*n.val = res
	return nil
}

type nullableDuration struct {
	val *time.Duration
}

func (n *nullableDuration) String() string {
	if n.val == nil {
		return "nil"
	}
	return n.val.String()
}

func (n *nullableDuration) Set(value string) error {
	res, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	n.val = &res
	return nil
}
