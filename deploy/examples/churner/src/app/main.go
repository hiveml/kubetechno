package main

import (
	"encoding/json"
	"fmt"
	"kubetechno-churner/common"
	"net/http"
	"os"
	"strconv"
)

func main() {
	toServe, port, err := config()
	if err != nil {
		fmt.Println("non-int port value")
		os.Exit(1)
	}
	server := http.Server{
		Addr: ":" + port,
	}
	toServeBytes, err := json.MarshalIndent(toServe, "", "  ")
	if err != nil {
		fmt.Println("string to return could not be created " + err.Error())
		os.Exit(1)
	}
	h := handler{string(toServeBytes)}
	server.Handler = h
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func config() (common.PodInfo, string, error) {
	portStr := os.Getenv("PORT0")
	port, err := strconv.Atoi(portStr)
	return common.PodInfo{
		Name:      os.Getenv("pod_name"),
		Namespace: os.Getenv("pod_namespace"),
		Node:      os.Getenv("node_name"),
		Port:      port,
		IP:        os.Getenv("pod_ip"),
	}, portStr, err
}
