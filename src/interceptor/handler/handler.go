// Handles http requests for pod binding requests.
package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"kubetechno/interceptor/orchestrator"
	"net/http"
)

// Handles pod binding requests from the k8s api server by it's wrapped orchestrator to assign ports.
type Handler struct {
	o      *orchestrator.Orchestrator
	logger *logrus.Logger
}

// ServeHTTP handles http requests to assign ports to kubetechno pods.
func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handling")
	ls := LogStruct{}
	ls.LastStep = "start"
	defer func() {
		if bytes, err := json.Marshal(ls); err != nil {
			logString := "{ \"log-print-err\": " + err.Error() + ", " + "\"log-string\":" + "}"
			s.logger.Info(logString)
		} else {
			s.logger.Info(string(bytes))
		}
	}()

	ls.LastStep = "parse input"
	uuid, nsName, poName, noName, apiVersion, err := parseInput(r)
	if err != nil {
		ls.Err = err.Error()
		w.WriteHeader(503)
		return
	}
	ls.Uuid = uuid
	ls.NsName = nsName
	ls.PoName = poName
	ls.NodeName = noName

	ls.LastStep = "assign ports"
	config, patches, err := s.o.AssignPorts(nsName, poName, noName)
	ls.Patches = patches
	ls.Config = config
	if err != nil {
		ls.Err = "error with port assignment " + err.Error()
		return
	}

	ls.LastStep = "create response"
	rsp := Response{
		APIVersion: apiVersion,
		Kind:       "AdmissionReview",
		Response: AllowRsp{
			UID:     uuid,
			Allowed: true,
		},
	}
	rspBytes, err := json.Marshal(rsp)
	if _, err = w.Write(rspBytes); err != nil {
		ls.Err = err.Error()
		w.WriteHeader(503)
		return
	}
	fmt.Println("rsp: " + string(rspBytes))
	w.Header().Add("Content-Type", "application/json")
}
