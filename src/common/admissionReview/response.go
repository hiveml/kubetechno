package admissionReview

type ReviewResponse struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	InnerResp  Response `json:"response"`
}

type Response struct {
	Uid       string         `json:"uid"`
	Allowed   bool           `json:"allowed"`
	Warnings  []string       `json:"warnings,omitempty"`
	PatchType string         `json:"patchType"`
	Patch     string         `json:"patch,omitempty"`
	Status    ResponseStatus `json:"status,omitempty"`
}

type ResponseStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
