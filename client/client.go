package client

import (
	"net/http"
	"os"
	"time"

	"go.uber.org/fx"
)

type HttpClient struct {
	Client *http.Client
	Base   string
}

// NewHttpClient creates an http client for
// calling ipfs http api
// reads os.Args[1] for ipfs "HOST:PORT"
// Default is "localhost:5001"
func NewHttpClient() fx.Option {
	return fx.Provide(func() *HttpClient {
		base := "http://localhost:5001"
		if len(os.Args) > 1 {
			base = os.Args[1]
		}
		return &HttpClient{
			Client: &http.Client{
				Transport:     nil,
				CheckRedirect: nil,
				Jar:           nil,
				Timeout:       time.Second * 10,
			},
			Base: base + "/api/v0/",
		}
	})
}
