package main

import (
	"github.com/vsychov/digitalocean-floating-ip-operator/pkg/config"
	"github.com/vsychov/digitalocean-floating-ip-operator/pkg/do"
	"github.com/vsychov/digitalocean-floating-ip-operator/pkg/k8s"
	"github.com/vsychov/digitalocean-floating-ip-operator/pkg/operator"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(config.CreateFromEnv),
		fx.Provide(do.NewDoClient),
		fx.Provide(k8s.NewInstance),
		fx.Invoke(operator.Handle),
	).Run()
}
