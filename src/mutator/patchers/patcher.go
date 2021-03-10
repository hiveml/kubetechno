package patchers

import (
	"k8s.io/api/core/v1"
)

type Patcher interface {
	UpdatePod(po *v1.Pod) (change bool, err error)
}
