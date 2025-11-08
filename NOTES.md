# 🧠 Personal Study Notes — Kubernetes Pod Watcher (Go + client-go)

## 📘 Objective

The goal of this project was to learn how to interact with Kubernetes programmatically using Go and the official client-go library.
Instead of using kubectl commands, I wanted to understand what happens behind the scenes and how Go applications can communicate directly with the Kubernetes API.

---

## ⚙️ What I Built

I created a simple Pod Watcher application that:
- Connects to the Kubernetes cluster using a kubeconfig file or in-cluster credentials.
- Lists all Pods in a given namespace and displays their names, phases, and assigned nodes.
- Can run both locally and inside a Kubernetes cluster.
- This forms the base for understanding how real-world controllers and operators interact with the API.

---

## 🔍 Key Concepts I Learned
### 1. Kubernetes API and client-go

- client-go is the official Go client library used by Kubernetes itself.
- Every call to kubectl eventually talks to the API server — client-go lets you do that directly in Go.
- kubernetes.NewForConfig(config) creates a clientset that lets you access core resources like Pods, Deployments, ConfigMaps, etc.

### 2. Kubeconfig and Authentication
- Kubernetes clients need a *rest.Config object to connect to the cluster.
- There are two ways to build this configuration:
  - clientcmd.BuildConfigFromFlags("", kubeconfigPath) → local kubeconfig (used outside the cluster)
  - rest.InClusterConfig() → used when your app runs inside the cluster as a Pod.
- If your app runs locally, it looks for ~/.kube/config.
- Inside a cluster, it uses the ServiceAccount token mounted at /var/run/secrets/kubernetes.io/serviceaccount.

### 3. RBAC (Role-Based Access Control)

- Even if your app connects to the cluster, it can only do what RBAC allows.
- I learned how to create minimal RBAC permissions for listing and watching Pods:
```bash
apiGroups: [""]
resources: ["pods"]
verbs: ["get", "list", "watch"]
```

### 4. Go Context and API Calls

- All Kubernetes API calls take a context.Context parameter for timeout and cancellation control.
- Example:
```bash
pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
```

### 5. Fallback Configuration Logic

I implemented a buildConfig() function that determines how to connect:
1. Use KUBECONFIG if set
2. Use --kubeconfig flag (if provided)
3. Use in-cluster configuration
4. Fallback to $HOME/.kube/config

This made my app portable — runnable both locally and in production clusters.

---

## 💡 Challenges and Learnings

| Challenge                                       | What I Learned                                                                                                                   |
| ----------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| The app panicked with “nil pointer dereference” | I learned that `rest.InClusterConfig()` only works inside Kubernetes Pods. For local testing, we must provide a kubeconfig file. |
| `KUBECONFIG` environment variable was empty     | I realized that Go doesn’t automatically fall back to `$HOME/.kube/config` unless coded explicitly.                              |
| Understanding RBAC permissions                  | I learned how ServiceAccounts and Roles control what an in-cluster app can do.                                                   |
| Handling multiple kubeconfig sources            | I learned to build more resilient fallback logic that matches how `kubectl` works.                                               |

---

## ✨ Personal Takeaways

- This project helped me move from “using Kubernetes” to “programmatically interacting with Kubernetes.”
- I now understand how Operators and Controllers watch resources and respond to events.
- Writing in Go deepened my understanding of types, interfaces, and error handling.
- This small project gave me the foundation to build more complex Kubernetes automation tools in the future.

---

## 🧰 Tools & Tech I Used

- Go (language)
- client-go (Kubernetes Go client)
- Docker (for building images)
- kubectl (for cluster access)
- minikube (for testing locally)
- YAML manifests for RBAC and Deployment

---

## 🔮 Next Steps

- Convert this Pod Watcher into a Pod Event Notifier that logs or sends alerts when Pods are created/deleted.
- Extend it to a custom controller or operator using Kubebuilder.
- Add unit tests using fake Kubernetes clients (k8s.io/client-go/kubernetes/fake).
