package bridge

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func TestCorrectProviderId(t *testing.T) {
	node := v1.Node{
		Spec: v1.NodeSpec{
			ProviderID: "digitalocean://123",
		},
	}

	dropletId, err := GetDoDropletId(&node)
	assert.Equal(t, 123, dropletId)
	assert.Nil(t, err)
}

func TestWrongProviderId(t *testing.T) {
	node := v1.Node{
		Spec: v1.NodeSpec{
			ProviderID: "wrong-prodiver://123",
		},
	}

	_, err := GetDoDropletId(&node)
	assert.NotNil(t, err)
}
