package handler

import (
	"encoding/json"
	"fmt"
	"kubetechno/common/admissionReview"
	"kubetechno/mutator/patchers"
	"kubetechno/mutator/server/log"
	"net/http"
)

func New(p []patchers.Patcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ls := log.NewLogEntry()
		defer func() {
			if bytes, err := json.Marshal(ls); err != nil {
				logString := "{ \"log-print-err\": " + err.Error() + ", " + "\"log-string\":" + "}"
				fmt.Println(logString)
			} else {
				fmt.Println(string(bytes))
			}
		}()
		rReq := admissionReview.ReviewRequest{}
		ls.LastStep = "decode input"
		if err := json.NewDecoder(r.Body).Decode(&rReq); err != nil {
			fmt.Println("decode error")
			w.WriteHeader(400)
			ls.Err = err.Error()
			return
		}
		ls.LastStep = "create response"
		patches, rRsp, reviewErr := rReq.CreateReviewResponse(p)
		if reviewErr != nil {
			ls.Err = reviewErr.Error()
		}
		ls.Patches = patches
		ls.LastStep = "marshal response"
		rspBytesToPrint, err := json.Marshal(rRsp)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(200)
			ls.Err = err.Error()
			errRspBytes, errRspMarshalErr := json.Marshal(
				admissionReview.ReviewResponse{
					ApiVersion: rReq.ApiVersion, Kind: rReq.Kind,
					InnerResp: admissionReview.Response{
						Uid:     rReq.InnerReq.Uid,
						Allowed: false,
						Status: admissionReview.ResponseStatus{
							Code:    403,
							Message: err.Error(),
						},
					},
				})
			if errRspMarshalErr == nil {
				w.Write(errRspBytes)
			}
			return
		}
		w.WriteHeader(200)
		w.Write(rspBytesToPrint)
	}
}
