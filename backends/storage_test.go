package backends

import (
	"fmt"
	"go.etcd.io/bbolt"
	"testing"
)

func Test_Storage(t *testing.T)  {
	db, err := bbolt.Open("my.db", 0600, nil)
	if err != nil {
		t.Error(err.Error())
	}
	db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("answer"), []byte("42"))
		return err
	})

	db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		v := b.Get([]byte("answer"))
		fmt.Printf("The answer is: %s\n", v)
		return nil
	})
}
