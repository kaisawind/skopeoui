package bbolt

import (
	"context"
	"encoding/json/v2"
	"math"

	pkgdb "github.com/kaisawind/skopeoui/pkg/db"
	"github.com/kaisawind/skopeoui/pkg/pb"
	"go.etcd.io/bbolt"
)

func (db *DB) CreateLog(ctx context.Context, in *pb.Log) (err error) {
	err = db.client.Update(func(tx *bbolt.Tx) error {
		id, _ := db.CLog(tx).NextSequence()
		in.Id = int32(id)
		buf, _ := json.Marshal(in)
		return db.CLog(tx).Put(itob(id), buf)
	})
	if err != nil {
		return
	}
	return
}

func (db *DB) DeleteLog(ctx context.Context, id int32) (err error) {
	err = db.client.Update(func(tx *bbolt.Tx) error {
		return db.CLog(tx).Delete(itob(id))
	})
	if err != nil {
		return
	}
	return
}

func (db *DB) DeleteLogByTaskId(ctx context.Context, taskId int32) (err error) {
	err = db.client.Update(func(tx *bbolt.Tx) error {
		return db.CLog(tx).ForEach(func(k, v []byte) error {
			log := &pb.Log{}
			err = json.Unmarshal(v, log)
			if err != nil {
				return err
			}
			if log.TaskId == taskId {
				return db.CLog(tx).Delete(k)
			}
			return nil
		})
	})
	if err != nil {
		return
	}
	return
}

func (db *DB) UpdateLog(ctx context.Context, in *pb.Log) (err error) {
	err = db.client.Update(func(tx *bbolt.Tx) error {
		buf, _ := json.Marshal(in)
		return db.CLog(tx).Put(itob(in.Id), buf)
	})
	if err != nil {
		return
	}
	return
}

func (db *DB) GetLog(ctx context.Context, id int32) (out *pb.Log, err error) {
	err = db.client.View(func(tx *bbolt.Tx) error {
		v := db.CLog(tx).Get(itob(id))
		if v == nil {
			return pkgdb.ErrNotFound
		}
		err = json.Unmarshal(v, out)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

func (db *DB) ListLog(ctx context.Context, skip, limit int32) (count int32, out []*pb.Log, err error) {
	if skip < 0 {
		skip = 0
	}
	if limit <= 0 {
		limit = math.MaxInt32
	}
	err = db.client.View(func(tx *bbolt.Tx) error {
		c := db.CLog(tx).Cursor()
		i := int32(0)
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if i >= skip && i < skip+limit {
				t := &pb.Log{}
				err = json.Unmarshal(v, t)
				if err != nil {
					continue
				}
				out = append(out, t)
			}
			i++
		}
		count = i
		return nil
	})
	if err != nil {
		return
	}
	return
}

func (db *DB) ListLogByTaskId(ctx context.Context, taskId int32, skip, limit int32) (count int32, out []*pb.Log, err error) {
	if skip < 0 {
		skip = 0
	}
	if limit <= 0 {
		limit = math.MaxInt32
	}
	err = db.client.View(func(tx *bbolt.Tx) error {
		c := db.CLog(tx).Cursor()
		i := int32(0)
		for k, v := c.First(); k != nil; k, v = c.Next() {
			log := &pb.Log{}
			err = json.Unmarshal(v, log)
			if err != nil {
				continue
			}
			if log.TaskId == taskId {
				if i >= skip && i < skip+limit {
					out = append(out, log)
				}
				i++
			}
		}
		count = i
		return nil
	})
	if err != nil {
		return
	}
	return
}
