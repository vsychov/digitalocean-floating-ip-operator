package k8s

import (
	"context"
	"encoding/json"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"strconv"
)

type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

// SetEgressReadyLabel set label "egress-ready" = false or true
func (k8s K8s) SetEgressReadyLabel(node *v1.Node, ready bool) {
	payload := []patchStringValue{{
		Op:    "replace",
		Path:  "/metadata/labels/egress-ready",
		Value: strconv.FormatBool(ready),
	}}

	payloadBytes, _ := json.Marshal(payload)

	patchResult, err := k8s.ClientSet.CoreV1().Nodes().Patch(context.TODO(), node.Name, types.JSONPatchType, payloadBytes, metav1.PatchOptions{})

	if err != nil {
		panic(err)
	}

	log.Printf("Node patched, labels: %+v", patchResult.Labels)
}
