package consul

import (
	"encoding/json"
	"errors"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubetechno/common/constants"
	"strconv"
	"testing"
)

func TestUpdatePod(t *testing.T) {
	annotations := createConsulAnnotations()
	input := CreateInputPod(annotations)
	NewTestCase(t, "Update Pod General", input, CreateOutputPod(annotations), true).Run()
}

func TestUpdatePod_NoLifeCycle(t *testing.T) {
	annotations := createConsulAnnotations()
	input := CreateInputPod(annotations)
	container := input.Spec.Containers[0]
	container.Lifecycle = nil
	input.Spec.Containers[0] = container
	NewTestCase(t, "Update Pod No Life Cycle", input, CreateOutputPod(annotations), true).Run()
}

func TestUpdatePod_NoName(t *testing.T) {
	annoataions := make(map[string]string, 1)
	annoataions["foo"] = "bar"
	input := CreateInputPod(annoataions)
	NewTestCase(t, "Update Pod No Name", input, input, false).Run()
}

func TestUpdatePod_BadPullPolicy(t *testing.T) {
	annotations := createConsulAnnotations()
	annotations[constants.ConsulClientImagePullPolicy] = "badPullPolicy"
	input := CreateInputPod(annotations)
	NewErrTestCase(
		t, "Update Pod Bad Pull Policy", input, errors.New("bad pull policy")).Run()
}

func TestUpdatePod_BadBufferSeconds(t *testing.T) {
	annotations := createConsulAnnotations()
	annotations[constants.ConsulBufferSecs] = "badBufferSeconds"
	input := CreateInputPod(annotations)
	NewErrTestCase(
		t, "Update Pod Bad Buffer Seconds", input, errors.New(
			"non int annotation val for kubetechno_consul_buffer_seconds")).Run()
}

func TestUpdatePod_NoConsulCheckPath(t *testing.T) {
	annotations := createConsulAnnotations()
	delete(annotations, constants.ConsulCheckPath)
	input := CreateInputPod(annotations)
	NewErrTestCase(t, "No Consul Check Path", input, errors.New(
		"could not get consul check path")).Run()
}

func TestUpdatePod_BadConsulServiceName(t *testing.T) {
	annotations := createConsulAnnotations()
	annotations[constants.ConsulServiceName] = ""
	input := CreateInputPod(annotations)
	NewErrTestCase(t, "Update Pod Bad Consul Service Name", input, errors.New("bad consul service name")).Run()
}

func TestUpdatePod_NoAnnotations(t *testing.T) {
	input := CreateInputPod(nil)
	NewTestCase(t, "Update Pod No Annotations", input, input, false).Run()
}

// types and non-test functions below

func createConsulAnnotations() map[string]string {
	annotations := make(map[string]string, 1)
	annotations[constants.ConsulServiceName] = "test-service"
	annotations[constants.ConsulBufferSecs] = "11"
	annotations[constants.ConsulTimeoutSeconds] = "11"
	annotations[constants.ConsulTimeoutSeconds] = "11"
	annotations[constants.ConsulInitialDelaySeconds] = "11"
	annotations[constants.ConsulCheckPath] = "/test.html"
	annotations[constants.ConsulConsulInitImage] = "test-kubetechno-init-image"
	return annotations
}

func CreateInputPod(annotations map[string]string) v1.Pod {
	return v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:        "test-pod",
			Namespace:   "test-namespace",
			Annotations: annotations,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "test-container",
					Image: "alpine",
				},
				{
					Name:  "sidecar",
					Image: "alpine",
				},
			},
		},
	}
}

