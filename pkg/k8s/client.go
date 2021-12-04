package k8s

import (
	"float-ip-do-k8s/pkg/config"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/klog/v2"
	"net"
	"os"
)

type K8s struct {
	ClientSet *kubernetes.Clientset
	Config    config.Config
}

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
	tokenFile, ok := os.LookupEnv("KUBERNTERS_TOKEN_FILE")
	if !ok {
		tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	}

	rootCAFile, ok := os.LookupEnv("KUBERNTERS_ROOT_CA_FILE")
	if !ok {
		rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	}

	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return nil, rest.ErrNotInCluster
	}

	token, err := ioutil.ReadFile(tokenFile)
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
