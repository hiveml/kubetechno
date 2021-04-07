package main

import (
	"errors"
	"fmt"
	"kubetechno/common/constants"
	"kubetechno/consulClient/clientWrapper"
	"kubetechno/consulClient/settings"
	"net/http"
	"os"
	"strconv"
)

func main() {
	action, err := getAction()
	s, err := settings.New()
	cw, err := clientWrapper.New(s)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if action == "check" {
		err = check(s, cw)
	} else if action == "dereg" {
		err = cw.Deregister()
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func getAction() (string, error) {
	action := os.Args[1]
	switch action {
	case constants.ConsulCheckCheckCmd:
		fallthrough
	case constants.ConsulCheckDeregCmd:
		return action, nil
	default:
		return "", errors.New("invalid action")
	}
}

func check(s settings.Settings, w clientWrapper.Wrapper) error {
	resp, err := http.Get(s.GetURL())
	if err != nil || (resp.StatusCode < 200 || resp.StatusCode > 299 && resp.StatusCode != 404) {
		if err == nil {
			err = errors.New("check status code was " + strconv.Itoa(resp.StatusCode))
		}
		if err2 := w.FailCheck(); err2 != nil {
			return errors.New(
				"check failed: " + err.Error() + " and consul update failed: " + err2.Error()) // err.Error exception
		}
		return errors.New("check failed: " + err.Error())
	}
	return w.PassCheck()
}
