# From-node exporter

The From-Node Exporter for Prometheus is designed to probe the accessibility of external endpoints from each node of the Kubernetes cluster over TCP, HTTP, and ~~ICMP~~.

## How will it be useful?

While there are other tools like the Blackbox Exporter, the Node Probe Exporter focuses specifically on simplicity and efficiency for Kubernetes node-level probing. It's designed to serve a specific use case - ensuring all required endpoints are accessible from every node of your cluster.

## Current state

The Node Probe Exporter is intentionally kept simple. Currently, no plans are in place to add additional functionality or metrics, except ICMP probe, as other tools like the Blackbox Exporter are already comprehensive in their feature set.

## Usage
### 1. Clone

```
git clone https://github.com/teymurgahramanov/from-node-exporter
```

### 2. Configure targets

Configure targets using the Helm values file. Refer to [values.yaml](./chart/values.yaml). Example:
```
config:
  targets:
    - target1:
        address: api.example.com:8080
        module: tcp
        timeout: 10
    - target2:
        address: https://example.com
        module: http
        interval: 60
        timeout: 5
```

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

## Run on local

1. Download binary
```
https://github.com/teymurgahramanov/from-node-exporter/releases/download/v0.1.0/from-node-exporter
```
2. Configure targets in __config.yaml__
```
config:
  targets:
    - target1:
        address: api.example.com:8080
        module: tcp
        timeout: 10
    - target2:
        address: https://example.com
        module: http
        interval: 60
        timeout: 5
```
3. Run
```
./from-node-exporter
```

## Contributing
Contributions to enhance or fix issues are welcome. Feel free to submit pull requests.