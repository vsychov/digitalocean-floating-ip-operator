# :warning: WARNING
This is in beta test, not tested in production, and can cause any issues.

[![Build Status](https://github.com/vsychov/digitalocean-floating-ip-operator/actions/workflows/ci.yml/badge.svg)](https://github.com/vsychov/digitalocean-floating-ip-operator/actions)
[![codecov](https://codecov.io/gh/vsychov/digitalocean-floating-ip-operator/branch/master/graph/badge.svg?token=7V853A3LYA)](https://codecov.io/gh/vsychov/digitalocean-floating-ip-operator)
[![Go Reference](https://pkg.go.dev/badge/github.com/vsychov/digitalocean-floating-ip-operator.svg)](https://pkg.go.dev/github.com/vsychov/digitalocean-floating-ip-operator)
[![Go Report Card](https://goreportcard.com/badge/github.com/vsychov/digitalocean-floating-ip-operator)](https://goreportcard.com/report/github.com/vsychov/digitalocean-floating-ip-operator)
[![Docker Repository on Quay](https://quay.io/repository/vsychov/digitalocean-floating-ip-operator/status "Docker Repository on Quay")](https://quay.io/repository/vsychov/digitalocean-floating-ip-operator)
---

# What is it?
This is K8S operator, that may be used for DigitalOcean managed Kubernetes, for provide static IP for nodes, it's using [floating IP](https://docs.digitalocean.com/products/networking/floating-ips/) for that.

#How it works?
Daemon starts in k8s and watch for all changes in selected pool, 
if new node added, and we have avalible float ip, 
it will be assigned to node, after that, special job will be launched on that node, 
that will override default traffic routing rule, and forward all outgoing traffic over floating ip.

E.g. it can be used together with istio, for create egress-gateway with static IP address. 

# Avalible ENV variables

| ENV variable name         | Description                                                                                                                                                                             |
|---------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| APP_EGRESS_POOL_NAME      | Name of the pool with ingress, used for select nodes to assign float-ip, search by label doks.digitalocean.com/node-pool=$APP_EGRESS_POOL_NAME                                          |
| DIGITALOCEAN_ACCESS_TOKEN | Access token, for digitalocean api, used for assign float ip to k8s node                                                                                                                |
| NODES_WATCH_TIMEOUT       | K8S api watch timeout, default is 1800                                                                                                                                                  |
| APP_EGRESS_ALLOWED_IPS    | Comma separated [floating IP](https://docs.digitalocean.com/products/networking/floating-ips/) list, that will be used by operator instance, without spaces, e.g. `127.0.0.1,127.0.0.2` |
| ROUTING_JOB_NAMESPACE     | Namespace for routing job, routing job will be started on each node of selected pool, and will add routing rules on k8s node, for made all egress traffic over floating ip              |
| ROUTING_SERVICE_ACCOUNT_NAME | Service account name for routing job & operator, should be splited in future                                                                                                            |

# Development note

1. start proxy

```
./start-proxy.sh
```

2. export env

```
export KUBERNETES_SERVICE_HOST=127.0.0.1
export KUBERNETES_SERVICE_PORT=8080
```

3. add token's & cert:
   or set KUBERNETES_ROOT_CA_FILE & KUBERNETES_TOKEN_FILE env

```
sudo mkdir -p /var/run/secrets/kubernetes.io/serviceaccount/
sudo cat token.json > /var/run/secrets/kubernetes.io/serviceaccount/token
sudo cat ca.crt > /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
```

# Tips & Tricks

### curl to api from pod:

```
TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
curl https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT/api/ -k -H "Authorization: Bearer $TOKEN"
```

#### list nodes

```
curl https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT/api/v1/nodes -k -H "Authorization: Bearer $TOKEN"
```

