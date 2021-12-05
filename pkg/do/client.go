package do

import (
	"github.com/digitalocean/godo"
	"github.com/vsychov/digitalocean-floating-ip-operator/pkg/config"
)

type Client struct {
	DoClient *godo.Client
}

func NewDoClient(config config.Config) Client {
	return Client{DoClient: godo.NewFromToken(string(config.DoToken))}
}
