package log

import (
	"kubetechno/common/patch"
)

func NewLogEntry() *Entry {
	return &Entry{
		LastStep: "start",
	}
}

type Entry struct {
	LastStep string
	Err      string
	Patches  []patch.Patch
}
