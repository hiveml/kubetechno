// Port assignment logic.
package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"kubetechno/common"
	"kubetechno/common/patch"
	"strconv"
	"sync"
)

type Orchestrator struct {
	superLock  *sync.Mutex // todo: move the locking system out of this struct
	nodeLocks  map[string]*sync.Mutex
	kClient    k8sClient
	lowerBound int
	upperBound int
}

type k8sClient interface {
	GetNodePodsInfo(nodeName string) (portsSet map[int]interface{}, err error)
	AssignPorts(nodeName, poNSName, poName string, ports []int) ([]patch.Patch, error)
	GetPodInfo(poNSName, poName string) (portCount int, err error)
}

func (o *Orchestrator) AssignPorts(nsName, poName, noName string) (*common.Config, []patch.Patch, error) {
	nLock := o.getNodeLock(noName)
	defer nLock.Unlock()
	podPortCount, err := o.kClient.GetPodInfo(nsName, poName)
	if err != nil {
		return nil, nil, errors.New("could not get pod info " + err.Error())
	}
	if podPortCount == 0 {
		return nil, nil, nil
	}
	occupiedPorts, err := o.kClient.GetNodePodsInfo(noName)
	if err != nil {
		return nil, nil, errors.New("could not get node info " + err.Error())
	}
	bytes, _ := json.Marshal(occupiedPorts)
	fmt.Println("occupied ports: " + string(bytes) + ", len: " +
		strconv.Itoa(len(occupiedPorts)) + " when considering " + poName)
	selectedPorts := o.selectPorts(occupiedPorts, podPortCount)
	if selectedPorts == nil {
		return nil, nil, errors.New("not enough ports")
	}
	config := &common.Config{Ports: selectedPorts}
	patches, err := o.kClient.AssignPorts(noName, nsName, poName, selectedPorts)
	if err != nil {
		return nil, nil, errors.New("could not assign ports " + err.Error())
	}
	return config, patches, nil
} // todo: add node lock pruning or a channel system currently it is a minor mem leak.

func (o *Orchestrator) getNodeLock(nodeName string) *sync.Mutex {
	o.superLock.Lock()
	defer o.superLock.Unlock()
	nLock, in := o.nodeLocks[nodeName]
	if !in {
		o.nodeLocks[nodeName] = &sync.Mutex{}
		nLock = o.nodeLocks[nodeName]
	}
	nLock.Lock()
	return nLock
}

func (o *Orchestrator) selectPorts(occupiedPorts map[int]interface{}, count int) []int {
	c := 0
	claimedPorts := make([]int, count)
	if count == 0 {
		return nil
	}
	for i := o.lowerBound; i < o.upperBound; i += 1 {
		if _, in := occupiedPorts[i]; !in {
			claimedPorts[c] = i
			c += 1
			if c == count {
				return claimedPorts
			}
		}
	}
	return nil
}
