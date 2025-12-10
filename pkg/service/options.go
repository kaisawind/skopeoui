package service

type Options struct {
	Address string // 服务监听地址
}

func NewOptions() (o *Options) {
	return &Options{}
}

func (o *Options) Apply(x *Options) {
	if x != o {
		if o.Address != "" {
			(*x).Address = o.Address
		}
	}
}

func (o *Options) SetAddress(address string) *Options {
	o.Address = address
	return o
}

type Option interface {
	Apply(*Options)
}

type funcOption struct {
	f func(*Options)
}

func (fdo *funcOption) Apply(do *Options) {
	fdo.f(do)
}

func newFuncOption(f func(*Options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func ServiceAddressOption(addr string) Option {
	return newFuncOption(func(o *Options) {
		o.Address = addr
	})
}
