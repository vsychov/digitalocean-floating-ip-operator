package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config is application configuration
type Config struct {
	TargetPool       nonEmptyStringDecoder   `envconfig:"APP_EGRESS_POOL_NAME"`
	AllowedEgressIps allowedIngressIpDecoder `envconfig:"APP_EGRESS_ALLOWED_IPS"`
	DoToken          nonEmptyStringDecoder   `envconfig:"DIGITALOCEAN_ACCESS_TOKEN"`
	WatchTimeout     int64                   `envconfig:"NODES_WATCH_TIMEOUT" default:"1800"`
	RoutingJob       routingJobConfig
}

type routingJobConfig struct {
	Namespace          nonEmptyStringDecoder `envconfig:"ROUTING_JOB_NAMESPACE"`
	ServiceAccountName nonEmptyStringDecoder `envconfig:"ROUTING_SERVICE_ACCOUNT_NAME"`
}

// CreateFromEnv create new config instance from env variables
func CreateFromEnv() (conf Config, err error) {
	err = envconfig.Process("", &conf)
	if err != nil {
		return
	}

	return
}
