package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/api/core/v1"
	"net/http"
)

// Parse the http request input
func parseInput(r *http.Request) (
	uid, nsName, poName, noName, apiVersion string, err error) {
	type BindRequest struct {
		Uid     string     `json:"uid"`
		Binding v1.Binding `json:"object"`
	}
	type BindReviewRequest struct {
		ApiVersion string      `json:"apiVersion"`
		Kind       string      `json:"kind"`
		InnerReq   BindRequest `json:"request"`
	}
	bytes, err := ioutil.ReadAll(r.Body)
	fmt.Println("below is the bind request")
	fmt.Println(string(bytes))
	brr := BindReviewRequest{}
	err = json.Unmarshal(bytes, &brr)
	if err != nil {
		return
	}
	apiVersion = brr.ApiVersion
	bind := brr.InnerReq.Binding
	metadata := bind.ObjectMeta
	noName = bind.Target.Name
	uid = brr.InnerReq.Uid
	poName = metadata.Name
	nsName = metadata.Namespace
	return
}
