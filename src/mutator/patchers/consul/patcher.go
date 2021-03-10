package consul

import (
	"k8s.io/api/core/v1"
)

type Patcher struct {
	defaultConsulInitImage string
}

func (p Patcher) UpdatePod(po *v1.Pod) (bool, error) {
	var err error = nil
	var hasName = false
	if _, hasName, err = getConsulServiceName(po.ObjectMeta.Annotations); err != nil {
		return false, err
	} else if !hasName {
		return false, nil
	} else if err = p.updateInitContainer(po); err != nil {
		return true, err
	} else if err = p.updateConsulContainer(po); err != nil {
		return true, err
	}
	p.updateVolumes(po)
	return true, nil
}
