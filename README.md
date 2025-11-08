# 🧠 Kubernetes Pod Watcher (Go + client-go)

A simple yet practical Go application that interacts with the Kubernetes API using the official **client-go** library.  
It lists all Pods in a specified namespace and serves as a foundation for learning how to communicate programmatically with a Kubernetes cluster.

---

## 🚀 Features

- Connects to Kubernetes clusters both **locally** and **in-cluster**
- Automatically detects kubeconfig source:
  - `KUBECONFIG` environment variable (if set)
  - Command-line `--kubeconfig` flag
  - In-cluster configuration (when deployed as a Pod)
  - Default `~/.kube/config`
- Lists all Pods in a namespace with their **name**, **phase**, and **node**
- Cleanly structured Go code using **client-go**
- Useful as a base for building custom Kubernetes controllers or operators

---

## 🛠️ Prerequisites

- Go 1.20 or higher  
- Access to a Kubernetes cluster (minikube, kind, EKS, etc.)
- A valid `kubeconfig` file in `~/.kube/config`

---

## 📦 Project Structure

```
├── main.go # Main application
├── go.mod # Go module definition
├── go.sum # Go dependencies lock file
├── Dockerfile # For containerizing the app (optional)
└── deploy/
    ├── rbac.yaml # Minimal RBAC permissions
    └── deployment.yaml # Sample Deployment manifest
```
---

## 🧩 How It Works

The application determines how to connect to Kubernetes in this order:

1. Uses `KUBECONFIG` environment variable if available  
2. Uses `--kubeconfig` flag (if provided at runtime)  
3. Falls back to in-cluster configuration when running inside a Pod  
4. Uses the default `$HOME/.kube/config` path  

Once connected, it lists all Pods in the target namespace and prints details to the console.

---

## 🧪 Run Locally

Make sure you have access to your cluster via `kubectl` first:
```bash
kubectl get pods -n default
```
Then run the app:
```bash
go run main.go --namespace=default
```
You should see output similar to:
```bash
kishore@Kishores-MacBook-Air k8s-Pod-watcher % go run main.go
Using config from: /Users/kishore/.kube/config
[ADDED]   default/demo-nginx  (phase=Running)
Pod watcher started. Listening for events... (Ctrl+C to stop)
[UPDATED] default/demo-nginx  (phase=Running)
```
If you want to specify a different kubeconfig:
```bash
go run main.go --kubeconfig=$HOME/.kube/config --namespace=kube-system
```

---

## 🐳 Run Inside a Cluster

1. Build the Docker image:
```bash
docker build -t pod-watcher:latest .
```
2. Apply RBAC and Deployment:
```bash
kubectl apply -f deploy/rbac.yaml
kubectl apply -f deploy/deployment.yaml
```
3. Check logs:
```bash
kubectl logs -l app=pod-watcher -n default -f
```
You should see similar log output showing Pod events or listings.

---

## 🔐 RBAC Configuration

The ServiceAccount used by this Pod needs permissions to get, list, and watch Pods.

Example minimal RBAC (already included in deploy/rbac.yaml):
```bash
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list", "watch"]
```

---


## 💡 Learning Goals

This project is designed to help you:
- Understand how Go programs interact with Kubernetes APIs
- Learn the use of client-go for resource management
- Build a foundation for more advanced projects such as:
  - Custom Kubernetes Controllers
  - Operators (using Kubebuilder or Operator SDK)
  - Automation tools for DevOps tasks

---

## 🧰 Tech Stack
- Language: Go
- Library: Kubernetes client-go
- Containerization: Docker
- Deployment: Kubernetes (optional)

---

## 📚 Learning Resources

- [Kubernetes client-go](https://github.com/kubernetes/client-go)
- [Kubernetes API Concepts](https://kubernetes.io/docs/concepts/)
- [Writing Controllers in Go](https://book.kubebuilder.io/)

---

## 🤝 Contributing

Feel free to fork, experiment, and raise pull requests!
This project is meant for learning and exploration.

