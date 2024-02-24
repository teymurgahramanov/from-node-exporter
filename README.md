# From-node exporter

The From-node Exporter for Prometheus is designed to probe the accessibility of external endpoints from each node of the Kubernetes cluster over TCP, HTTP, and ICMP.

## How will it be useful?

While there are other tools like the Blackbox Exporter, the From-node Exporter focuses specifically on simplicity and efficiency for Kubernetes node-level probing. It's designed to serve a specific use case - ensuring all required endpoints are accessible from every node of your cluster.

For instance, consider a situation where your pods were evicted to a different server because of a node failure. The node has been added to the cluster recently. However, it resulted in errors and service unavailability because of the lack of access to essential external endpoints. It turned out that the security department had recently performed firewall maintenance and an error was made in configuring access rules that affected new cluster nodes. 

Or another case ...
> In my org we have several k8s clusters and quite unreliable security department, who has control over firewall and have a habbit of corrupting the rules on said firewall. The confusion is immense. The issue is that at any point in time *one or several nodes can lose access to one or several external resources*.
https://www.mail-archive.com/prometheus-users@googlegroups.com/msg06409.html

## Current state

The From-node Exporter is intentionally kept simple. Currently, no plans are in place to add additional functionality or metrics as other tools like the Blackbox Exporter are already comprehensive in their feature set.

## Run
### Kubernetes

1. Clone repo
2. Create config.yaml
3. Install Helm chart

### Docker
1. Create config.yaml
2. Run

    ```
    docker run --rm -v $(pwd)/config.yaml:/config.yaml teymurgahramanov/from-node-exporter:v0.2.0
    ```

### Binary

1. Download binary from [the releases tab](https://github.com/teymurgahramanov/from-node-exporter/releases)
2. Create config.yaml
3. Run ```./from-node-exporter```

## config.yaml

## Metrics

By default metrics and their descriptions are available on ```:8080/metrics```.

Prometheus job example:
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

## Note
-  __The ICMP probe requires elevated privileges to function__ \
Refer to https://github.com/prometheus-community/pro-bing?tab=readme-ov-file#supported-operating-systems



## Contributing

Contributions to enhance or fix issues are welcome. Feel free to submit pull requests.
