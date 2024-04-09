# KWOK test driver

The KWOK test driver creates an instance of the NuoDB Control Plane running in KWOK, which is a Kubernetes test cluster that simulates scheduling of pods.
The KWOK cluster does not support running actual containers and consists of Etcd, Kubernetes API server, controller-manager, and KWOK itself.

The NuoDB Control Plane Operator and REST API are deployed separately as Docker containers that have access to the Kubernetes API server.
To enable scheduling of statefulset pods, a noop volume provisioner is also deployed.
