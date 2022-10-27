// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hhzhhzhhz/local-dns/anticorruption"
	"github.com/hhzhhzhhz/local-dns/api"
	"github.com/hhzhhzhhz/local-dns/backends"
	"github.com/hhzhhzhhz/local-dns/log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hhzhhzhhz/local-dns/metrics"
	"github.com/hhzhhzhhz/local-dns/msg"
	"github.com/hhzhhzhhz/local-dns/server"

	"github.com/miekg/dns"
)

var (
	tlskey     = ""
	tlspem     = ""
	cacert     = ""
	username   = ""
	password   = ""
	config     = &server.Config{ReadTimeout: 0, Domain: "", DnsAddr: "", DNSSEC: ""}
	nameserver = ""
	machine    = ""
	stub       = false
	ctx        = context.Background()
)

func env(key, def string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return def
}

func intEnv(key string, def int) int {
	if x := os.Getenv(key); x != "" {
		if v, err := strconv.ParseInt(x, 10, 0); err == nil {
			return int(v)
		}
	}
	return def
}

func boolEnv(key string, def bool) bool {
	if x := os.Getenv(key); x != "" {
		if v, err := strconv.ParseBool(x); err == nil {
			return v
		}
	}
	return def
}

func SetEnv()  {
	os.Setenv("SKYDNS_DOMAIN", "cluster.local.")
	os.Setenv("SKYDNS_NAMESERVERS", "114.114.114.114:53,8.8.8.8:53")
	os.Setenv("SKYDNS_API", "0.0.0.0:80")
	os.Setenv("DB_FILE", "D:\\dns\\dns_db_new")
	os.Setenv("LOG_DIR", "D:\\dns\\log")
	os.Setenv("IP_DB_FILE", "./var/ip2region.xdb")
	os.Setenv("SKYDNS_ADDR", "0.0.0.0:53")
	os.Setenv("STATIC_PATH", "./api/static/html/*")
	//os.Setenv("SKYDNS_DOMAIN", "cluster.local")
}

func init() {
	os.Setenv("env", "dev")
	if os.Getenv("env") == "dev" {
		SetEnv()
	}
	flag.StringVar(&config.Domain, "domain", env("SKYDNS_DOMAIN", "skydns.local."), "domain to anchor requests to (SKYDNS_DOMAIN)")
	flag.StringVar(&config.DnsAddr, "addr", env("SKYDNS_ADDR", "127.0.0.1:53"), "ip:port to bind to (SKYDNS_ADDR)")
	flag.StringVar(&config.ApiAddr, "api", env("SKYDNS_API", "127.0.0.1:9876"), "ip:port to bind to (API_ADDR)")
	flag.StringVar(&config.DbFile, "dbfile", env("DB_FILE", "/var/dnsCore"), "ip:port to bind to (API_ADDR)")
	flag.StringVar(&config.LogDir, "logdir", env("LOG_DIR", "/var/log"), "ip:port to bind to (API_ADDR)")
	flag.StringVar(&config.IpDbFile, "IpDbFile", env("IP_DB_FILE", "/var/ip2region.xdb"), "ip:port to bind to (API_ADDR)")
	flag.StringVar(&config.StaticPath, "StaticPath", env("STATIC_PATH", "/var/templates/*"), "static path")


	flag.StringVar(&nameserver, "nameservers", env("SKYDNS_NAMESERVERS", ""), "nameserver address(es) to forward (non-local) queries to e.g. 8.8.8.8:53,8.8.4.4:53")
	flag.BoolVar(&config.NoRec, "no-rec", false, "do not provide a recursive service")
	flag.StringVar(&config.DNSSEC, "dnssec", "", "basename of DNSSEC key file e.q. Kskydns.local.+005+38250")
	flag.StringVar(&config.Local, "local", ".*", "optional unique value for this skydns instance")
	flag.DurationVar(&config.ReadTimeout, "rtimeout", 2*time.Second, "read timeout")
	flag.BoolVar(&config.RoundRobin, "round-robin", true, "round robin A/AAAA replies")
	flag.BoolVar(&config.NSRotate, "ns-rotate", true, "round robin selection of nameservers from among those listed")
	flag.BoolVar(&stub, "stubzones", false, "support stub zones")
	flag.BoolVar(&config.Verbose, "verbose", false, "log queries")
	flag.BoolVar(&config.Systemd, "systemd", boolEnv("SKYDNS_SYSTEMD", false), "bind to socket(s) activated by systemd (ignore -addr)")

	// Version
	flag.BoolVar(&config.Version, "version", false, "Print the version and exit.")

	// TTl
	// Minttl
	flag.StringVar(&config.Hostmaster, "hostmaster", "hostmaster@skydns.local.", "hostmaster email address to use")
	flag.IntVar(&config.SCache, "scache", server.SCacheCapacity, "capacity of the signature cache")
	flag.IntVar(&config.RCache, "rcache", 0, "capacity of the response cache") // default to 0 for now
	flag.IntVar(&config.RCacheTtl, "rcache-ttl", server.RCacheTtl, "TTL of the response cache")

	// Ndots
	flag.IntVar(&config.Ndots, "ndots", intEnv("SKYDNS_NDOTS", server.Ndots), "How many labels a name should have before we allow forwarding")

	flag.StringVar(&msg.PathPrefix, "path-prefix", env("SKYDNS_PATH_PREFIX", "skydns"), "backend(etcd) path prefix, default: skydns")
}

