package do

import (
	"float-ip-do-k8s/pkg/config"
	"github.com/digitalocean/godo"
)

type Client struct {
	DoClient *godo.Client
}

func NewDoClient(config config.Config) Client {
	return Client{DoClient: godo.NewFromToken(string(config.DoToken))}
}
