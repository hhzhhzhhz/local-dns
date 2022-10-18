// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package server

import (
	"fmt"
	"github.com/hhzhhzhhz/local-dns/log"
	"github.com/miekg/dns"
	"strings"
)

// printf calls log.Printf with the parameters given.
func logf() func(format string, v ...interface{}) {
	return log.Logger().Info
}

// fatalf calls log.Fatalf with the parameters given.
func fatalf() func(format string, v ...interface{}) {
	return log.Logger().Info
}


func request_log(w dns.ResponseWriter, req *dns.Msg, resp *dns.Msg) {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("received DNS Request from %q question=[", w.RemoteAddr()))
	for _, r := range req.Question {
		buf.WriteString(fmt.Sprintf("%s %s %s;", r.Name, dns.Type(r.Qtype).String(), dns.Class(r.Qclass).String()))
	}
	buf.WriteString("] answer=[")
	for _, r := range resp.Answer {
		buf.WriteString(r.String()+ ";")
	}
	buf.WriteString("] ns=[")
	for _, r := range resp.Ns {
		buf.WriteString(r.String()+ ";")
	}
	buf.WriteString("] extra=[")
	for _, r := range resp.Extra {
		buf.WriteString(r.String()+ ";")
	}
	buf.WriteString("]")
	log.Logger().Info(buf.String())
}


//func uniq(w dns.ResponseWriter, req *dns.Msg, resp *dns.Msg) {
//	var name, qtype string
//	for _, r := range req.Question {
//		name = r.Name
//		qtype = dns.Type(r.Qtype).String()
//	}
//	var item *msg.DnsItem
//	v, ok := Lru().Get(name)
//	if !ok {
//		count := int64(0)
//		item = &msg.DnsItem{Dns: map[string]cache.Set{}, Count: &count, Domain: name, Type: qtype}
//	} else {
//		item = v.(*msg.DnsItem)
//	}
//	*item.Count++
//	for _, r := range resp.Answer {
//		rcd, ok := item.Dns[dns.Type(r.Header().Rrtype).String()]
//		if ok {
//			rcd.Add(parseAnswer(r))
//		} else {
//			rcd = cache.NewSet()
//			rcd.Add(parseAnswer(r))
//			item.Dns[dns.Type(r.Header().Rrtype).String()] = rcd
//		}
//	}
//
//	for _, r := range resp.Ns {
//		rcd, ok := item.Dns[dns.Type(r.Header().Rrtype).String()]
//		if ok {
//			rcd.Add(parseNS(r))
//		} else {
//			rcd = cache.NewSet()
//			rcd.Add(parseNS(r))
//			item.Dns[dns.Type(r.Header().Rrtype).String()] = rcd
//		}
//	}
//
//	for _, r := range resp.Extra {
//		rcd, ok := item.Dns[dns.Type(r.Header().Rrtype).String()]
//		if ok {
//			rcd.Add(parseExt(r))
//		} else {
//			rcd = cache.NewSet()
//			rcd.Add(parseExt(r))
//			item.Dns[dns.Type(r.Header().Rrtype).String()] = rcd
//		}
//	}
//
//	Lru().Add(name, item)
//}

func parseAnswer(r dns.RR) string {
	switch r.(type) {
	case *dns.A:
		a := r.(*dns.A)
		return a.A.String()
	case *dns.AAAA:
		a := r.(*dns.AAAA)
		return a.AAAA.String()
	case *dns.SOA:
		a := r.(*dns.SOA)
		return a.Ns
	case *dns.CNAME:
		a := r.(*dns.CNAME)
		return a.Hdr.Name
	case *dns.NS:
		a := r.(*dns.NS)
		return a.Ns
	default:
		return r.String()
	}
}

func parseNS(r dns.RR) string {
	return parseAnswer(r)
}

func parseExt(r dns.RR) string {
	return parseAnswer(r)
}

