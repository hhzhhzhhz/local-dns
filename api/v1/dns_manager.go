package v1

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hhzhhzhhz/local-dns/constant"
	"github.com/hhzhhzhhz/local-dns/log"
	"github.com/hhzhhzhhz/local-dns/msg"
	"github.com/hhzhhzhhz/local-dns/server"
	"github.com/hhzhhzhhz/local-dns/utils"
	"io/ioutil"
)

func NewDnsManager(b server.Backend) *dnsManager {
	return &dnsManager{backend: b}
}

type dnsManager struct {
	backend server.Backend
}

func (d *dnsManager) AddDnsRecord(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		utils.HttpFailed(c, constant.ErrorReadBody)
		return
	}
	var records []msg.Service
	if err := json.Unmarshal(b, &records); err != nil {
		utils.HttpFailed(c, constant.ErrorParams)
		return
	}
	if err := utils.VarifyDnsMgrRequst(records); err != nil {
		return
	}
	rmap := map[string]*[]msg.Service{}
	for _, r := range records {
		v , ok := rmap[r.Key]
		if ok {
			*v = append(*v, r)
			continue
		}
		rmap[r.Key] = &[]msg.Service{r}
	}
	for k, r := range rmap {
		if err := d.backend.SaveRecords(context.TODO(), k, *r); err != nil {
			utils.HttpFailed(c, constant.ErrorDbOp)
			return
		}
	}
	log.Logger().Info("add dns record %+v", rmap)
	utils.HttpSuccess(c, "success")
}
// RemoveDnsRecord 移除域名Dns 记录
func (d *dnsManager) RemoveDnsRecord(c *gin.Context) {
	domains := []string{c.Param("domains")}
	if err := utils.VarifyDnsDomain(domains); err != nil {
		return
	}
	if len(domains) == 0 {
		utils.HttpFailed(c, constant.ErrorParams)
		return
	}
	if err := d.backend.RemoveRecords(context.TODO(), domains); err != nil {
		utils.HttpFailed(c, constant.ErrorDbOp)
		return
	}
	log.Logger().Info("remove dns record %+v", domains)

	utils.HttpSuccess(c, "success")

}