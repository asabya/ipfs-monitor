package client

import (
	"net/http"
	"time"

	"go.uber.org/fx"
)

type HttpClient struct {
	Client *http.Client
	Base   string
}

func NewHttpClient() fx.Option {
	return fx.Provide(func() *HttpClient {
		return &HttpClient{
			Client: &http.Client{
				Transport:     nil,
				CheckRedirect: nil,
				Jar:           nil,
				Timeout:       time.Second * 10,
			},
			Base: "http://localhost:5001/api/v0/",
		}
	})
}
