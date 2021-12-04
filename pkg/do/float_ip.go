package do

import (
	"context"
	"github.com/digitalocean/godo"
	"log"
)

func (client *Client) getFloatIpList() (ips []godo.FloatingIP, err error) {
	ctx := context.TODO()

	//TODO: add pagination
	ips, _, err = client.DoClient.FloatingIPs.List(ctx, &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	})

	return
}

// GetAvailableEgressIps fetch info from DO about tracked floating ips
func (client *Client) GetAvailableEgressIps(trackedIps map[string]bool) (doIps []godo.FloatingIP, err error) {
	floatingIps, err := client.getFloatIpList()
	if err != nil {
		return
	}

	for _, ip := range floatingIps {
		if trackedIps[ip.IP] {
			doIps = append(doIps, ip)
			log.Printf("%s allowed for use as egress IP", ip.IP)
		}
	}

	return
}

func (client *Client) AssignFloatIpToDroplet(doDropletId int, ip string) (err error) {
	action, _, err := client.DoClient.FloatingIPActions.Assign(context.TODO(), ip, doDropletId)
	if err != nil {
		return
	}

	log.Printf("IP %s assigned to %d, actionId: %d", ip, doDropletId, action.ID)

	return
}
