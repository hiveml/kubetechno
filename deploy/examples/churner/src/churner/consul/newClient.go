package consul

import (
	"github.com/hashicorp/consul/api"
	"os"
)

// Create a new consul api client
func NewClient() (*api.Client, error) {
	consulClientConfig := api.DefaultConfig()
	consulClientConfig.Address = os.Getenv("NODE_NAME") + ":8500"
	return api.NewClient(consulClientConfig) // todo: make the consul config more flexible
}
