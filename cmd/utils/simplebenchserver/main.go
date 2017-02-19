package main

import (
	"log"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/valyala/fasthttp"
)

var serverPort = kingpin.Flag("port", "port to use for benchmarks").
	Default("8080").
	Short('p').
	String()
var responseSize = kingpin.Flag("size", "size of response in bytes").
	Default("1024").
	Short('s').
	Uint()

func main() {
	kingpin.Parse()
	response := strings.Repeat("a", int(*responseSize))
	addr := "localhost:" + *serverPort
	log.Println("Starting HTTP server on:", addr)
	err := fasthttp.ListenAndServe(addr, func(c *fasthttp.RequestCtx) {
		_, werr := c.WriteString(response)
		if werr != nil {
			log.Println(werr)
		}
	})
	if err != nil {
		log.Println(err)
	}
}
