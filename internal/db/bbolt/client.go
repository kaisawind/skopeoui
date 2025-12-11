package bbolt

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"

	pkgdb "github.com/kaisawind/skopeoui/pkg/db"
	"go.etcd.io/bbolt"
)

// 数据库定义
const (
	dbName = "bbolt.db"
	dbTask = "task"
	dbLog  = "log"
)

// DB redis database
type DB struct {
	opts   *pkgdb.Options
	client *bbolt.DB
}

// NewDB 创建redis连接
func NewDB(ops *pkgdb.Options) (*DB, error) {
	if ops.Address == "" {
		return nil, fmt.Errorf("db address is empty")
	}
	if ops.Database == "" {
		ops.Database = dbName
	}
	dbFile := path.Join(ops.Address, ops.Database)
	client, err := bbolt.Open(dbFile, os.ModePerm, bbolt.DefaultOptions)
	if err != nil {
		return nil, err
	}
	db := &DB{
		opts:   ops,
		client: client,
	}
	err = db.createBuckets()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) Close() (err error) {
	if db.client != nil {
		db.client.Close()
	}
	return
}

func (db *DB) createBuckets() (err error) {
	list := []string{dbTask, dbLog}
	err = db.client.Batch(func(tx *bbolt.Tx) error {
		for _, v := range list {
			b := tx.Bucket([]byte(v))
			if b != nil {
				continue
			}
			b, err := tx.CreateBucketIfNotExists([]byte(v))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

func (db *DB) CTask(tx *bbolt.Tx) *bbolt.Bucket {
	return tx.Bucket([]byte(dbTask))
}

func (db *DB) CLog(tx *bbolt.Tx) *bbolt.Bucket {
	return tx.Bucket([]byte(dbLog))
}

// itob returns an 8-byte big endian representation of v.
func itob[T int | int32 | uint64](v T) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
