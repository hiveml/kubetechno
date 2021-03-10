package main

import (
	"fmt"
	"kubetechno-churner/churner/consul"
	"os"
	"strconv"
	"sync"
	"time"

	"k8s.io/api/core/v1"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	ns := os.Getenv("namespace")
	replicas, err := strconv.Atoi(os.Getenv("replicas"))
	// setup clients
	cc, err := consul.NewClient()
	kc, err := NewK8sClient()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	logger := NewLog()
	// run checks in a loop
	for reqCount := 1; ; reqCount++ {
		log := logger.WithField(
			"run_start_time", time.Now().Format(time.Stamp)).WithField("request_count", reqCount)
		log.Info("new test")
		run(replicas, ns, cc, kc, log)
		time.Sleep(time.Second * time.Duration(30)) // todo: parameterize this
	}
}

// Get info, compare results, churn pods
func run(replicas int, ns string, cc *api.Client, kc *kubernetes.Clientset, log *logrus.Entry) {
	cpis, k8sPods, err := getInfo(replicas, ns, cc, kc, log)
	if err != nil {
		log.Error("failure to get consul and/or k8s info " + err.Error())
		return
	}
	pass := true
	comparisonsWG := sync.WaitGroup{}
	comparisonsWG.Add(2)
	go func(pis []consul.AppInstanceInfo, pods []v1.Pod, wg *sync.WaitGroup) {
		defer wg.Done()
		if err = compareK8sPodsToConsulServices(cpis, k8sPods); err != nil {
			pass = false
			log.Error("error comparing k8s pods to consul service instances " + err.Error())
		}
	}(cpis, k8sPods, &comparisonsWG)
	go func(ns string, pis []consul.AppInstanceInfo, pods []v1.Pod, wg *sync.WaitGroup) {
		defer wg.Done()
		if err = compareAppDataToConsulServices(cpis); err != nil {
			pass = false
			log.Error("error comparing pod data to consul services " + err.Error())
		}
	}(ns, cpis, k8sPods, &comparisonsWG)
	comparisonsWG.Wait()
	if pass {
		log.Info("comparisons passed")
	} else {
		log.Info("comparisons failed")
	}
	churnPods(k8sPods, kc, log)
	time.Sleep(time.Second * time.Duration(5))
}

func NewK8sClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewLog() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	return logger
}
