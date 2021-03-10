package consul

import (
	"errors"
	"k8s.io/api/core/v1"
	"kubetechno/common/constants"
	"strconv"
)

func (p Patcher) updateConsulContainer(po *v1.Pod) error {
	var err error = nil
	if err = p.addDataFromLabels(po); err != nil {
		return err
	}
	p.addVolumeMounts(po)
	p.addDeregister(po)
	return nil
}

func (p Patcher) addDeregister(po *v1.Pod) {
	container := &po.Spec.Containers[0]
	if container.Lifecycle == nil {
		container.Lifecycle = &v1.Lifecycle{}
	}
	container.Lifecycle.PreStop = &v1.Handler{
		Exec: &v1.ExecAction{
			Command: []string{constants.ConsulClientPath, constants.ConsulCheckDeregCmd},
		},
	}
}

func (p Patcher) addVolumeMounts(po *v1.Pod) {
	container := po.Spec.Containers[0]
	if container.VolumeMounts == nil {
		container.VolumeMounts = []v1.VolumeMount{}
	}
	container.VolumeMounts = append(container.VolumeMounts,
		v1.VolumeMount{Name: "kubetechno-consul", MountPath: "/kubetechnoConsul"})
	po.Spec.Containers[0] = container

}

func (p Patcher) addDataFromLabels(po *v1.Pod) error {
	if po.Spec.Containers[0].Env == nil {
		po.Spec.Containers[0].Env = []v1.EnvVar{}
	}
	checkPath, err := getConsulCheckPath(po.Annotations)
	if err != nil {
		return errors.New("could not get consul check path")
	}
	addToEnvVars(po, constants.ConsulCheckPath, checkPath)
	// this is called before to check for presence
	// and errors so 2 of the rtn vals are not needed
	serviceName, _, _ := getConsulServiceName(po.Annotations)
	addToEnvVars(po, constants.ConsulServiceName, serviceName)

	intSettings := make(map[string]int, 4)
	intSettings[constants.ConsulBufferSecs] = 10
	intSettings[constants.ConsulTimeoutSeconds] = 30
	intSettings[constants.ConsulPeriodSeconds] = 150
	intSettings[constants.ConsulInitialDelaySeconds] = 30
	setIntValues := make(map[string]int, 4)

	for _, settingName := range []string{constants.ConsulBufferSecs,
		constants.ConsulTimeoutSeconds, constants.ConsulPeriodSeconds, constants.ConsulInitialDelaySeconds} {
		defaultVal, _ := intSettings[settingName]
		valToUse := defaultVal
		if annotationsVal, in := po.Annotations[settingName]; in {
			var err error = nil
			if valToUse, err = strconv.Atoi(annotationsVal); err != nil {
				return errors.New("non int annotation val for " + settingName)
			}

		}
		addToEnvVars(po, settingName, strconv.Itoa(valToUse))
		setIntValues[settingName] = valToUse
	}
	po.Spec.Containers[0].ReadinessProbe = &v1.Probe{
		Handler: v1.Handler{
			Exec: &v1.ExecAction{
				Command: []string{constants.ConsulClientPath, constants.ConsulCheckCheckCmd},
			},
		},
		InitialDelaySeconds: int32(setIntValues[constants.ConsulInitialDelaySeconds]),
		TimeoutSeconds:      int32(setIntValues[constants.ConsulTimeoutSeconds]),
		PeriodSeconds:       int32(setIntValues[constants.ConsulPeriodSeconds]),
		SuccessThreshold:    1,
		FailureThreshold:    1,
	}
	return nil
}

func getConsulServiceName(annotations map[string]string) (serviceName string, hasName bool, err error) {
	if annotations == nil {
		return
	}
	serviceName, hasName = annotations[constants.ConsulServiceName]
	if !hasName {
		return
	}
	if serviceName == "" {
		err = errors.New("bad consul service name")
		return
	}
	return
}

func getConsulCheckPath(annotations map[string]string) (string, error) {
	if path, in := annotations[constants.ConsulCheckPath]; in {
		return path, nil
	}
	return "", errors.New("no check path present")
}

func addToEnvVars(po *v1.Pod, envName, envVal string) {
	po.Spec.Containers[0].Env = append(po.Spec.Containers[0].Env,
		v1.EnvVar{Name: envName, Value: envVal})
}
