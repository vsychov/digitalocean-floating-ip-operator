package operator

import (
	"encoding/json"
	k8s_bridge_do "float-ip-do-k8s/pkg/k8s-bridge-do"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
)

func (op *operator) newOrModify(node v1.Node) (err error) {
	//    - check, if it ready for egress ("egress-ready=false" should be seted if it not ready)
	//    - try to find free floating ip (from operator list)
	//    - if free ip not avalible - skip node
	//    - assign ip to node via DO API

	//    - get default internal gateway from metadata (curl http://169.254.169.254/metadata/v1/interfaces/public/0/anchor_ipv4/gateway)
	//    - update node default route (ip route replace default via 10.19.0.1 dev eth0)
	//    - verify external ip changed (curl to somewhere external, e.g. do api)
	//    - mark node as ready for egress ("egress-ready=true")
	log.Printf("Node: %s (name: %s)", node.Spec.ProviderID, node.Name)

	isHaveAllowedFloatIp, ip, err := k8s_bridge_do.IsNodeHaveAllowedFloatingIp(op.Config.AllowedEgressIps, &op.DoClient, &node)
	if err != nil {
		log.Printf("Skip event, API error received: %s", err.Error())
		return
	}

	if isHaveAllowedFloatIp {
		log.Printf("Node have valid floating ip assigned: %s", ip.IP)
		op.K8s.AddRouteForFloatIp(&node, ip.IP)
		return
	}

	log.Printf("Node have no assigned float IP's, let's assign free IP to this node")
	op.K8s.SetEgressReadyLabel(&node, false)

	ips, err := op.DoClient.GetAvailableEgressIps(op.Config.AllowedEgressIps)
	if err != nil {
		return
	}

	for _, ip := range ips {
		//ip not assigned to Droplet, assign it to current Droplet
		if ip.Droplet == nil {
			dropletId, err := k8s_bridge_do.GetDoDropletId(&node)
			if err != nil {
				return err
			}

			log.Printf("IP: %s not used, assign it %s (%d)", ip.IP, node.Name, dropletId)
			err = op.DoClient.AssignFloatIpToDroplet(dropletId, ip.IP)
			if err != nil {
				log.Printf("Skip event, API error received: %s", err.Error())
				return err
			}

			op.K8s.AddRouteForFloatIp(&node, ip.IP)
			return nil
		}
	}

	//TODO: trigrer alert, no avalible IP address
	log.Printf("ALARM! No free IP for new node!!!")
	return
}

func (op *operator) handleNodeEvent(event watch.Event) error {
	log.Printf("Node Event Type %v", event.Type)

	switch event.Type {
	case watch.Added, watch.Deleted, watch.Modified:
		node, err := op.unmarshalNode(event.Object)
		if err != nil {
			return err
		}

		return op.routeNodeEvent(event.Type, node)
	default:
		log.Printf("Skip %s event", event.Type)
	}

	return nil
}

func (op *operator) routeNodeEvent(eventType watch.EventType, node v1.Node) error {
	switch eventType {
	case watch.Added:
		err := op.newOrModify(node)
		if err != nil {
			return err
		}
	case watch.Modified:
		err := op.newOrModify(node)
		if err != nil {
			return err
		}
	case watch.Deleted:
		op.removeNode(node)
	default:
		return fmt.Errorf("Wrong event type: %s", eventType)
	}

	return nil
}

func (op *operator) removeNode(node v1.Node) {
	//    - make sure, all egress-gateways evicted from node
	//    - if not - force evict it
	//    - unassign floating ip from node via DOSK api
	//TODO: implement node removal
	log.Printf("Node Removed: %s, but nothing happens", node.Name)
}

func (op *operator) unmarshalNode(obj interface{}) (node v1.Node, err error) {
	b, err := json.Marshal(obj)
	err = json.Unmarshal(b, &node)
	return
}
