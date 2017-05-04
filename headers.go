package main

import (
	"fmt"
	"strings"
)

type header struct {
	key, value string
}

type headersList []header

func (h *headersList) String() string {
	return fmt.Sprint(*h)
}

func (h *headersList) IsCumulative() bool {
	return true
}

func (h *headersList) Set(value string) error {
	res := strings.SplitN(value, ":", 2)
	if len(res) != 2 {
		return errInvalidHeaderFormat
	}
	*h = append(*h, header{
		res[0], strings.Trim(res[1], " "),
	})
	return nil
}
