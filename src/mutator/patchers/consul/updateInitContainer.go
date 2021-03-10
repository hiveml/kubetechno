package consul

import (
	"errors"
	"k8s.io/api/core/v1"
	"kubetechno/common/constants"
)

func (p Patcher) updateInitContainer(po *v1.Pod) error {
	pullPolicy, err := getPullPolicy(po.ObjectMeta.Annotations)
	if err != nil {
		return err
	}
	if po.Spec.InitContainers == nil {
		po.Spec.InitContainers = []v1.Container{}
	}
	consulInitImage := p.defaultConsulInitImage
	if po.ObjectMeta.Annotations != nil {
		if annotationConsulInitImage, in := po.ObjectMeta.Annotations[constants.ConsulConsulInitImage]; in {
			consulInitImage = annotationConsulInitImage
		}
	}
	kubetechnoInitConsul := v1.Container{
		Name:            "kubetechno-init-consul",
		Image:           consulInitImage,
		ImagePullPolicy: v1.PullPolicy(pullPolicy),
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "kubetechno-consul",
				MountPath: "/mvTarget",
			},
		},
	}
	po.Spec.InitContainers = append(
		po.Spec.InitContainers, kubetechnoInitConsul)
	return nil
}

func getPullPolicy(annotations map[string]string) (pullPolicy string, err error) {
	pullPolicy, in := annotations[constants.ConsulClientImagePullPolicy]
	if in &&
		pullPolicy != string(v1.PullAlways) &&
		pullPolicy != string(v1.PullNever) &&
		pullPolicy != string(v1.PullIfNotPresent) {
		err = errors.New("bad pull policy")
		return
	} else if !in {
		pullPolicy = string(v1.PullIfNotPresent)
	}
	return
}
