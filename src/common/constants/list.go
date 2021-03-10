package constants

import "strconv"

const (
	Kubetechno                  = "kubetechno"
	ConsulServiceName           = Kubetechno + "_consul_service"
	ConsulBufferSecs            = Kubetechno + "_consul_buffer_seconds"
	ConsulTimeoutSeconds        = Kubetechno + "_consul_timeout_seconds"
	ConsulPeriodSeconds         = Kubetechno + "_consul_period_seconds"
	ConsulInitialDelaySeconds   = Kubetechno + "_consul_initial_delay_seconds"
	ConsulConsulInitImage       = Kubetechno + "_consul_init_image"
	ConsulClientImagePullPolicy = Kubetechno + "_consul_client_pull_policy"
	ConsulCheckPath             = Kubetechno + "_consul_check_path"
	ConsulCheckDeregCmd         = "dereg"
	ConsulCheckCheckCmd         = "check"
	ConsulClientPath            = "/" + Kubetechno + "Consul/client"
	KubetechnoDashNode          = Kubetechno + "-node"
	KubetechnoSlashPort         = Kubetechno + "/port"
	KubetechnoActiveStatus      = "active"
)

func PortN(i int) string {
	return "PORT" + strconv.Itoa(i)
}
