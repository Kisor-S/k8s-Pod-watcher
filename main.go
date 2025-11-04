package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	// Define command-line flags
	kubeconfig := flag.String("kubeconfig", homeDir()+"/.kube/config", "(optional) absolute path to the kubeconfig file")
	namespace := flag.String("namespace", "default", "namespace to watch (empty means all namespaces)")
	flag.Parse()

	// Build kuubernetes config (local kubeconfig or in-cluster config)
	config, err := buildConfig(*kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build kube config: %v\n", err)
		os.Exit(1)
	}

	// Build Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build kube config: %v\n", err)
		os.Exit(1)
	}

	// Create a stop channel and handle SIGINT/SIGTERM for graceful shutdown
	stopCh := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("\nReceived shutdown signal, stopping...")
		close(stopCh)
	}()

	// SharedInformerFactory: namespace-scoped if namespace set, otherwise cluster-wide
	var factory informers.SharedInformerFactory
	if *namespace == "" {
		factory = informers.NewSharedInformerFactory(clientset, 0)
	} else {
		factory = informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace(*namespace))
	}

	// Get Pod informer
	podInformer := factory.Core().V1().Pods().Informer()

}

// homeDir returns the home directory for the current user
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// buildConfig loads kubeconfig or falls back to in-cluster config

func buildConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath != "" {
		return clientcmd.BuildConfigFromFlags("", filepath.Clean(kubeconfigPath))
	}
	// Fallback to in-cluster config
	return rest.InClusterConfig()
}
