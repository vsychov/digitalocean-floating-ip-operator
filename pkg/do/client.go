package do

import (
	"github.com/digitalocean/godo"
	"github.com/vsychov/digitalocean-floating-ip-operator/pkg/config"
)

// Client DigitalOcean client wrapper
type Client struct {
	DoClient *godo.Client
}

// NewDoClient create new instance of Client
func NewDoClient(config config.Config) Client {
	return Client{DoClient: godo.NewFromToken(string(config.DoToken))}
}
