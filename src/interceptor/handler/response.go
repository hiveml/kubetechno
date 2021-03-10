package handler

// Returned to the api server.
// Implements https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#response
type Response struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Response   AllowRsp `json:"response"`
}

type AllowRsp struct {
	UID     string `json:"uid"`
	Allowed bool   `json:"allowed"`
}
