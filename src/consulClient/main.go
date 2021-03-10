package main

import (
	"errors"
	"fmt"
	"kubetechno/common/constants"
	"kubetechno/consulClient/clientWrapper"
	"kubetechno/consulClient/settings"
	"net/http"
	"os"
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
		fmt.Println("running check")
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
	fmt.Println(s.GetURL())
	resp, err := http.Get(s.GetURL())
	if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		if err2 := w.FailCheck(); err2 != nil {
			return errors.New(
				"check failed: " + err.Error() + " and consul update failed: " + err2.Error())
		}
		return errors.New("check failed: " + err.Error())
	}
	return w.PassCheck()
}
