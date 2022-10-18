package backends

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hhzhhzhhz/local-dns/msg"
	"github.com/hhzhhzhhz/local-dns/singleflight"
	"strings"
)

func NewBackend(ctx context.Context, storage Storage, config *Config) *backend {
	return &backend{
		ctx: ctx,
		config: config,
		db: storage,
		inflight: &singleflight.Group{},
	}
}

type Config struct {
	Ttl      uint32
	Priority uint16
}

type backend struct {
	ctx context.Context
	db Storage
	config   *Config
	inflight *singleflight.Group

}

func (bd *backend) HasSynced() bool {
	return true
}

func (bd *backend) Records(name string, exact bool) ([]msg.Service, error) {
	path, star := msg.PathWithWildcard(name)
	r, err := bd.get(path, true)
	if err != nil {
		return nil, err
	}
	segments := strings.Split(msg.Path(name), "/")

	return bd.loopNodes(r, segments, star, nil)
}

func (bd *backend) ReverseRecord(name string) (*msg.Service, error) {
	path, star := msg.PathWithWildcard(name)
	if star {
		return nil, fmt.Errorf("reverse can not contain wildcards")
	}

	r, err := bd.get(path, true)
	if err != nil {
		return nil, err
	}

	segments := strings.Split(msg.Path(name), "/")
	records, err := bd.loopNodes(r, segments, star, nil)
	if err != nil {
		return nil, err
	}
	if len(records) != 1 {
		return nil, fmt.Errorf("must be only one service record")
	}
	return &records[0], nil
}

func (bd *backend) SaveRecords(ctx context.Context, name string, rds []msg.Service) error {
	path, _ := msg.PathWithWildcard(name)
	if err := bd.db.Save(ctx,"", path, rds); err != nil {
		return err
	}
	return nil
}

func (bd *backend) RemoveRecords(ctx context.Context, domain []string) error {
	for _, dmn := range domain {
		path, _ := msg.PathWithWildcard(dmn)
		if err := bd.db.Delete(ctx,"",  path); err != nil {
			return err
		}
	}
	return nil
}

func (bd *backend) AllRecord(ctx context.Context) (sx []msg.Service, err error) {
	b, err := bd.inflight.Do("all", func() (interface{}, error) {
		return bd.db.All(context.Background(), "")
	})
	if err != nil {
		return nil, err
	}
	return bd.loopNodes(b.(RangeResponse).Kvs, []string{}, false, nil)
}

func (bd *backend) get(path string, recursive bool) (kv []*KeyValue, err error) {
	b, err := bd.inflight.Do(path, func() (interface{}, error) {
		return bd.db.Get(context.Background(), "", path)
	})
	if err != nil {
		return nil, err
	}

	return b.(RangeResponse).Kvs, nil
}

func (bd *backend) loopNodes(kv []*KeyValue, nameParts []string, star bool, bx map[bareService]bool) (sx []msg.Service, err error) {
	if bx == nil {
		bx = make(map[bareService]bool)
	}
Nodes:
	for _, item := range kv {

		if star {
			s := string(item.Key[:])
			keyParts := strings.Split(s, "/")
			for i, n := range nameParts {
				if i > len(keyParts)-1 {
					continue Nodes
				}
				if n == "*" || n == "any" {
					continue
				}
				if keyParts[i] != n {
					continue Nodes
				}
			}
		}

		var servs []*msg.Service
		if err := json.Unmarshal(item.Value, &servs); err != nil {
			return nil, err
		}
		for _, serv := range servs {
			b := bareService{serv.Host,
				serv.Port,
				serv.Priority,
				serv.Weight,
				serv.Text}

			bx[b] = true
			serv.Key = string(item.Key)
			//TODO: another call (LeaseRequest) for TTL when RPC in etcdv3 is ready
			serv.Ttl = bd.calculateTtl(item, serv)

			if serv.Priority == 0 {
				serv.Priority = int(bd.config.Priority)
			}

			sx = append(sx, *serv)
		}


	}
	return sx, nil
}


func (bd *backend) calculateTtl(kv *KeyValue, serv *msg.Service) uint32 {
	etcdTtl := uint32(10) //TODO: default value for now, should be an rpc call for least request when it becomes available in etcdv3's api

	if etcdTtl == 0 && serv.Ttl == 0 {
		return bd.config.Ttl
	}
	if etcdTtl == 0 {
		return serv.Ttl
	}
	if serv.Ttl == 0 {
		return etcdTtl
	}
	if etcdTtl < serv.Ttl {
		return etcdTtl
	}
	return serv.Ttl
}

type bareService struct {
	Host     string
	Port     int
	Priority int
	Weight   int
	Text     string
}
