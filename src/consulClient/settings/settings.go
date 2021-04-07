package settings

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
)

type Settings struct {
	serviceName    string
	consulNodeName string
	port           int
	timeoutSecs    int
	bufferSecs     int
	periodSeconds  int
	path           string
}

func (s Settings) ConsulNode() string {
	return s.consulNodeName
}

func (s Settings) Registration() api.AgentServiceRegistration {
	reg := api.AgentServiceRegistration{
		Name: s.ServiceName(),
		ID:   s.Id(),
		Port: s.Port(),
		Checks: api.AgentServiceChecks{
			{
				Name:                           s.ServiceName(),
				TTL:                            s.TTL(),
				DeregisterCriticalServiceAfter: s.DeregisterCriticalServiceAfter(),
			},
		},
	}
	b, _ := json.Marshal(reg)
	fmt.Println(string(b))
	return reg
}

func (s Settings) ServiceName() string {
	return s.serviceName
}

func (s Settings) Id() string {
	return s.serviceName + ":" + strconv.Itoa(s.port)
}

func (s Settings) Port() int {
	return s.port
}

func (s Settings) TTL() string {
	return strconv.Itoa(s.periodSeconds+s.bufferSecs) + "s"
}

func (s Settings) DeregisterCriticalServiceAfter() string {
	return strconv.Itoa(s.periodSeconds*2+s.bufferSecs) + "s"
}

func (s Settings) Checks() api.AgentServiceChecks {
	check := api.AgentServiceCheck{
		CheckID:                        s.CheckID(),
		Name:                           s.ServiceName(),
		TTL:                            s.TTL(),
		DeregisterCriticalServiceAfter: s.DeregisterCriticalServiceAfter(),
	}
	return []*api.AgentServiceCheck{&check}
}

func (s Settings) CheckID() string {
	return "service:" + s.Id()
}

func (s Settings) GetURL() string {
	return "http://" + s.consulNodeName +  ":" + strconv.Itoa(s.Port()) + s.path
}
