package main

import (
	"fmt"
	"net/http"
)

type handler struct {
	piString string
}

func (h handler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(writer, h.piString)
}
