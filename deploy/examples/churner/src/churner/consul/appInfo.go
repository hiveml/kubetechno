package consul

import (
	"encoding/json"
	"errors"
	"kubetechno-churner/common"
	"net/http"
	"strconv"
)

type AppInstanceInfo struct {
	Port int
	IP   string
}

func (aii AppInstanceInfo) TestReq() (*common.PodInfo, error) {
	getTarget := "http://" + aii.IP + ":" + strconv.Itoa(aii.Port)
	resp, err := http.Get(getTarget)
	if err != nil {
		return nil, err
	}
	pi := common.PodInfo{}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err = json.NewDecoder(resp.Body).Decode(&pi); err != nil {
		return nil, err
	}
	if pi.IP == aii.IP && pi.Port == aii.Port {
		return &pi, nil
	}
	return &pi, errors.New("pod info does not match what was expected")
}
