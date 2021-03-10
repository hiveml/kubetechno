package clientWrapper

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"kubetechno/consulClient/settings"
)

type Wrapper struct {
	c *api.Client
	s settings.Settings
}

func (w Wrapper) Register() error {
	reg := w.s.Registration()
	return w.c.Agent().ServiceRegister(&reg)
}

func (w Wrapper) Deregister() error {
	return w.c.Agent().ServiceDeregister(w.s.Id())
}

func (w Wrapper) FailCheck() error {
	return w.c.Agent().FailTTL(w.s.CheckID(), "")
}

func (w Wrapper) PassCheck() error {
	var err error = nil
	if err = w.c.Agent().PassTTL(w.s.CheckID(), ""); err == nil {
		return nil
	}
	fmt.Println("could not pass ttl first time " + err.Error())
	if err = w.Register(); err != nil {
		fmt.Println("could not reg service" + err.Error())
		return err
	}
	fmt.Println("service registered")
	if err = w.c.Agent().PassTTL(w.s.CheckID(), ""); err == nil {
		return nil
	}
	fmt.Println("could not pass ttl second time " + err.Error())
	return err
}
