package operator

import (
	"context"
	"github.com/vsychov/digitalocean-floating-ip-operator/pkg/config"
	"github.com/vsychov/digitalocean-floating-ip-operator/pkg/do"
	"github.com/vsychov/digitalocean-floating-ip-operator/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

type operator struct {
	Config   config.Config
	K8s      k8s.K8s
	DoClient do.Client
}

func Handle(config config.Config, doClient do.Client, k8s k8s.K8s) {
	log.Println("Operator started")
	log.Printf("Config: %+v\n", config)

	watchNodesInterface, err := k8s.ClientSet.CoreV1().Nodes().Watch(context.TODO(), metav1.ListOptions{
		LabelSelector:  "doks.digitalocean.com/node-pool=" + string(config.TargetPool),
		Watch:          true,
		TimeoutSeconds: &config.WatchTimeout,
	})

	op := operator{
		Config:   config,
		DoClient: doClient,
		K8s:      k8s,
	}

	if err != nil {
		panic(err.Error())
	}

	for {
		event, ok := <-watchNodesInterface.ResultChan()
		err := op.handleNodeEvent(event)
		if err != nil {
			panic(err)
		}

		if !ok {
			panic("timeout, exiting")
		}
	}
}
