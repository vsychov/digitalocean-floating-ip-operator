package do

import (
	"digitalocean-floating-ip-operator/pkg/config"
	"github.com/digitalocean/godo"
)

type Client struct {
	DoClient *godo.Client
}

func NewDoClient(config config.Config) Client {
	return Client{DoClient: godo.NewFromToken(string(config.DoToken))}
}
