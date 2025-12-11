package bbolt

import (
	"context"
	"encoding/json/v2"
	"math"

	pkgdb "github.com/kaisawind/skopeoui/pkg/db"
	"github.com/kaisawind/skopeoui/pkg/pb"
	"go.etcd.io/bbolt"
)

func (db *DB) CreateTask(ctx context.Context, in *pb.Task) (err error) {
	err = db.client.Update(func(tx *bbolt.Tx) error {
		id, _ := db.CTask(tx).NextSequence()
		in.Id = int32(id)
		buf, _ := json.Marshal(in)
		return db.CTask(tx).Put(itob(id), buf)
	})
	if err != nil {
		return
	}
	return
}

func (db *DB) DeleteTask(ctx context.Context, id int32) (err error) {
	err = db.client.Update(func(tx *bbolt.Tx) error {
		return db.CTask(tx).Delete(itob(id))
	})
	if err != nil {
		return
	}
	return
}

func (db *DB) UpdateTask(ctx context.Context, in *pb.Task) (err error) {
	err = db.client.Batch(func(tx *bbolt.Tx) error {
		buf, _ := json.Marshal(in)
		return db.CTask(tx).Put(itob(in.Id), buf)
	})
	if err != nil {
		return
	}
	return
}

func (db *DB) GetTask(ctx context.Context, id int32) (out *pb.Task, err error) {
	out = &pb.Task{}
	err = db.client.View(func(tx *bbolt.Tx) error {
		buf := db.CTask(tx).Get(itob(id))
		if len(buf) == 0 {
			return pkgdb.ErrNotFound
		}
		return json.Unmarshal(buf, out)
	})
	if err != nil {
		return
	}
	return
}

func (db *DB) ListTask(ctx context.Context, skip, limit int32) (count int32, out []*pb.Task, err error) {
	if skip < 0 {
		skip = 0
	}
	if limit <= 0 {
		limit = math.MaxInt32
	}
	err = db.client.View(func(tx *bbolt.Tx) error {
		c := db.CTask(tx).Cursor()
		i := int32(0)
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if i >= skip && i < skip+limit {
				t := &pb.Task{}
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
