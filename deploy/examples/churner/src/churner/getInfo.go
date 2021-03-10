package main

import (
	"context"
	"errors"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kubetechno-churner/churner/consul"
	"strconv"
	"time"
)

func getInfo(appCount int, ns string, cc *consulAPI.Client, kc *kubernetes.Clientset, log *logrus.Entry) ([]consul.AppInstanceInfo, []v1.Pod, error) {
	log.Info("get info")
	var consulErr error = nil
	var k8sErr error = nil
	cpis := []consul.AppInstanceInfo{}
	poList := []v1.Pod{}
	getK8sInfo(appCount, ns, kc, &poList, log, &k8sErr)
	getConsulInfo(appCount, &cpis, cc, log, &consulErr)
	rtnErrMsg := ""
	if consulErr != nil {
		rtnErrMsg += "consul retrieval err: " + consulErr.Error()
	}
	if k8sErr != nil {
		if rtnErrMsg != "" {
			rtnErrMsg += " "
		}
		rtnErrMsg += "k8s retrieval err: " + k8sErr.Error()
	}
	var err error = nil
	if rtnErrMsg != "" {
		err = errors.New(rtnErrMsg)
	}
	return cpis, poList, err
}

func getConsulInfo(appCount int, cpis *[]consul.AppInstanceInfo, cc *consulAPI.Client, log *logrus.Entry, rtnErr *error) {
	log.Info("waiting for all consul instances to be ready")
	const intervalLenSecs = 5 // todo: move to params or make part of a struct
	for {
		passingOnlyCpis, err := consul.GetAIIs(cc, "kubetechno-churner-app", true)
		if err != nil {
			*rtnErr = err
			return
		}
		if len(passingOnlyCpis) == appCount {
			log.Info("All consul service instances are passing!")
			*cpis = passingOnlyCpis
			break
		}
		log.Info("only found " + strconv.Itoa(len(passingOnlyCpis)) + " passing consul services")
		time.Sleep(time.Second * time.Duration(intervalLenSecs))
	}
}

func getK8sInfo(appCount int, ns string, kc *kubernetes.Clientset, poList *[]v1.Pod, log *logrus.Entry, rtnErr *error) {
	log.Info("waiting for all consul containers to be ready")
	const intervalLenSecs = 5 // todo: move to params or make part of a struct
	var err error = nil
	for {
		var newPoList []v1.Pod
		var returnedList *v1.PodList
		returnedList, err = kc.CoreV1().Pods(ns).List(context.Background(),
			metaV1.ListOptions{LabelSelector: "app=kubetechno-churner-app"})
		if err != nil {
			*rtnErr = err
			log.Info(err.Error())
			break
		} else if returnedList != nil {
			newPoList = returnedList.Items
			readyCount := 0
			for _, po := range newPoList {
				if po.Status.ContainerStatuses != nil && po.Status.ContainerStatuses[0].Ready {
					readyCount += 1
				}
			}
			if readyCount == appCount {
				log.Info("All pods are ready!")
				*poList = newPoList
				break
			}
		}
		time.Sleep(time.Second * time.Duration(intervalLenSecs))
	}
}
