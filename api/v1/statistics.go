package v1

import (
	"context"
	"encoding/json"
	"github.com/hhzhhzhhz/local-dns/backends"
	"github.com/hhzhhzhhz/local-dns/constant"
	"github.com/hhzhhzhhz/local-dns/entity"
	"github.com/hhzhhzhhz/local-dns/server"
	"github.com/hhzhhzhhz/local-dns/utils"
	"net/http"
	"sort"
)

func NewStatistics(b backends.Storage) *statistics {
	return &statistics{db: b}
}

type statistics struct {
	db backends.Storage
}

func (s *statistics) DnsRequestStatistics(w http.ResponseWriter, r *http.Request) {

	ks, err  := s.db.All(context.TODO(), server.RequestTotal)
	if err != nil {
		utils.HttpFailed(w, constant.ErrorDbOp)
		return
	}
	var items []*entity.DnsItem
	for _, k := range ks.Kvs {
		item := &entity.DnsItem{}
		if err := json.Unmarshal(k.Value, item); err != nil {
			continue
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return *items[i].Count > *items[j].Count
	})

	utils.HttpSuccess(w, items)
}


