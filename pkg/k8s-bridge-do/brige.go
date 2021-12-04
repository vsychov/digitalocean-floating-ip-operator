package k8s_bridge_do

import (
	"context"
	"digitalocean-floating-ip-operator/pkg/do"
	"github.com/digitalocean/godo"
	v1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

func GetDoDropletId(node *v1.Node) (doDropletId int, err error) {
	doDropletId, err = strconv.Atoi(strings.Replace(node.Spec.ProviderID, "digitalocean://", "", 1))
	if err != nil {
		return
	}

	return
}

func IsNodeHaveAllowedFloatingIp(trackedIps map[string]bool, doClient *do.Client, node *v1.Node) (result bool, ip godo.FloatingIP, err error) {
	doDropletId, err := GetDoDropletId(node)
	result = false

	if err != nil {
		return
	}

	allowedEgressIps, err := doClient.GetAvailableEgressIps(trackedIps)
	if err != nil {
		return
	}

	droplet, _, err := doClient.DoClient.Droplets.Get(context.TODO(), doDropletId)
	if err != nil {
		return
	}

	for _, v4Ip := range droplet.Networks.V4 {
		for _, floatIp := range allowedEgressIps {
			if v4Ip.IPAddress == floatIp.IP {
				ip = floatIp
				result = true

				return
			}
		}
	}

	return
}