func main() {
	p := &program{}
	if err := p.Run(); err != nil {
		log.Logger().Error("start server failed %s", err.Error())
		os.Exit(-1)
	}
}


func validateHostPort(hostPort string) error {
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return err
	}
	if ip := net.ParseIP(host); ip == nil {
		return fmt.Errorf("bad IP address: %s", host)
	}

	if p, _ := strconv.Atoi(port); p < 1 || p > 65535 {
		return fmt.Errorf("bad port number %s", port)
	}
	return nil
}

type program struct {
	dns *server.Server
	api *api.HttpServer
}

func (p *program) Run() error {
	flag.Parse()
	log.Init()
	if config.Version {
		fmt.Printf("skydns server version: %s\n", server.Version)
		os.Exit(0)
	}

	if err := server.NewLru(); err != nil {
		fmt.Printf("skydns newlru failed: %s\n", err)
		os.Exit(-1)
	}
	anticorruption.Init(config.IpDbFile)
	if nameserver != "" {
		for _, hostPort := range strings.Split(nameserver, ",") {
			if err := validateHostPort(hostPort); err != nil {
				log.Logger().Warn("skydns: nameserver is invalid: %s", err)
			}
			config.Nameservers = append(config.Nameservers, hostPort)
		}
	}
	if err := validateHostPort(config.DnsAddr); err != nil {
		log.Logger().Error("skydns: addr is invalid: %s", err)
	}

	if err := validateHostPort(config.ApiAddr); err != nil {
		log.Logger().Error("api: addr is invalid: %s", err)
	}

	if err := server.SetDefaults(config); err != nil {
		log.Logger().Warn("skydns: defaults could not be set from /etc/resolv.conf: %v", err)
	}
	if config.Local != "" {
		config.Local = dns.Fqdn(config.Local)
	}
	var backend server.Backend
	db, err := backends.NewDb(&backends.DbConfig{FullFileName: config.DbFile})
	if err != nil {
		log.Logger().Error("skydns: new Bucket failed %s", err)
		return err
	}
	backend = backends.NewBackend(ctx, db, &backends.Config{Ttl: config.Ttl, Priority: config.Priority})
	s := server.New(backend, db, config)
	if stub {
		s.UpdateStubZones()
	}
	if err := metrics.Metrics(); err != nil {
		log.Logger().Warn("skydns: %s", err)
	} else {
		log.Logger().Info("skydns: metrics enabled on :%s%s", metrics.Port, metrics.Path)
	}

	api := api.NewApiServer(config, backend, db)
	go func() {
		if err := api.Run(); err != nil {
			log.Logger().Warn("api: %s", err)
		}
	}()

	if err := s.Run(); err != nil {
		log.Logger().Warn("skydns: %s", err)
	}
	p.dns = s
	return nil
}

func (p *program) Stop() error {
	return log.Logger().Close()
}