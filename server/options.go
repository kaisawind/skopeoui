package server

import (
	"encoding/json/v2"
	"time"

	"github.com/kaisawind/skopeoui/pkg/configs"
)

type Options struct {
	HttpAddress  string
	DbType       string
	DbAddress    string
	DbUsername   string
	DbPassword   string
	DbDatabase   string // 数据库名
	DbExpiration time.Duration
}

func NewOptions() *Options {
	return &Options{
		HttpAddress:  configs.DefaultHttp.Address,
		DbType:       configs.DefaultDB.Type,
		DbAddress:    configs.DefaultDB.Address,
		DbUsername:   configs.DefaultDB.Username,
		DbPassword:   configs.DefaultDB.Password,
		DbDatabase:   configs.DefaultDB.Database,
		DbExpiration: configs.DefaultDB.Expiration,
	}
}

func (o *Options) Apply(x *Options) {
	if x != o {
		if o.HttpAddress != "" {
			(*x).HttpAddress = o.HttpAddress
		}
		if o.DbAddress != "" {
			(*x).DbAddress = o.DbAddress
		}

		(*x).DbUsername = o.DbUsername
		(*x).DbPassword = o.DbPassword
		(*x).DbDatabase = o.DbDatabase
		if o.DbExpiration != 0 {
			(*x).DbExpiration = o.DbExpiration
		}
	}
}

func (opts *Options) ApplyDBConfig(cfg configs.DB) *Options {
	opts.SetDbType(cfg.Type)
	opts.SetDbAddress(cfg.Address)
	opts.SetDbUsername(cfg.Username)
	opts.SetDbPassword(cfg.Password)
	opts.SetDbDatabase(cfg.Database)
	opts.SetDbExpiration(cfg.Expiration)
	return opts
}

func (opts *Options) ApplyHttpConfig(cfg configs.Http) *Options {
	opts.SetHttpAddress(cfg.Address)
	return opts
}

func (opts *Options) SetHttpAddress(addr string) *Options {
	opts.HttpAddress = addr
	return opts
}

func (opts *Options) SetDbType(dbType string) *Options {
	opts.DbType = dbType
	return opts
}
func (opts *Options) SetDbAddress(dbAddress string) *Options {
	opts.DbAddress = dbAddress
	return opts
}

func (opts *Options) SetDbUsername(dbUsername string) *Options {
	opts.DbUsername = dbUsername
	return opts
}
func (opts *Options) SetDbPassword(dbPassword string) *Options {
	opts.DbPassword = dbPassword
	return opts
}
func (opts *Options) SetDbDatabase(dbDatabase string) *Options {
	opts.DbDatabase = dbDatabase
	return opts
}
func (opts *Options) SetDbExpiration(dbExpiration time.Duration) *Options {
	opts.DbExpiration = dbExpiration
	return opts
}

func (opts *Options) String() string {
	buff, _ := json.Marshal(opts)
	return string(buff)
}
