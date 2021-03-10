package main

import (
	"fmt"
	"kubetechno/mutator/patchers"
	"kubetechno/mutator/patchers/consul"
	"kubetechno/mutator/patchers/core"
	"kubetechno/mutator/server/handler"
	"net/http"
	"os"
)

func main() {
	settings, err := New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	http.HandleFunc("/pod-mutate-hook", handler.New([]patchers.Patcher{
		consul.NewPatcher(settings.DefaultConsulImage), core.Patcher{}}))
	if err := http.ListenAndServeTLS(":443",
		settings.CertFilePath, settings.KeyFilePath, nil); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
