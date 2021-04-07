package settings

import (
	"fmt"
	"kubetechno/common/constants"
	"os"
	"strconv"
)

func New() (s Settings, err error) {
	consulService := os.Getenv(constants.ConsulServiceName)
	bufferSecsStr := os.Getenv(constants.ConsulBufferSecs)
	bufferSecs, err := strconv.Atoi(bufferSecsStr)
	timeoutSecondsStr := os.Getenv(constants.ConsulTimeoutSeconds)
	timeoutSeconds, err := strconv.Atoi(timeoutSecondsStr)
	periodSecondsStr := os.Getenv(constants.ConsulPeriodSeconds)
	periodSeconds, err := strconv.Atoi(periodSecondsStr)
	portStr := os.Getenv("PORT0")
	port, err := strconv.Atoi(portStr)
	consulNodeName := os.Getenv("CONSUL_NODE_NAME")
	path := os.Getenv(constants.ConsulCheckPath)
	fmt.Println(consulNodeName)
	return Settings{
		serviceName:    consulService,
		consulNodeName: consulNodeName,
		port:           port,
		bufferSecs:     bufferSecs,
		periodSeconds:  periodSeconds,
		timeoutSecs:    timeoutSeconds,
		path: path,
	}, err
}
