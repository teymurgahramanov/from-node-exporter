# From-node exporter

Prometheus exporter to probe external endpoints from each Kubernetes node. It makes you sure that all required endpoints are accessible from every node of cluster, it is very useful for newly provisioned clusters, Currently it allows TCP and HTTP probes and returns single metric: successful or not. I am not planning to add any extra functionality as Blackbox Exporter already doing all the things perfectly.

