package core

import (
	"encoding/json"
	"errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"testing"
)

func TestUpdatePod_BadPortCount(t *testing.T) {
	in := CreateBaseInPod()
	in.ObjectMeta.Annotations["kubetechno_port_count"] = "fail"
	TestCase{
		input:                   in,
		expected:                v1.Pod{},
		t:                       t,
		expectedErr:             errors.New("non-int value for kubetechno port count"),
		expectedPossiblyChanged: false,
	}.Run()
}

func TestUpdatePod_NoPortCount(t *testing.T) {
	in := CreateBaseInPod()
	delete(in.ObjectMeta.Annotations, "kubetechno_port_count")
	expected := CreateBaseInPod()
	delete(expected.ObjectMeta.Annotations, "kubetechno_port_count")
	TestCase{
		input:                   in,
		expected:                expected,
		t:                       t,
		expectedPossiblyChanged: false,
	}.Run()
}

func TestUpdatePod_NoLabels(t *testing.T) {
	in := CreateBaseInPod()
	in.ObjectMeta.Labels = nil
	expected := CreateBaseOutPod()
	expected.ObjectMeta.Labels = make(map[string]string, 1)
	TestCase{
		input:                   in,
		expected:                expected,
		t:                       t,
		expectedPossiblyChanged: true,
	}.Run()
}

func TestUpdatePod(t *testing.T) {
	TestCase{
		input:                   CreateBaseInPod(),
		expected:                CreateBaseOutPod(),
		t:                       t,
		expectedPossiblyChanged: true,
	}.Run()
}

func CreateBaseOutPod() v1.Pod {
	po := CreateBaseInPod()
	po.Labels["kubetechno"] = "user"
	for i, c := range po.Spec.InitContainers {
		c.Env = append(c.Env, []v1.EnvVar{
			{
				Name: "PORT0",
				ValueFrom: &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{
						FieldPath: "metadata.annotations['PORT0']",
					},
				},
			}, {
				Name: "PORT1",
				ValueFrom: &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{
						FieldPath: "metadata.annotations['PORT1']",
					},
				},
			},
		}...)
		po.Spec.InitContainers[i] = c
	}
	for i, c := range po.Spec.Containers {
		c.Env = append(c.Env, []v1.EnvVar{
			{
				Name: "PORT0",
				ValueFrom: &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{
						FieldPath: "metadata.annotations['PORT0']",
					},
				},
			}, {
				Name: "PORT1",
				ValueFrom: &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{
						FieldPath: "metadata.annotations['PORT1']",
					},
				},
			},
		}...)
		po.Spec.Containers[i] = c
	}
	c0 := po.Spec.Containers[0]
	kubetechnoOnlyResourceList := make(v1.ResourceList, 1)
	kubetechnoOnlyResourceList["kubetechno/port"] = *resource.NewQuantity(int64(2), "DecimalExponent")
	c0.Resources = v1.ResourceRequirements{
		Limits:   kubetechnoOnlyResourceList,
		Requests: kubetechnoOnlyResourceList,
	}
	po.Spec.Containers[0] = c0
	return po
}

func CreateBaseInPod() v1.Pod {
	return v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:        "test-pod",
			Namespace:   "test-namespace",
			Annotations: CreateBaseAnnotations(),
			Labels:      CreateBaseLabels(),
		},
		Spec: v1.PodSpec{
			HostNetwork: true,
			InitContainers: []v1.Container{
				{
					Name:  "init-c0",
					Image: "alpine",
				},
				{
					Name:  "init-c1",
					Image: "alpine",
				},
			},
			Containers: []v1.Container{
				{
					Name:  "c0",
					Image: "alpine",
					Env: []v1.EnvVar{
						{
							Name:  "foo",
							Value: "bar",
						}, {
							Name:  "bar",
							Value: "foo",
						},
					},
				}, {
					Name:  "c1",
					Image: "alpine",
					Resources: v1.ResourceRequirements{
						Limits:   CreateBaseResourceMap(),
						Requests: CreateBaseResourceMap(),
					},
				},
				{
					Name:  "c2",
					Image: "alpine",
					Resources: v1.ResourceRequirements{
						Requests: CreateBaseResourceMap(),
					},
				},
			},
		},
	}
}

func CreateBaseLabels() map[string]string {
	labels := make(map[string]string, 3)
	labels["a"] = "A"
	labels["b"] = "B"
	labels["c"] = "C"
	labels["kubetechno"] = "user"
	return labels
}

func CreateBaseAnnotations() map[string]string {
	annotations := make(map[string]string, 2)
	annotations["kubetechno_port_count"] = "2"
	annotations["foo"] = "bar"
	return annotations
}

func CreateBaseResourceMap() v1.ResourceList {
	rl := make(v1.ResourceList, 2)
	rl["foo"] = *resource.NewQuantity(int64(1), "DecimalExponent")
	rl["bar"] = *resource.NewQuantity(int64(2), "DecimalExponent")
	return rl
}

type TestCase struct {
	input                   v1.Pod
	expected                v1.Pod
	expectedErr             error
	t                       *testing.T
	expectedPossiblyChanged bool
}

func (tc TestCase) Run() {
	underTest := Patcher{}
	possiblyChanged, err := underTest.UpdatePod(&tc.input)
	if tc.expectedErr != nil {
		if err == nil {
			tc.t.Fail()
			tc.t.Log("no error returned as expected")
		} else if err.Error() != tc.expectedErr.Error() {
			tc.t.Fail()
			tc.t.Log("error is incorrect: '" + err.Error() + "'")
		}
		return
	}
	if err != nil && tc.expectedErr == nil {
		tc.t.Fail()
		tc.t.Log("there was an error " +
			err.Error() + " and there shouldn't be")
		return
	}
	if tc.expectedPossiblyChanged != possiblyChanged {
		tc.t.Fail()
		tc.t.Log("possibly changed does not equal what the expected val of " +
			strconv.FormatBool(tc.expectedPossiblyChanged))
		return
	}
	actualBytes, _ := json.Marshal(tc.input)
	actualString := string(actualBytes)
	expectedBytes, _ := json.Marshal(tc.expected)
	expectedString := string(expectedBytes)
	if actualString != expectedString {
		tc.t.Fail()
		tc.t.Log("based on their JSON marshaled strings, actual (above) != expected (below)")
		tc.t.Log(actualString)
		tc.t.Log(expectedString)
	}
}
