package main

import (
	"flag"
	"log"
	"strings"

	"github.com/valyala/fasthttp"
)

var serverPort = flag.String("port", "8080", "port to use for benchmarks")
var responseSize = flag.Uint("size", 1024, "size of response in bytes")

func main() {
	flag.Parse()
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
