package entity

type DnsItem struct {
	Domain string `json:"domain"`
	Source []string `json:"source"`
	Type   string `json:"type"`
	Count *int64 `json:"count"`
	LastTime string `json:"last_time"`
	Dns map[string][]string `json:"dns"`
}

