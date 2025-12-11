package db

import "errors"

var (
	ErrUnsupportedDBType = errors.New("unsupported db type")
	ErrNotFound          = errors.New("not found")
)

// 数据库类型枚举
const (
	MongoDB       = "mongodb"
	Redis         = "redis"
	Memory        = "memory"
	InfluxDB      = "influxdb"
	MySQL         = "mysql"
	Sqlite3       = "sqlite3"
	ElasticSearch = "elasticsearch"
	Bolt          = "bolt"
)

type IDB interface {
	Close() error

	ITask
	ILog
}
