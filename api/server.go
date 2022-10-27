package api

import (
	"github.com/gin-gonic/gin"
	v12 "github.com/hhzhhzhhz/local-dns/api/v1"
	v1 "github.com/hhzhhzhhz/local-dns/api/web"
	"github.com/hhzhhzhhz/local-dns/backends"
	"github.com/hhzhhzhhz/local-dns/constant"
	"github.com/hhzhhzhhz/local-dns/log"
	"github.com/hhzhhzhhz/local-dns/server"
	"github.com/hhzhhzhhz/local-dns/utils"
	"net/http"
	_ "net/http/pprof"
)

func NewApiServer(cfg *server.Config, backend server.Backend, db backends.Storage) *HttpServer {
	return &HttpServer{
		cfg: cfg,
		backend: backend,
		db: db,
	}
}

type HttpServer struct {
	cfg     *server.Config
	r       *gin.Engine
	backend server.Backend
	db      backends.Storage
}

func (h *HttpServer) Run() error  {
	dnsMgr := v12.NewDnsManager(h.backend)
	web := v1.NewWeb(h.db, h.backend)
	gin.SetMode(gin.ReleaseMode)
	h.r = gin.Default()
	h.Default()
	h.r.GET("/", web.Root)
	h.r.GET("/web/dns", web.DnsRequestStatistics)
	h.r.GET("/web/ListRecord", web.ListRecord)
	h.r.POST("/api/dnsMgr/addRecord", dnsMgr.AddDnsRecord)
	h.r.GET("/api/dnsMgr/RemoveDnsRecord", dnsMgr.RemoveDnsRecord)
	log.Logger().Info("Api is listening and serving on %s", h.cfg.ApiAddr)
	if err := h.r.Run(h.cfg.ApiAddr); err != nil {
		log.Logger().Error("Api.server stop cause=%s", err.Error())
		return err
	}
	return nil
}

func (h *HttpServer) Default()  {
	h.r.LoadHTMLGlob(h.cfg.StaticPath)
	h.r.StaticFS("/static", http.Dir("./api/static"))
	h.r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.tmpl", "")
	})
	h.r.NoMethod(func(c *gin.Context) {
		utils.HttpFailed(c, constant.ApiResponse{Code: 405, Info: "method does not allow"})
	})
}

func (h *HttpServer) Stop() error {
	return nil
}