// assumes a lot about the input pod
func CreateOutputPod(annotations map[string]string) v1.Pod {
	po := CreateInputPod(annotations)
	po.Spec.Volumes = []v1.Volume{{
		Name: "kubetechno-consul",
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{},
		},
	}}
	po.Spec.InitContainers = []v1.Container{
		{
			Name:  "kubetechno-init-consul",
			Image: "test-kubetechno-init-image",
			VolumeMounts: []v1.VolumeMount{{
				Name:      "kubetechno-consul",
				MountPath: "/mvTarget",
			}},
			ImagePullPolicy: v1.PullIfNotPresent,
		},
	}
	po.Spec.Containers[0].ReadinessProbe = &v1.Probe{
		InitialDelaySeconds: 11,
		TimeoutSeconds:      11,
		SuccessThreshold:    1,
		FailureThreshold:    1,
		Handler: v1.Handler{
			Exec: &v1.ExecAction{
				Command: []string{"/kubetechnoConsul/client", "check"},
			},
		},
	}
	po.Spec.Containers[0].Env = []v1.EnvVar{
		{
			Name:  constants.ConsulCheckPath,
			Value: "/test.html",
		}, {
			Name:  constants.ConsulServiceName,
			Value: "test-service",
		}, {
			Name:  constants.ConsulBufferSecs,
			Value: "11",
		}, {
			Name:  constants.ConsulTimeoutSeconds,
			Value: "11",
		}, {
			Name:  constants.ConsulPeriodSeconds,
			Value: "150",
		},
		{
			Name:  constants.ConsulInitialDelaySeconds,
			Value: "11",
		},
	}
	po.Spec.Containers[0].VolumeMounts = []v1.VolumeMount{
		{
			Name:      "kubetechno-consul",
			MountPath: "/kubetechnoConsul",
		},
	}
	po.Spec.Containers[0].ReadinessProbe = &v1.Probe{
		Handler: v1.Handler{
			Exec: &v1.ExecAction{
				Command: []string{"/kubetechnoConsul/client", "check"},
			},
		},
		InitialDelaySeconds: 11,
		TimeoutSeconds:      11,
		PeriodSeconds:       150,
		SuccessThreshold:    1,
		FailureThreshold:    1,
	}
	po.Spec.Containers[0].Lifecycle = &v1.Lifecycle{
		PreStop: &v1.Handler{
			Exec: &v1.ExecAction{
				Command: []string{constants.ConsulClientPath, "dereg"},
			},
		},
	}
	return po
}

func NewTestCase(t *testing.T, message string, input, expected v1.Pod, changed bool) TestCase {
	return TestCase{
		Message:  message,
		Input:    input,
		Expected: expected,
		T:        t,
		Changed:  changed,
	}
}

func NewErrTestCase(t *testing.T, message string, input v1.Pod, err error) TestCase {
	return TestCase{
		Message:     message,
		Input:       input,
		Expected:    v1.Pod{}, // not checked
		ExpectedErr: err,
		T:           t,
	}
}

type TestCase struct {
	Message     string
	Input       v1.Pod
	Expected    v1.Pod
	ExpectedErr error
	Changed     bool
	T           *testing.T
}

func (tc TestCase) Run() {
	tc.T.Log(tc.Message)
	underTest := NewPatcher("test-kubetechno-init-image")
	changed, err := underTest.UpdatePod(&tc.Input)
	if tc.ExpectedErr == nil && changed != tc.Changed {
		tc.T.Fail()
		tc.T.Log("unexpected changed return value of " + strconv.FormatBool(changed))
		return
	}
	if tc.ExpectedErr != nil {
		if err == nil {
			tc.T.Fail()
			tc.T.Log("no error returned as expected")
		} else if err.Error() != tc.ExpectedErr.Error() {
			tc.T.Fail()
			tc.T.Log("error is incorrect: '" + err.Error() + "'")
		}
		return
	}
	if err != nil && err != tc.ExpectedErr {
		tc.T.Fail()
		tc.T.Log("error is incorrect")
		return
	}
	actualBytes, _ := json.Marshal(tc.Input)
	actualString := string(actualBytes)
	expectedBytes, _ := json.Marshal(tc.Expected)
	expectedString := string(expectedBytes)
	if actualString != expectedString {
		tc.T.Fail()
		tc.T.Log("based on their JSON marshaled strings, actual (above) != expected (below)")
		tc.T.Log(actualString)
		tc.T.Log(expectedString)
	}
}
