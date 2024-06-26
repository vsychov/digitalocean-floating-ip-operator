package k8s

import (
	"github.com/vsychov/digitalocean-floating-ip-operator/pkg/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/klog/v2"
	"log"
	"net"
	"os"
	"strings"
)

// K8s client wrapper
type K8s struct {
	ClientSet *kubernetes.Clientset
	Config    config.Config
}

// NewInstance create new instance of K8s
func NewInstance(applicationConfig config.Config) (instance K8s, err error) {
	clientset, err := newClientset()
	if err != nil {
		return
	}

	return K8s{
		ClientSet: clientset,
		Config:    applicationConfig,
	}, nil
}

func newClientset() (client *kubernetes.Clientset, err error) {
	// creates the in-cluster config
	cfg, err := inClusterCustomConfig()
	if err != nil {
		return
	}

	//todo: for local proxy only
	//config.Host = strings.Replace(config.Host, "https", "http", 1)

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return
	}

	return clientset, nil
}

// copy-paste from pkg/mod/k8s.io/client-go@v0.21.3/rest/config.go
// add way to override tokenFile & rootCAFile path
func inClusterCustomConfig() (*rest.Config, error) {
	tokenFile, ok := os.LookupEnv("KUBERNETES_TOKEN_FILE")
	if !ok {
		tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	}

	rootCAFile, ok := os.LookupEnv("KUBERNETES_ROOT_CA_FILE")
	if !ok {
		rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	}

	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return nil, rest.ErrNotInCluster
	}

	log.Println("Kubernetes host and port is: ", host, port)
	token, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}

	tlsClientConfig := rest.TLSClientConfig{}

	if _, err := certutil.NewPool(rootCAFile); err != nil {
		klog.Errorf("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
	} else {
		tlsClientConfig.CAFile = rootCAFile
	}

	return &rest.Config{
		Host:            "https://" + net.JoinHostPort(host, port),
		TLSClientConfig: tlsClientConfig,
		BearerToken:     string(token),
		BearerTokenFile: tokenFile,
	}, nil
}

func (k8s *K8s) getK8sServerVersion() string {

	serverVersion, err := k8s.ClientSet.ServerVersion()
	if err != nil {
		log.Fatalf("Unable to get k8s server version: %s", err)
	}

	//for values like v1.26.3-gke.1000
	versions := strings.Split(serverVersion.String(), "-")
	version := strings.Replace(versions[0], "v", "", 1)

	return version
}
