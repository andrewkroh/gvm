package httpproxy

import "net/http"

type Proxy interface {
	NewTransport() *http.Transport
}
