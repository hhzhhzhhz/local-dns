package web

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hhzhhzhhz/local-dns/backends"
	"github.com/hhzhhzhhz/local-dns/constant"
	"github.com/hhzhhzhhz/local-dns/entity"
	"github.com/hhzhhzhhz/local-dns/server"
	"github.com/hhzhhzhhz/local-dns/utils"
	"net/http"
	"sort"
)

func NewWeb(b backends.Storage, backend server.Backend) *web {
	return &web{db: b, backend: backend}
}

type web struct {
	db backends.Storage
	backend server.Backend
}

func (s *web) DnsRequestStatistics(c *gin.Context) {
	ks, err  := s.db.All(context.TODO(), server.RequestTotal)
	if err != nil {
		c.HTML(http.StatusOK, "show_dns_query.tmpl", gin.H{"data": nil})
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
	c.HTML(http.StatusOK, "show_dns_query.tmpl", gin.H{"data": items})
}

func (s *web) Root(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"data": ""})
}

func (d *web) ListRecord(c *gin.Context) {
	rds, err := d.backend.AllRecord(context.TODO())
	if err != nil {
		utils.HttpFailed(c, constant.ErrorDbOp)
		return
	}
	c.HTML(http.StatusOK, "show_dns_record.tmpl", gin.H{"data": rds})
}



