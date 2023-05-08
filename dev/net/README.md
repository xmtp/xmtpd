# Create a local XMTP devnet

**⚠️ Experimental:** This software is in early development. Expect frequent changes and unresolved issues.

Use this component to configure and provision a cluster of [experimental `xmtpd` nodes](/README.md) running on [Kubernetes](https://kubernetes.io/) in a devnet. This component creates the cluster based on [Terraform](https://terraform.io/) configurations.

At this time, you can experiment with `xmtpd` nodes running in a devnet.

## Create a local devnet

The [default Terraform configuration](/dev/terraform/plans/devnet-local/main.tf) uses [kind](https://kind.sigs.k8s.io/) to create a local Kubernetes cluster using Docker container nodes.

To provision a devnet, run:

```sh
dev/net/up
```

To bring down the devnet, run:

```sh
dev/net/down
```

### Interact with the `nodes` API

You can interact with the `nodes` API via `localhost` (port 80). 

You can also interact with each node individually via `${NODE_NAME}.localhost`. For example:

```sh
curl -s -XPOST node1.localhost/message/v1/query -d '{"content_topics":["topic"]}' | jq
```

### Monitor the devnet

To access your local [Prometheus](https://prometheus.io/) instance to explore metrics, run:

```sh
open http://prometheus.localhost
```

To access your local [Grafana](https://grafana.com/) instance to explore and build dashboards, run:

```sh
dev/net/copy-grafana-password
open http://grafana.localhost
```

### Interact with the Kubernetes cluster

To interact with the Kubernetes cluster directly using `kubectl`, export the `KUBECONFIG`. To do this, run:

```sh
source dev/net/k8s-env
```

This `k8s-env` script also creates the following command-line aliases that you can use to interact with specific namespaces in the cluster:

```sh
alias kn="kubectl -n xmtp-nodes"
alias ks="kubectl -n xmtp-system"
alias kt="kubectl -n xmtp-tools"
```
