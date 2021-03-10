package main

import (
	"errors"
	"k8s.io/api/core/v1"
	"kubetechno-churner/churner/consul"
	"kubetechno-churner/common"
	"strconv"
	"strings"
	"sync"
)

func compareK8sPodsToConsulServices(cpis []consul.AppInstanceInfo, pods []v1.Pod) error {
	if len(pods) != len(cpis) {
		return errors.New("length of pods " + strconv.Itoa(len(pods)) +
			" does not equal " + strconv.Itoa(len(cpis)))
	}
	length := len(pods)
	addressesFromConsul := make(map[string]interface{}, length)
	for _, po := range pods {
		address := po.Status.PodIP + ":" + po.ObjectMeta.Annotations["PORT0"]
		addressesFromConsul[address] = true
	}
	addressesFromK8s := make(map[string]interface{}, length)
	for _, cpi := range cpis {
		address := cpi.IP + ":" + strconv.Itoa(cpi.Port)
		addressesFromK8s[address] = true
	}

	errStrList := []string{}
	for addressFromConsul := range addressesFromConsul {
		if _, in := addressesFromK8s[addressFromConsul]; !in {
			errStrList = append(errStrList, addressFromConsul+" not in k8s")
		}
	}
	for addressFromK8s := range addressesFromK8s {
		if _, in := addressesFromConsul[addressFromK8s]; !in {
			errStrList = append(errStrList, addressFromK8s+" not in consul")
		}
	}
	if len(errStrList) != 0 {
		return errors.New(strings.Join(errStrList, ", "))
	}
	return nil
}

func compareAppDataToConsulServices(cpis []consul.AppInstanceInfo) error {
	wg := sync.WaitGroup{}
	errAccMutex := sync.Mutex{}
	errMsgs := []string{}
	wg.Add(len(cpis))
	for _, cpi := range cpis {
		go func(cpi consul.AppInstanceInfo, wg *sync.WaitGroup) {
			defer wg.Done()
			var errMsg string = ""
			var err error = nil
			var pi *common.PodInfo = nil
			if pi, err = cpi.TestReq(); pi == nil && err != nil {
				errMsg = "error querying " + pi.IP + ":" + strconv.Itoa(pi.Port) + err.Error()
			} else if err != nil {
				errMsg = "app info does not match consul info " +
					err.Error() + " for " + pi.IP + ":" + strconv.Itoa(pi.Port)
			}
			if errMsg != "" {
				errAccMutex.Lock()
				defer errAccMutex.Unlock()
				errMsgs = append(errMsgs, errMsg)
			}
		}(cpi, &wg)
	}
	if len(errMsgs) == 0 {
		return nil
	}
	return errors.New(strings.Join(errMsgs, ", "))
}
