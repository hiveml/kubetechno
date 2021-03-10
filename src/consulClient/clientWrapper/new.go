package clientWrapper

import (
	"github.com/hashicorp/consul/api"
	"kubetechno/consulClient/settings"
)

func New(s settings.Settings) (Wrapper, error) {
	c, err := api.NewClient(api.DefaultConfig())
	return Wrapper{
		c: c,
		s: s,
	}, err
}
