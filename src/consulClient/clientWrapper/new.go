package clientWrapper

import (
	"github.com/hashicorp/consul/api"
	"kubetechno/consulClient/settings"
)

func New(s settings.Settings) (Wrapper, error) {
	config := api.DefaultConfig()
	config.Address = s.ConsulNode() + ":8500"
	c, err := api.NewClient(config)
	return Wrapper{
		c: c,
		s: s,
	}, err
}
