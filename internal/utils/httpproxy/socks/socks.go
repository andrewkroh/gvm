package socks

import (
	"net/http"
	"net/url"
)

const (
	defaultProxyAddr = "socks5://127.0.0.1:10808"
)

func NewSocks5Proxy(opts ...Option) *Proxy {
	var (
		o     = &option{}
		socks = &Proxy{
			addr: defaultProxyAddr,
		}
	)

	for _, opt := range opts {
		opt(o)
	}

	o.apply(socks)

	return socks
}

type Proxy struct {
	addr string
}

func (s *Proxy) NewTransport() *http.Transport {
	return &http.Transport{
		Proxy: func(_ *http.Request) (*url.URL, error) {
			return url.Parse(s.addr)
		},
	}
}
