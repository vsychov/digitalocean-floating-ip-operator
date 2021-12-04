package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestWithValidEnv(t *testing.T) {
	_ = os.Setenv("APP_EGRESS_POOL_NAME", "egress-pool-name")
	_ = os.Setenv("APP_EGRESS_ALLOWED_IPS", "127.0.0.1,127.0.0.2")
	_ = os.Setenv("DIGITALOCEAN_ACCESS_TOKEN", "do-token-ololo")
	_ = os.Setenv("ROUTING_SERVICE_ACCOUNT_NAME", "default")
	_ = os.Setenv("ROUTING_JOB_NAMESPACE", "default")

	config, err := CreateFromEnv()

	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, string(config.TargetPool), "egress-pool-name")
	assert.Equal(t, string(config.DoToken), "do-token-ololo")
	assert.Equal(t, string(config.RoutingJob.Namespace), "default")
	assert.Equal(t, string(config.RoutingJob.ServiceAccountName), "default")
	assert.True(t, config.AllowedEgressIps["127.0.0.1"])
	assert.True(t, config.AllowedEgressIps["127.0.0.2"])
}

func TestWithMissingPoolName(t *testing.T) {
	_ = os.Setenv("APP_EGRESS_POOL_NAME", "")
	_ = os.Setenv("APP_EGRESS_ALLOWED_IPS", "127.0.0.1,127.0.0.2")
	_ = os.Setenv("DIGITALOCEAN_ACCESS_TOKEN", "do-token-ololo")
	_ = os.Setenv("ROUTING_SERVICE_ACCOUNT_NAME", "default")
	_ = os.Setenv("ROUTING_JOB_NAMESPACE", "default")

	_, err := CreateFromEnv()

	assert.NotNil(t, err)
}

func TestWithMissingIps(t *testing.T) {
	_ = os.Setenv("APP_EGRESS_POOL_NAME", "egress-pool-name")
	_ = os.Setenv("APP_EGRESS_ALLOWED_IPS", "")
	_ = os.Setenv("DIGITALOCEAN_ACCESS_TOKEN", "do-token-ololo")
	_ = os.Setenv("ROUTING_SERVICE_ACCOUNT_NAME", "default")
	_ = os.Setenv("ROUTING_JOB_NAMESPACE", "default")

	_, err := CreateFromEnv()

	assert.NotNil(t, err)
}

func TestWithMissingDoToken(t *testing.T) {
	_ = os.Setenv("APP_EGRESS_POOL_NAME", "egress-pool-name")
	_ = os.Setenv("APP_EGRESS_ALLOWED_IPS", "127.0.0.1")
	_ = os.Setenv("DIGITALOCEAN_ACCESS_TOKEN", "")
	_ = os.Setenv("ROUTING_SERVICE_ACCOUNT_NAME", "default")
	_ = os.Setenv("ROUTING_JOB_NAMESPACE", "default")

	_, err := CreateFromEnv()

	assert.NotNil(t, err)
}

func TestWithMissingRoutingServiceAccountName(t *testing.T) {
	_ = os.Setenv("APP_EGRESS_POOL_NAME", "egress-pool-name")
	_ = os.Setenv("APP_EGRESS_ALLOWED_IPS", "127.0.0.1")
	_ = os.Setenv("DIGITALOCEAN_ACCESS_TOKEN", "do-token-ololo")
	_ = os.Setenv("ROUTING_SERVICE_ACCOUNT_NAME", "")
	_ = os.Setenv("ROUTING_JOB_NAMESPACE", "default")

	_, err := CreateFromEnv()

	assert.NotNil(t, err)
}

func TestWithMissingJobNamespace(t *testing.T) {
	_ = os.Setenv("APP_EGRESS_POOL_NAME", "egress-pool-name")
	_ = os.Setenv("APP_EGRESS_ALLOWED_IPS", "127.0.0.1")
	_ = os.Setenv("DIGITALOCEAN_ACCESS_TOKEN", "do-token-ololo")
	_ = os.Setenv("ROUTING_SERVICE_ACCOUNT_NAME", "default")
	_ = os.Setenv("ROUTING_JOB_NAMESPACE", "")

	_, err := CreateFromEnv()

	assert.NotNil(t, err)
}

func TestWithWrongIps(t *testing.T) {
	_ = os.Setenv("APP_EGRESS_POOL_NAME", "egress-pool-name")
	_ = os.Setenv("APP_EGRESS_ALLOWED_IPS", "327.0.0.1")
	_ = os.Setenv("DIGITALOCEAN_ACCESS_TOKEN", "do-token-ololo")
	_ = os.Setenv("ROUTING_SERVICE_ACCOUNT_NAME", "default")
	_ = os.Setenv("ROUTING_JOB_NAMESPACE", "default")

	_, err := CreateFromEnv()

	assert.NotNil(t, err)
}

func TestWithIpv6(t *testing.T) {
	_ = os.Setenv("APP_EGRESS_POOL_NAME", "egress-pool-name")
	_ = os.Setenv("APP_EGRESS_ALLOWED_IPS", "2a00:1450:4016:80b::200e")
	_ = os.Setenv("DIGITALOCEAN_ACCESS_TOKEN", "do-token-ololo")
	_ = os.Setenv("ROUTING_SERVICE_ACCOUNT_NAME", "default")
	_ = os.Setenv("ROUTING_JOB_NAMESPACE", "default")

	_, err := CreateFromEnv()

	assert.NotNil(t, err)
}
