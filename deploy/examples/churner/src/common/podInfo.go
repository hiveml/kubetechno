package common

type PodInfo struct {
	Port      int    `json:"port"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Node      string `json:"node"`
	IP        string `json:"ip"`
}
