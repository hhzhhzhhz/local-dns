package utils

import (
	"fmt"
	"github.com/hhzhhzhhz/local-dns/msg"
)

func VarifyDnsMgrRequst(ms []msg.Service) error {
	for _, m := range ms {
		if m.Host == "" || m.Key == "" {
			return fmt.Errorf("params is empty")
		}
	}
	return nil
}

func VarifyDnsDomain(ms []string) error {
	if len(ms) == 0 {
		return fmt.Errorf("params is empty")
	}
	return nil
}
