package k8s_bridge_do

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/vsychov/digitalocean-floating-ip-operator/pkg/do"
	v1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

// GetDoDropletId return DO specific droplet ID
func GetDoDropletId(node *v1.Node) (doDropletId int, err error) {
	doDropletId, err = strconv.Atoi(strings.Replace(node.Spec.ProviderID, "digitalocean://", "", 1))
	if err != nil {
		return
	}

	return
}

// IsNodeHaveAllowedFloatingIp check that specific node have assigned float ip
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
