package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"math/rand"
	"sync"
)

func churnPods(pods []v1.Pod, kc *kubernetes.Clientset, log *logrus.Entry) {
	wg := sync.WaitGroup{}
	wg.Add(len(pods))
	log.Info("churning started")
	for _, pod := range pods {
		go func(pod v1.Pod, wg *sync.WaitGroup) {
			defer wg.Done()
			if rand.Intn(10) < 2 { // todo: parameterize the chance of a churn
				kc.CoreV1().Pods(pod.Namespace).Delete(context.Background(), pod.Name, metaV1.DeleteOptions{})
				log.Info("deleted " + pod.Name)
			}
		}(pod, &wg)
	}
	wg.Wait()
}
