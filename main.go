package main

import (
	"digitalocean-floating-ip-operator/pkg/config"
	"digitalocean-floating-ip-operator/pkg/do"
	"digitalocean-floating-ip-operator/pkg/k8s"
	"digitalocean-floating-ip-operator/pkg/operator"
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
