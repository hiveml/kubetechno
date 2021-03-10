package admissionReview

import (
	"encoding/base64"
	"encoding/json"
	"k8s.io/api/core/v1"
	"kubetechno/common/patch"
	"kubetechno/mutator/patchers"
)

type ReviewRequest struct {
	ApiVersion string  `json:"apiVersion"`
	Kind       string  `json:"kind"`
	InnerReq   Request `json:"request"`
}

type Request struct {
	Uid    string `json:"uid"`
	Object v1.Pod `json:"object"`
}

func (req ReviewRequest) CreateReviewResponse(
	patchers []patchers.Patcher) (patches []patch.Patch, resp ReviewResponse, err error) {
	resp.ApiVersion = req.ApiVersion
	resp.Kind = req.Kind
	resp.InnerResp = Response{
		Uid:       req.InnerReq.Uid,
		Allowed:   false,
		PatchType: "JSONPatch",
		Patch:     "",
		Status: ResponseStatus{
			Code:    200,
			Message: "Non-kubetechno Pod",
		},
	}
	pod := req.InnerReq.Object
	podPossiblyAltered := false
	for _, patcher := range patchers {
		changed := false
		if changed, err = patcher.UpdatePod(&pod); err != nil {
			resp.InnerResp.Status = ResponseStatus{
				Code:    403,
				Message: err.Error(),
			}
			return
		} else if changed {
			podPossiblyAltered = true
		}
	}
	if !podPossiblyAltered {
		resp.InnerResp.Allowed = true
		return
	}
	patches = []patch.Patch{
		{
			Operation: "replace",
			Path:      "/spec",
			Value:     pod.Spec,
		},
		{
			Operation: "replace",
			Path:      "/metadata",
			Value:     pod.ObjectMeta,
		},
	}
	if patchList, err := json.Marshal(patches); err == nil {
		resp.InnerResp.Patch = base64.StdEncoding.EncodeToString(patchList)
		resp.InnerResp.Allowed = true
		resp.InnerResp.Status.Message = "kubetechno pod with alterations"
	}
	return
}
