package domain

type Log struct {
	Service       string `json:"service"`
	ContainerName string `json:"container_name"`
	Time          string `json:"time"`
	RemoteIP      string `json:"remote_ip"`
	Host          string `json:"host"`
	Method        string `json:"method"`
	Uri           string `json:"uri"`
	UserAgent     string `json:"user_agent"`
	Status        string `json:"status"`
	Latency       string `json:"latency"`
	LatencyHuman  string `json:"latency_human"`
	Error         string `json:"error"`
}
