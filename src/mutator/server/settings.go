package main

import (
	"errors"
	"os"
)

type Settings struct {
	CertFilePath       string
	KeyFilePath        string
	DefaultConsulImage string
}

func New() (s Settings, err error) {
	s.CertFilePath = "/etc/kubetechno/pems/cert.pem"
	s.KeyFilePath = "/etc/kubetechno/pems/key.pem"
	s.DefaultConsulImage = os.Getenv("default_consul_image")
	if s.DefaultConsulImage == "" {
		err = errors.New("no default consul image")
	}
	return s, err
}
