package handler

import (
	"github.com/sirupsen/logrus"
	"kubetechno/interceptor/orchestrator"
	"os"
)

// New creates a new server instance
func New(o *orchestrator.Orchestrator) *Handler {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	handler := Handler{
		o:      o,
		logger: logger,
	}
	return &handler
}
