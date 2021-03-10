// Server that intercepts pod bind requests in order to add port assigment information.
package main

import (
	"fmt"
	"kubetechno/interceptor/config"
	"kubetechno/interceptor/handler"
	"kubetechno/interceptor/k8s"
	"kubetechno/interceptor/orchestrator"
	"net/http"
	"os"
)

func main() {
	settings, err := config.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	kClient, err := k8s.NewClient()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	o := orchestrator.New(settings.LowerBound, settings.UpperBound, kClient)
	bindHandler := handler.New(o)
	http.HandleFunc("/bind-hook", bindHandler.ServeHTTP)
	if err := http.ListenAndServeTLS(":443", settings.CertFilePath, settings.KeyFilePath, nil); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
