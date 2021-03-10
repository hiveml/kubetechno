package consul

import "github.com/hashicorp/consul/api"

func GetAIIs(client *api.Client, serviceName string, passingOnly bool) ([]AppInstanceInfo, error) {
	query := api.QueryOptions{}
	list, _, err := client.Health().Service(serviceName, "", passingOnly, &query)
	if err != nil {
		return nil, err
	}
	cpis := []AppInstanceInfo{}
	for _, entry := range list {
		cpis = append(cpis, AppInstanceInfo{
			Port: entry.Service.Port,
			IP:   entry.Node.Address,
		})
	}
	return cpis, err
}
