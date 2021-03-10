package core

import (
	"errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"strconv"
)

type Patcher struct {
}

func (p Patcher) UpdatePod(po *v1.Pod) (bool, error) {
	portCountStr, in := po.ObjectMeta.Annotations["kubetechno_port_count"]
	if !in {
		return false, nil
	}
	portCount, err := strconv.Atoi(portCountStr)
	if err != nil {
		return false, errors.New("non-int value for kubetechno port count")
	}
	po.Spec.HostNetwork = true // todo: revisit making this required instead of added
	addPortEnvs(po, portCount)
	kubetechnoPortResources(po, portCount)
	return true, nil
}

func addPortEnvs(po *v1.Pod, portCount int) {
	containers := po.Spec.Containers
	initContainers := po.Spec.InitContainers
	envVarList := []v1.EnvVar{}
	for portIndex := 0; portIndex < portCount; portIndex += 1 {
		portIndexName := "PORT" + strconv.Itoa(portIndex)
		envVarList = append(envVarList, v1.EnvVar{
			Name: portIndexName,
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "metadata.annotations['" + portIndexName + "']",
				},
			},
		})
	}
	for i := range containers {
		containers[i].Env = append(containers[i].Env, envVarList...)
	}
	for i := range initContainers {
		initContainers[i].Env = append(initContainers[i].Env, envVarList...)
	}
}

func kubetechnoPortResources(pod *v1.Pod, portCount int) {
	containers := pod.Spec.Containers
	c0 := containers[0]
	if c0.Resources.Limits == nil {
		c0.Resources.Limits = make(v1.ResourceList, 1)
	}
	if c0.Resources.Requests == nil {
		c0.Resources.Requests = make(v1.ResourceList, 1)
	}
	c0.Resources.Limits["kubetechno/port"] = *resource.NewQuantity(int64(portCount), "DecimalExponent")
	c0.Resources.Requests["kubetechno/port"] = *resource.NewQuantity(int64(portCount), "DecimalExponent")
	containers[0] = c0
}
