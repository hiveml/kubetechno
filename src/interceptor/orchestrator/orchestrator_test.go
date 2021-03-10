package orchestrator

import (
	"encoding/json"
	"kubetechno/common"
	"reflect"
	"strconv"
	"testing"
)

const testNode1 = "testNode1"
const testPod1 = "testPod1"
const testNS1 = "testNS1"

func TestOrchestrator_AssignPorts(t *testing.T) {
	tkc := NewTestK8sClient(t, []int{9001, 9002, 9004}, 3)
	o := New(9000, 9009, tkc)
	actual, _, err := o.AssignPorts(testNS1, testPod1, testNode1)
	if err != nil {
		t.Log("o.AssignPorts err is not nil and should be. Err msg: " + err.Error())
		t.Fail()
	}
	expected := common.Config{Version: "", Ports: []int{9000, 9003, 9005}}
	if !reflect.DeepEqual(actual.Ports, expected.Ports) {
		expectedBytes, _ := json.Marshal(expected.Ports)
		resultBytes, _ := json.Marshal(actual.Ports)
		t.Log("Assigned ports are not " + string(expectedBytes) + " they are " + string(resultBytes))
		t.Fail()
	}
	checkCalls(t, tkc, 1, 1)
}

func TestOrchestrator_AssignPorts_NotEnough(t *testing.T) {
	tkc := NewTestK8sClient(t, []int{9001, 9002}, 3)
	o := New(9000, 9003, tkc)
	_, _, err := o.AssignPorts(testNS1, testPod1, testNode1)
	if err == nil {
		t.Log("o.AssignPorts err is nil and shouldn't be.")
		t.Fail()
	} else if err.Error() != "not enough ports" {
		t.Log("err msg is '" + err.Error() + "' not 'not enough ports'")
		t.Fail()
	}
	checkCalls(t, tkc, 1, 0)
}

func checkCalls(t *testing.T, tkc *TestK8sClient, expectedCallsToAssignPorts, expectedCallsToGetPodNodePortInfo int) {
	if tkc.callsToGetPodNodePortInfo != expectedCallsToAssignPorts {
		t.Log(strconv.Itoa(tkc.callsToGetPodNodePortInfo) +
			" calls to GetPodNodePortInfo, not " + strconv.Itoa(expectedCallsToAssignPorts))
		t.Fail()
	}
	if tkc.callsToAssignPorts != expectedCallsToGetPodNodePortInfo {
		t.Log(strconv.Itoa(tkc.callsToAssignPorts) +
			" calls to AssignPorts, not " + strconv.Itoa(expectedCallsToGetPodNodePortInfo))
		t.Fail()
	}
}
