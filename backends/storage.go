package backends

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"time"
)

const (
	defaultBucket = "dns_records"
)

type DbConfig struct {
	FullFileName string
}

type Storage interface {
	Save(ctx context.Context, bucket, key string, value interface{}) error
	Delete(ctx context.Context, bucket, key string) error
	Get(ctx context.Context, bucket, key string) (RangeResponse, error)
	All(ctx context.Context, bucket string) (RangeResponse, error)
}

func NewDb(cfg *DbConfig) (Storage, error) {
	db, err := bbolt.Open(cfg.FullFileName, 0600, &bbolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return nil, err
	}
	return &storage{db: db, cfg: cfg}, nil
}

type storage struct {
	db *bbolt.DB
	cfg *DbConfig
}

func (s *storage) Save(ctx context.Context, bucket, key string, value interface{}) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bt, err := tx.CreateBucketIfNotExists(s.hash(bucket))
		if err != nil {
			return err
		}
		vb, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return bt.Put([]byte(key), vb)
	})
}

func (s *storage) Delete(ctx context.Context, bucket, key string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(s.hash(bucket)).Delete([]byte(key))
	})
}

func (s *storage) Get(ctx context.Context, bucket,  key string) (RangeResponse, error) {
	var b []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		bt := tx.Bucket(s.hash(bucket))
		if bt == nil {
			return fmt.Errorf("bucket %s not found", s.hash(bucket))
		}
		b = bt.Get([]byte(key))
		return nil
	})
	if b == nil {
		return RangeResponse{}, nil
	}
	return RangeResponse{Kvs: []*KeyValue{&KeyValue{Key: []byte(key), Value: b}}}, err
}



func (s *storage) All(ctx context.Context, bucket string) (RangeResponse, error) {
	var kv []*KeyValue
	err := s.db.View(func(tx *bbolt.Tx) error {
		bt := tx.Bucket(s.hash(bucket))
		if bt == nil {
			return fmt.Errorf("bucket %s not found", s.hash(bucket))
		}
		cs := bt.Cursor()
		for k, v := cs.First(); k != nil; k, v = cs.Next() {
			kv = append(kv, &KeyValue{Key: k, Value: v})
		}
		return nil
	})
	if err != nil {
		return RangeResponse{}, err
	}
	return RangeResponse{Kvs: kv}, err
}

func (s *storage) hash(bucket string) []byte {
	if bucket != "" {
		return []byte(bucket)
	}
	return []byte(defaultBucket)
	//h :=fnv.New64()
	//h.Write([]byte(v))
	//var buf = make([]byte, 8)
	//binary.BigEndian.PutUint64(buf, h.Sum64())
	//return buf
}
