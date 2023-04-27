package main

import (
	"bytes"
	"log"
	"net/http"

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
var stdHTTP = kingpin.Flag("std-http", "use standard http library").
	Default("false").
	Bool()

func main() {
	kingpin.Parse()
	response := bytes.Repeat([]byte("a"), int(*responseSize))
	addr := "localhost:" + *serverPort
	log.Println("Starting HTTP server on:", addr)
	var lserr error
	if *stdHTTP {
		lserr = http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, werr := w.Write(response)
			if werr != nil {
				log.Println(werr)
			}
		}))
	} else {
		lserr = fasthttp.ListenAndServe(addr, func(c *fasthttp.RequestCtx) {
			_, werr := c.Write(response)
			if werr != nil {
				log.Println(werr)
			}
		})
	}
	if lserr != nil {
		log.Println(lserr)
	}
}
