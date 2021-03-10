package orchestrator

import (
	"kubetechno/common/patch"
	"testing"
)

// Instantiates, sets up, and returns a TestK8sClient instance.
func NewTestK8sClient(t *testing.T, occupiedPorts []int, podPortCount int) *TestK8sClient {
	tkc := &TestK8sClient{
		t:            t,
		occupied:     make(map[int]interface{}),
		podPortCount: podPortCount,
	}
	for _, port := range occupiedPorts {
		tkc.occupied[port] = nil
	}
	return tkc
}

// Testing implementation of k8sClient
type TestK8sClient struct {
	t                         *testing.T
	occupied                  map[int]interface{}
	podPortCount              int
	callsToAssignPorts        int
	callsToGetPodNodePortInfo int
}

func (c *TestK8sClient) GetNodePodsInfo(nodeName string) (portsSet map[int]interface{}, err error) {
	if nodeName != testNode1 {
		c.t.Log("GetPodNodePortInfo nodeName: " + nodeName + " != " + testNode1)
		c.t.Fail()
	}
	c.callsToGetPodNodePortInfo += 1
	return c.occupied, nil

}
func (c *TestK8sClient) AssignPorts(nodeName, poNSName, poName string, ports []int) ([]patch.Patch, error) {
	if poNSName != testNS1 {
		c.t.Log("GetPodNodePortInfo podNSName: " + poNSName + " != " + testNS1)
		c.t.Fail()
	}
	if poName != testPod1 {
		c.t.Log("GetPodNodePortInfo podName: " + poName + " != " + testPod1)
		c.t.Fail()
	}
	c.callsToAssignPorts += 1
	return nil, nil
}

func (c *TestK8sClient) GetPodInfo(poNSName, poName string) (portCount int, err error) {
	return c.podPortCount, nil
}
