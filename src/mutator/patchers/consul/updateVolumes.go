package consul

import (
	"k8s.io/api/core/v1"
)

func (p Patcher) updateVolumes(po *v1.Pod) {
	if po.Spec.Volumes == nil {
		po.Spec.Volumes = []v1.Volume{}
	}
	po.Spec.Volumes = append(po.Spec.Volumes, v1.Volume{
		Name: "kubetechno-consul",
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{},
		},
	})
}
