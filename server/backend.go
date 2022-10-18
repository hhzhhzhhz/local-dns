// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package server

import (
	"context"
	"github.com/hhzhhzhhz/local-dns/msg"
)

type Backend interface {
	HasSynced() bool
	Records(name string, exact bool) ([]msg.Service, error)
	ReverseRecord(name string) (*msg.Service, error)
	SaveRecords(ctx context.Context, name string, rds []msg.Service) error
	RemoveRecords(ctx context.Context, domain []string) error
	AllRecord(ctx context.Context) ([]msg.Service, error)
}