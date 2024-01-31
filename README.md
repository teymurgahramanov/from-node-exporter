# From-node exporter

The From-Node Exporter for Prometheus is designed to probe the accessibility of external endpoints from each node of the Kubernetes cluster over TCP, HTTP, and ~~ICMP~~.

## How will it be useful?

While there are other tools like the Blackbox Exporter, the From-node Exporter focuses specifically on simplicity and efficiency for Kubernetes node-level probing. It's designed to serve a specific use case - ensuring all required endpoints are accessible from every node of your cluster.

Assume such case, during node failure, your pods as normally have been evicted to another node but after that your services got broken because the node has not network access to essential external services such as API of government service, SMS gateway and etc.

Or ...
You have during disaster you have switched to DR site, and surprised that your endpoints of your company's partrners are not accessible because of latest firewall maintenance works.

Here are some cases where the From-node exporter will be useful:

> In my org we have several k8s clusters and quite unreliable security department, who has control over firewall and have a habbit of corrupting the rules on said firewall. The confusion is immense. The issue is that at any point in time *one or several nodes can lose access to one or several external resources*.
https://www.mail-archive.com/prometheus-users@googlegroups.com/msg06409.html

> In our DevOps environment, where rapid development and deployment occur, the access controls and permissions are in constant flux. The security team struggles to keep up with the dynamic nature of resource allocation and permissions. Consequently, developers and operations staff experience frequent instances where specific nodes or services lose connectivity to external dependencies due to misconfigured access controls or inconsistent policies.

> Within our company, we manage multiple cloud environments, and the network configuration is frequently altered by a dynamic group of administrators. The challenge we face is that our team relies heavily on specific API connections, and due to the unpredictable nature of network changes, we experience intermittent disruptions. These disruptions result in critical services losing connectivity, impacting our overall operational efficiency.

## Current state

The From-node Exporter is intentionally kept simple. Currently, no plans are in place to add additional functionality or metrics, except ICMP probe, as other tools like the Blackbox Exporter are already comprehensive in their feature set.

## Metrics

By default metrics and their descriptions are available on ```:8080/metrics```.

## Install

### Helm chart

#### 1. Clone

```
git clone https://github.com/teymurgahramanov/from-node-exporter
```

#### 2. Configure targets

Configure targets in the Helm values file. Refer to [example.config.yaml](./example.config.yaml).

#### 3. Install Helm chart

```
# I'm sure you know how
```

#### 3. Configure Prometheus

Configuration snippet will be provided in Helm output upon the chart installation. Refer to [NOTES.txt](chart/templates/NOTES.txt).

### Binary

1. Download binary from [the releases tab](https://github.com/teymurgahramanov/from-node-exporter/releases).
2. Configure targets in __config.yaml__. Refer to [example.config.yaml](./example.config.yaml).
3. Run ```./from-node-exporter```

## Contributing
Contributions to enhance or fix issues are welcome. Feel free to submit pull requests.
