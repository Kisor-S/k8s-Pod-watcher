package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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

	// Add event handlers for Add / Update / Delete
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if p, ok := obj.(*corev1.Pod); ok {
				fmt.Printf("[ADDED]   %s/%s  (phase=%s)\n", p.Namespace, p.Name, p.Status.Phase)
			}
		},

		UpdateFunc: func(oldObj, newObj interface{}) {
			if p, ok := newObj.(*corev1.Pod); ok {
				fmt.Printf("[UPDATED] %s/%s  (phase=%s)\n", p.Namespace, p.Name, p.Status.Phase)
			}
		},

		DeleteFunc: func(obj interface{}) {
			// Deletion may come as a tombstone; handle both cases
			switch t := obj.(type) {
			case *corev1.Pod:
				fmt.Printf("[DELETED] %s/%s\n", t.Namespace, t.Name)
			case cache.DeletedFinalStateUnknown:
				if p, ok := t.Obj.(*corev1.Pod); ok {
					fmt.Printf("[DELETED] %s/%s\n", p.Namespace, p.Name)
				}
			default:
				fmt.Printf("[DELETED] unknown object type\n")
			}
		},
	})

	// Start informer factory (runs in background goroutines)
	factory.Start(stopCh)

	// Wait for all caches to sync
	if ok := cache.WaitForCacheSync(stopCh, podInformer.HasSynced); !ok {
		fmt.Fprintf(os.Stderr, "failed to wait for caches to sync\n")
		os.Exit(1)
	}

	fmt.Println("Pod watcher started. Listening for events... (Ctrl+C to stop)")

	// Block until stop signal is received
	<-stopCh

	// give a small grace period for goroutines to finish printing
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Exited.")
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
