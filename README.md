# From-node exporter

The From-Node Exporter for Prometheus is designed to probe the accessibility of external endpoints from each node of the Kubernetes cluster over TCP, HTTP, and ~~ICMP~~.

## How will it be useful?

While there are other tools like the Blackbox Exporter, the From-node Exporter focuses specifically on simplicity and efficiency for Kubernetes node-level probing. It's designed to serve a specific use case - ensuring all required endpoints are accessible from every node of your cluster.

Example case
> In my org we have several k8s clusters and quite unreliable security 
department, who has control over firewall and have a habbit of corrupting the rules on said firewall. 
The confusion is immense. The issue is that at any point in time *one or 
several nodes can lose access to one or several external resources*.
https://www.mail-archive.com/prometheus-users@googlegroups.com/msg06409.html

## Current state

The From-node Exporter is intentionally kept simple. Currently, no plans are in place to add additional functionality or metrics, except ICMP probe, as other tools like the Blackbox Exporter are already comprehensive in their feature set.

## Usage
### 1. Clone

```
git clone https://github.com/teymurgahramanov/from-node-exporter
```

### 2. Configure targets

Configure targets in the Helm values file. Refer to [example.config.yaml](./example.config.yaml).

### 3. Install Helm chart

```
# I'm sure you know how
```

### 3. Configure Prometheus

Configuration snippet will be provided in Helm output upon the chart installation. Example:
```
- job_name: from-node-exporter
  kubernetes_sd_configs:
    - role: endpoints
  relabel_configs:
  - source_labels: [__meta_kubernetes_endpoints_name]
    regex: from-node-exporter
    action: keep
  - source_labels: [__meta_kubernetes_endpoint_node_name]
    action: replace
    target_label: instance
```

## Run binary

1. Download from [the releases tab](https://github.com/teymurgahramanov/from-node-exporter/releases).
2. Configure targets in __config.yaml__. Refer to [example.config.yaml](./example.config.yaml).
3. Run ```./from-node-exporter```

## Contributing
Contributions to enhance or fix issues are welcome. Feel free to submit pull requests.
