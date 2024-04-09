# K8s test driver

The K8s test driver creates an instance of the NuoDB Control Plane into the configured Kubernetes cluster, i.e. the one specified by KUBECONFIG or ~/.kube/config.
This installs the Helm charts for the NuoDB Control Plane, along with the following ancillary services:
- Cert-manager to enable webhooks.
- Nginx Ingress to enable external connectivity.
