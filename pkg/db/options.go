package db

import (
	"os"
	"os/exec"
	"path/filepath"
)

type Option interface {
	Apply(*Options)
}

type Options struct {
	Type     string
	Address  string
	Username string
	Password string
	Database string // 数据库名
}

func NewOptions() (o *Options) {
	address := "./"
	file, err := exec.LookPath(os.Args[0])
	if err == nil {
		address, _ = filepath.Abs(file)
		address = filepath.Dir(address)
	}
	return &Options{
		Type:    Bolt,
		Address: address,
	}
}

func (o *Options) Apply(x *Options) {
	if x != o {
		if o.Type != "" {
			(*x).Type = o.Type
		}
		if o.Address != "" {
			(*x).Address = o.Address
		}

		(*x).Username = o.Username
		(*x).Password = o.Password
		(*x).Database = o.Database
	}
}

func (o *Options) SetUsername(username string) *Options {
	o.Username = username
	return o
}

func (o *Options) SetPassword(password string) *Options {
	o.Password = password
	return o
}

func (o *Options) SetDatabase(database string) *Options {
	o.Database = database
	return o
}

func (o *Options) SetType(t string) *Options {
	o.Type = t
	return o
}

func (o *Options) SetAddress(address string) *Options {
	o.Address = address
	return o
}
