package entity

type DnsItem struct {
	Source   string                  `json:"source"`
	Answers  map[string][]*DnsAnswer `json:"answers"`
	Type     string                  `json:"type"`
	Count    *int64                  `json:"count"`
	LastTime string                  `json:"last_time"`
}

type DnsAnswer struct {
	Qty    string   `json:"qty"`
	Result []string `json:"result"`
}
