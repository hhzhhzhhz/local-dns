package api

import (
	v1 "github.com/hhzhhzhhz/local-dns/api/v1"
	"github.com/hhzhhzhhz/local-dns/backends"
	"github.com/hhzhhzhhz/local-dns/log"
	"github.com/hhzhhzhhz/local-dns/server"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func NewApiServer(addr string, backend server.Backend, db backends.Storage) *HttpServer {
	return &HttpServer{
		Addr: addr,
		backend: backend,
		db: db,
	}
}

type HttpServer struct {
	Addr string
	Server *http.Server
	backend server.Backend
	db backends.Storage
}

func (h *HttpServer) Run() error  {
	mux := http.DefaultServeMux
	dnsMgr := v1.NewDnsManager(h.backend)
	dnsStc := v1.NewStatistics(h.db)
	mux.HandleFunc("/api/dns", dnsStc.DnsRequestStatistics)
	mux.HandleFunc("/api/dnsMgr/addRecord", dnsMgr.AddDnsRecord)
	mux.HandleFunc("/api/dnsMgr/RemoveDnsRecord", dnsMgr.RemoveDnsRecord)
	mux.HandleFunc("/api/dnsMgr/ListRecord", dnsMgr.ListRecord)
	h.Server = &http.Server{
		Addr:           h.Addr,
		Handler:        mux,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Logger().Info("Api is listening and serving on %s", h.Addr)
	if err := h.Server.ListenAndServe(); err != nil {
		log.Logger().Error("Api.server stop cause=%s", err.Error())
		return err
	}
	return nil
}

func (h *HttpServer) Stop() error {
	return nil
}

type Resp struct {
	Domain string
	Info interface{}
}