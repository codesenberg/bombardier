package main

import (
	"strconv"
	"time"
)

const (
	nilStr = "nil"
)

type nullableUint64 struct {
	val *uint64
}

func (n *nullableUint64) String() string {
	if n.val == nil {
		return nilStr
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
		return nilStr
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

type nullableString struct {
	val *string
}

func (n *nullableString) String() string {
	if n.val == nil {
		return nilStr
	}
	return *n.val
}

func (n *nullableString) Set(value string) error {
	n.val = new(string)
	*n.val = value
	return nil
}
