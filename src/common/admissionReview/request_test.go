package admissionReview

import (
	"encoding/json"
	"kubetechno/mutator/patchers"
	"kubetechno/mutator/patchers/consul"
	"kubetechno/mutator/patchers/core"
	"testing"

	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const TestApiVersion = "test-api-version"
const TestKind = "test-kind"

type ReviewRequestTestCase struct {
	Message  string
	Req      ReviewRequest
	Expected string
	T        *testing.T
}

func (rrtc ReviewRequestTestCase) runTest() {
	_, resp, _ := rrtc.Req.CreateReviewResponse([]patchers.Patcher{
		core.Patcher{}, consul.NewPatcher("test-patcher")})
	actualBytes, _ := json.Marshal(resp)
	actual := string(actualBytes)
	if rrtc.Expected != actual {
		rrtc.T.Fail()
		rrtc.T.Log("expected does not equal actual (respectively below):")
		rrtc.T.Log(rrtc.Expected)
		rrtc.T.Log(actual)
	}
}

func TestReviewRequest_CreateReviewResponse_WithoutkubetechnoAnnotations(t *testing.T) {
	ReviewRequestTestCase{
		Req: ReviewRequest{
			ApiVersion: TestApiVersion,
			Kind:       TestKind,
			InnerReq: Request{
				Uid: "test-uid",
				Object: v1.Pod{
					ObjectMeta: metaV1.ObjectMeta{
						Annotations: createAnnotations("foo", "bar"),
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{{}},
					},
				},
			},
		},
		Expected: "{\"apiVersion\":\"test-api-version\",\"kind\":\"test-kind\",\"response\":{\"uid\":\"test-uid\"" +
			",\"allowed\":true,\"patchType\":\"JSONPatch\",\"status\":{\"code\":200,\"message\":\"Non-kubetechno Pod\"}}}",
		T: t,
	}.runTest()
}

func TestReviewRequest_CreateReviewResponse_WithkubetechnoAnnotations(t *testing.T) {
	ReviewRequestTestCase{
		Req: reviewRequestGenerator("kubetechno_port_count", "1", true),
		Expected: "{\"apiVersion\":\"test-api-version\",\"kind\":\"test-kind\",\"response\":{\"uid\":\"test-uid\"," +
			"\"allowed\":true,\"patchType\":\"JSONPatch\",\"patch\":\"W3sib3AiOiJyZXBsYWNlIiwicGF0aCI6Ii9zcGVjIiwidm" +
			"FsdWUiOnsiY29udGFpbmVycyI6W3sibmFtZSI6IiIsImVudiI6W3sibmFtZSI6IlBPUlQwIiwidmFsdWVGcm9tIjp7ImZpZWxkUmVmIj" +
			"p7ImZpZWxkUGF0aCI6Im1ldGFkYXRhLmFubm90YXRpb25zWydQT1JUMCddIn19fV0sInJlc291cmNlcyI6eyJsaW1pdHMiOnsia3ViZX" +
			"RlY2huby9wb3J0IjoiMSJ9LCJyZXF1ZXN0cyI6eyJrdWJldGVjaG5vL3BvcnQiOiIxIn19fV0sImhvc3ROZXR3b3JrIjp0cnVlfX0sey" +
			"JvcCI6InJlcGxhY2UiLCJwYXRoIjoiL21ldGFkYXRhIiwidmFsdWUiOnsiY3JlYXRpb25UaW1lc3RhbXAiOm51bGwsImFubm90YXRpb2" +
			"5zIjp7Imt1YmV0ZWNobm9fcG9ydF9jb3VudCI6IjEifX19XQ==\"," +
			"\"status\":{\"code\":200,\"message\":\"kubetechno pod with alterations\"}}}",
		T: t,
	}.runTest()
}
func TestReviewRequest_CreateReviewResponse_WithBadkubetechnoAnnotations(t *testing.T) {
	ReviewRequestTestCase{
		Req: reviewRequestGenerator("kubetechno_port_count", "z", true),
		Expected: "{\"apiVersion\":\"test-api-version\",\"kind\":\"test-kind\",\"response\":{\"uid\":\"test-uid\"," +
			"\"allowed\":false,\"patchType\":\"JSONPatch\",\"status\":{\"code\":403,\"message\":\"non-int value for " +
			"kubetechno port count\"}}}",
		T: t,
	}.runTest()
}

func reviewRequestGenerator(key, val string, addAnnotations bool) ReviewRequest {
	var annotations map[string]string = nil
	if addAnnotations {
		annotations = createAnnotations(key, val)
	}
	return ReviewRequest{
		ApiVersion: TestApiVersion,
		Kind:       TestKind,
		InnerReq: Request{
			Uid: "test-uid",
			Object: v1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Annotations: annotations,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{}},
				},
			},
		},
	}
}

func createAnnotations(key, val string) map[string]string {
	annotations := make(map[string]string, 1)
	annotations[key] = val
	return annotations
}
