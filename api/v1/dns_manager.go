package v1

import (
	"context"
	"encoding/json"
	"github.com/hhzhhzhhz/local-dns/constant"
	"github.com/hhzhhzhhz/local-dns/log"
	"github.com/hhzhhzhhz/local-dns/msg"
	"github.com/hhzhhzhhz/local-dns/server"
	"github.com/hhzhhzhhz/local-dns/utils"
	"io/ioutil"
	"net/http"
)

func NewDnsManager(b server.Backend) *dnsManager {
	return &dnsManager{backend: b}
}

type dnsManager struct {
	backend server.Backend
}

func (d *dnsManager) AddDnsRecord(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.HttpFailed(w, constant.ErrorReadBody)
		return
	}
	var records []msg.Service
	if err := json.Unmarshal(b, &records); err != nil {
		utils.HttpFailed(w, constant.ErrorParams)
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
			utils.HttpFailed(w, constant.ErrorDbOp)
			return
		}
	}
	log.Logger().Info("add dns record %+v", rmap)

	utils.HttpSuccess(w, "success")
}
// RemoveDnsRecord 移除域名Dns 记录
func (d *dnsManager) RemoveDnsRecord(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	domains := r.Form["domains"]
	if err := utils.VarifyDnsDomain(domains); err != nil {
		return
	}
	if len(domains) == 0 {
		utils.HttpFailed(w, constant.ErrorParams)
		return
	}
	if err := d.backend.RemoveRecords(context.TODO(), domains); err != nil {
		utils.HttpFailed(w, constant.ErrorDbOp)
		return
	}
	log.Logger().Info("remove dns record %+v", domains)

	utils.HttpSuccess(w, "success")

}

func (d *dnsManager) ListRecord(w http.ResponseWriter, r *http.Request) {
	rds, err := d.backend.AllRecord(context.TODO())
	if err != nil {
		utils.HttpFailed(w, constant.ErrorDbOp)
		return
	}
	utils.HttpSuccess(w, rds)
}