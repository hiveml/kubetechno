package consul

import "github.com/hashicorp/consul/api"

// Create a new consul api client
func NewClient() (*api.Client, error) {
	return api.NewClient(api.DefaultConfig()) // todo: make the consul config more flexible
}
