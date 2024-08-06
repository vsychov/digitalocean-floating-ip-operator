package k8s

import (
	"context"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

const fieldManager = "do-floating-ip-operator"

// AddRouteForFloatIp execute job to create routing rules on node
func (k8s *K8s) AddRouteForFloatIp(node *v1.Node, ip string) {
	//TODO: add watch and alert for failed job (float-ip-*)
	jobName := fmt.Sprintf("float-ip-%s", node.Name) //job name should unique for node, for prevent RC
	jobNamespace := string(k8s.Config.RoutingJob.Namespace)
	saName := string(k8s.Config.RoutingJob.ServiceAccountName)

	existingJobs, err := k8s.ClientSet.BatchV1().Jobs(jobNamespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s,status.successful!=1", jobName),
	})

	if err != nil {
		panic(err)
	}

	if len(existingJobs.Items) != 0 {
		log.Printf("Still have not completed jobs on node %s, skip this.", node.Name)
		return
	}

	jobBackoffLimit := int32(0)
	jobParallelRuns := int32(1)
	jobTtlRemoveAfterSuccess := int32(60) //sec

	priority := int32(2000001000)
	privileged := true
	capabilities := v1.Capabilities{
		Add: []v1.Capability{
			"NET_ADMIN",
		},
	}

	jobSecurityContext := v1.SecurityContext{
		Privileged:   &privileged,
		Capabilities: &capabilities,
	}

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: jobNamespace,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &jobTtlRemoveAfterSuccess,
			Parallelism:             &jobParallelRuns,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					PriorityClassName:  "system-node-critical",
					Priority:           &priority,
					NodeName:           node.Name,
					ServiceAccountName: saName,
					Containers: []v1.Container{
						{
							ImagePullPolicy: v1.PullAlways,
							SecurityContext: &jobSecurityContext,
							Name:            jobName,
							Image:           string(k8s.Config.RoutingJob.ImageName) + ":" + k8s.getK8sServerVersion(),
							//ip route replace default via 10.19.0.1 dev eth0
							Command: []string{
								"/bin/sh",
								"-c",
								"GW=`curl -sS http://169.254.169.254/metadata/v1/interfaces/public/0/anchor_ipv4/gateway` && " +
									"echo ip route replace default via $GW dev eth0 && " +
									"sysctl -w net.ipv4.ip_forward=1 && sysctl -w net.ipv6.conf.all.forwarding=1 && " +
									"ip route replace default via $GW dev eth0 && " +
									//"sleep 10000 && " +
									"IP_S1=`curl -sS 2ip.io`; " +
									"IP_S2=`curl -sS ifconfig.io`; " +
									"if [ \"$IP_S1\" = \"" + ip + ", \" ] || [ \"$IP_S2\" = \"" + ip + "\" ]; then " +
									"echo IP Changed success to " + ip + " ; " +
									"kubectl label nodes " + node.Name + " egress-ready=true --overwrite ; " +
									"else " +
									"echo Failed to change IP ; " +
									"fi;",
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
					HostNetwork:   true,
				},
			},
			BackoffLimit: &jobBackoffLimit,
		},
	}

	_, err = k8s.ClientSet.BatchV1().Jobs(jobNamespace).Create(context.TODO(), jobSpec, metav1.CreateOptions{
		FieldManager: fieldManager,
	})

	if err != nil {
		panic(err)
	}

	log.Printf("Launch route-assign job for verify")
}
