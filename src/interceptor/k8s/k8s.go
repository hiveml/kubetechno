// Wraps calls to the k8s api
package k8s

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kubetechno/common/constants"
	"kubetechno/common/patch"
	"strconv"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type Client struct {
	k *kubernetes.Clientset
}

func (c *Client) GetNodePodsInfo(nodeName string) (portsSet map[int]interface{}, err error) {
	pList, err := c.k.CoreV1().Pods("").List(context.Background(),
		metaV1.ListOptions{LabelSelector: constants.KubetechnoDashNode + "=" + nodeName})
	if err != nil {
		return
	}
	fmt.Println(len(pList.Items))
	portsSet = make(map[int]interface{})
	for _, pod := range pList.Items {
		for portIndex := 0; ; portIndex += 1 {
			if portStr, in := pod.Annotations[constants.PortN(portIndex)]; in {
				var port int
				port, err = strconv.Atoi(portStr)
				if err != nil {
					return
				}
				portsSet[port] = nil
			} else {
				break
			}
		}
	}
	return
}

// Determines the number of ports to be assigned to the pod passed as the arg.
func (c *Client) countPodkubetechnoPorts(pod *v1.Pod) int {
	pc := 0
	for _, container := range pod.Spec.Containers {
		pc += int(container.Resources.Requests.Name(constants.KubetechnoSlashPort, resource.DecimalSI).Value())
	}
	return pc
}

func (c *Client) GetPodInfo(poNSName, poName string) (portCount int, err error) {
	pod, err := c.k.CoreV1().Pods(poNSName).Get(context.Background(), poName, metaV1.GetOptions{})
	if err != nil {
		return
	} else if pod == nil {
		err = errors.New("pod was nil")
		return
	}
	labels := pod.ObjectMeta.Labels
	if labels == nil {
		return 0, nil
	} else if ktLabelVal, _ := labels[constants.Kubetechno]; ktLabelVal != "user" {
		return 0, nil
	}
	portCount = c.countPodkubetechnoPorts(pod)
	return
}

// Assigns ports via the annotations and changes the kubetechno label value to 'active'.
func (c *Client) AssignPorts(nodeName, poNSName, poName string, ports []int) (patches []patch.Patch, err error) {
	patches = createPatches(ports, nodeName)
	patchesBytes, err := json.Marshal(patches)
	if err != nil {
		return nil, err
	}
	_, err = c.k.CoreV1().Pods(poNSName).Patch(
		context.Background(), poName, types.JSONPatchType, patchesBytes, metaV1.PatchOptions{})
	return patches, err
}

// Changes the kubetechno label to active and assigns ports.
func createPatches(ports []int, nodeName string) []patch.Patch {
	patchLength := len(ports) + 3
	patchsStartIndex := 3
	patches := make([]patch.Patch, patchLength)
	patches[0] = createAnnotationsPatchComponent(constants.Kubetechno, constants.KubetechnoActiveStatus)
	patches[1] = createAnnotationsPatchComponent(constants.KubetechnoDashNode, nodeName)
	patches[2] = createLabelsPatchComponent(constants.KubetechnoDashNode, nodeName)
	for i, port := range ports {
		patches[i+patchsStartIndex] = createAnnotationsPatchComponent(constants.PortN(i), strconv.Itoa(port))
	}
	return patches
} // todo: integrate container names with the port names

func createAnnotationsPatchComponent(key, value string) patch.Patch {
	return patch.Patch{
		Operation: "add",
		Path:      "/metadata/annotations/" + key,
		Value:     value,
	}
}

func createLabelsPatchComponent(label, value string) patch.Patch {
	return patch.Patch{
		Operation: "add",
		Path:      "/metadata/labels/" + label,
		Value:     value,
	}
}
