# XMTP devnet

This component configures and provisions clusters of XMTP nodes on [Kubernetes](https://kubernetes.io/) via [Terraform](https://terraform.io/).

The default creates a local cluster running on [kind](https://kind.sigs.k8s.io/), with the option to manage remote clusters on clouds like AWS and GCP defined in `dev/terraform/plans` and using the `PLAN` environment variable.

## Local ([kind](https://kind.sigs.k8s.io/))

Provision a local devnet with:

```sh
dev/net/up
```

Tear it down with:

```sh
dev/net/down
```

### Nodes API

Interact with the nodes API via `localhost` (port 80), or each node individually via `${NODE_NAME}.localhost`, for example:

```sh
curl -s -XPOST node1.localhost/message/v1/query -d '{"content_topics":["topic"]}' | jq
```

### Monitoring

Visit the [Grafana](https://prometheus.io/) UI to explore and build dashboards:

```sh
dev/net/copy-grafana-password
open http://grafana.localhost
```

Visit the [Prometheus](https://prometheus.io/) UI to explore metrics:

```sh
open http://prometheus.localhost
```

### Kubernetes

Interact with the Kubernetes cluster directly via `kubectl` by exporting the `KUBECONFIG` with:

```sh
source dev/net/k8s-env
```

This `k8s-env` script also creates a few command-line aliases for interacting with specific namespaces in the cluster:

```sh
alias kn="kubectl -n xmtp-nodes"
alias ks="kubectl -n xmtp-system"
alias kt="kubectl -n xmtp-tools"
```
