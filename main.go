package main

import (
	"float-ip-do-k8s/pkg/config"
	"float-ip-do-k8s/pkg/do"
	"float-ip-do-k8s/pkg/k8s"
	"float-ip-do-k8s/pkg/operator"
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
