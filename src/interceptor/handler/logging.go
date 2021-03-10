package handler

import (
	"kubetechno/common"
	"kubetechno/common/patch"
)

// Printable struct describing the outcome of handling a binding request.
type LogStruct struct {
	Uuid     string         `json:"pod-uuid"`
	NsName   string         `json:"namespace"`
	PoName   string         `json:"pod-name"`
	NodeName string         `json:"node"`
	Config   *common.Config `json:"config,omitempty"`
	Patches  []patch.Patch  `json:"patches,omitempty"`
	Err      string         `json:"err,omitempty"`
	LastStep string         `json:"lastStep"`
}
