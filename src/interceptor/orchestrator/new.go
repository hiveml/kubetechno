package orchestrator

import (
	"sync"
)

func New(lowerBound, upperBound int, kc k8sClient, disallowedPorts map[int]interface{}) *Orchestrator {
	return &Orchestrator{
		superLock:       &sync.Mutex{},
		nodeLocks:       make(map[string]*sync.Mutex),
		kClient:         kc,
		lowerBound:      lowerBound,
		upperBound:      upperBound,
		disallowedPorts: disallowedPorts,
	}
}
