package db

import (
	"github.com/kaisawind/skopeoui/internal/db/bbolt"
	pkgdb "github.com/kaisawind/skopeoui/pkg/db"
)

func New(opts ...pkgdb.Option) (pkgdb.IDB, error) {
	o := pkgdb.NewOptions()
	for _, opt := range opts {
		if opt != nil {
			opt.Apply(o)
		}
	}
	switch o.Type {
	case pkgdb.Bolt:
		return bbolt.NewDB(o)
	case pkgdb.MongoDB:
		return nil, pkgdb.ErrUnsupportedDBType
	case pkgdb.Redis:
		return nil, pkgdb.ErrUnsupportedDBType
	case pkgdb.MySQL:
		return nil, pkgdb.ErrUnsupportedDBType
	case pkgdb.Sqlite3:
		return nil, pkgdb.ErrUnsupportedDBType
	case pkgdb.Memory:
		return nil, pkgdb.ErrUnsupportedDBType
	default:
		return nil, pkgdb.ErrUnsupportedDBType
	}
}
