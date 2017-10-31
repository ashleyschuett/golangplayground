package main

import (
	"io"
	"log"
	"net/http"

	"server/lib"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	port := ":8000"
	log.Printf("Webserver has started on port %v \n", port)

	server := http.Server{
		Addr:    port,
		Handler: &myHandler{},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = lib.Hello
	mux["/parsePage"] = lib.ParsePage

	server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "My server: "+r.URL.String())

	if handler, ok := mux[r.URL.String()]; ok {
		handler(w, r)
		return
	}
}
