package socks

type Option func(opt *option)

type option struct {
	addr string
}

func (o *option) apply(s *Proxy) {
	if o.addr != "" {
		s.addr = o.addr
	}
}

func WithProxyAddr(addr string) Option {
	return func(opt *option) {
		opt.addr = addr
	}
}
