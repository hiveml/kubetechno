// Interceptor configuration.
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Settings struct {
	LowerBound      int
	UpperBound      int
	CertFilePath    string
	KeyFilePath     string
	DisallowedPorts map[int]interface{}
}

// Determine settings from env vars.
func New() (Settings, error) {
	settings := Settings{}
	var err error = nil
	if settings.LowerBound, err = osIntEnv("LOWER_BOUND"); err != nil {
		return Settings{}, err
	}
	if settings.UpperBound, err = osIntEnv("UPPER_BOUND"); err != nil {
		return Settings{}, err
	}
	settings.CertFilePath = "/etc/kubetechno/pems/cert.pem"
	settings.KeyFilePath = "/etc/kubetechno/pems/key.pem"
	ports, err := disallowedPorts()
	settings.DisallowedPorts = ports
	return settings, err
}

// Gets an env var as an int.
func osIntEnv(key string) (int, error) {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		err = errors.New("could not parse env with key " + key + " due to error " + err.Error())
	}
	return val, err
}

func disallowedPorts() (map[int]interface{}, error) {
	envVal := os.Getenv("DISALLOWED_PORTS")
	if envVal == "" {
		return make(map[int]interface{}), nil
	}
	portsStrList := strings.Split(envVal, ",")
	portsSet := make(map[int]interface{}, len(portsStrList))
	for _, strPort := range portsStrList {
		intPort, err := strconv.Atoi(strPort)
		if err != nil {
			return nil, err
		}
		portsSet[intPort] = nil
	}
	return portsSet, nil
}
