# From-node exporter

The From-Node Exporter for Prometheus is designed to probe external endpoints from each Kubernetes node.

## Why From-node Exporter?

While there are other tools like the Blackbox Exporter, the Node Probe Exporter focuses specifically on simplicity and efficiency for Kubernetes node-level probing. It's designed to serve a specific use case - ensuring all required endpoints are accessible from every node.

## Limitations

- __Scope__: The Node Probe Exporter focuses exclusively on TCP and HTTP probes and is intentionally kept simple. Currently no plans are in place to add additional functionality or metrics, as other tools like the Blackbox Exporter are already comprehensive in their feature set.

## Install



## Usage
