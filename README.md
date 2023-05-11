# xmtpd

[![Test](https://github.com/xmtp/xmtpd/actions/workflows/test.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/test.yml)
[![Build](https://github.com/xmtp/xmtpd/actions/workflows/build.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/build.yml)
[![Publish](https://github.com/xmtp/xmtpd/actions/workflows/publish-node.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/publish-node.yml)

**⚠️ Experimental:** This software is in early development. Expect frequent changes and unresolved issues.

`xmtpd` (XMTP daemon) is an experimental version of XMTP node software. It is **not** the node software that currently forms the XMTP network.

After `xmtpd` meets specific functional requirements, the plan is for it to become the node software that powers the XMTP network. 

Some of these requirements include reaching functional parity with the current node software and reliably performing data replication without data loss. 

To keep up with and provide feedback about `xmtpd` development, see the [Issues tab](https://github.com/xmtp/xmtpd/issues) in this repo.

At this time, you can experiment with `xmtpd` as follows:

- Run a local `xmtpd` node.
- Run a cluster of `xmtpd` nodes in a local devnet.

In the future, anyone will be able to run an `xmtpd` node that participates in the XMTP network.

## Prerequisites

- [Go](https://go.dev/doc/install)
- [Docker](https://www.docker.com/get-started/)

## Install and build dependencies

To install dependencies and build background services needed by the node, run:

```sh
dev/up
```

This starts the following background services in Docker:

- A [Postgres](https://www.postgresql.org/) database that is used to store data needed to run the node

- A [Prometheus](https://prometheus.io/) instance that enables you to view metrics and data about the node, such as resource usage and errors

- A [Grafana](https://grafana.com/) instance that enables you to monitor the health and performance of the node

To stop the background services running in Docker, run:

```sh
dev/down
```

## Start the node

To start the `xmtpd` node, run:

```sh
dev/start
```

## Test the node

To run tests against the `xmtpd` node, run:

```sh
dev/test
```

These tests provide a full suite of unit and integration tests for the `xmtpd` repo to help ensure and maintain correctness of the code over time and to avoid regressions as the code evolves. You can explore the tests by taking a look at any files with the suffix `_test.go`. 

## Connect the node to the experimental XMTP testnet

You can connect the `xmtpd` node to an experimental XMTP testnet.

You can use https://pizza.xmtp.chat/ to try messaging on the XMTP testnet.

**Warning:** The XMTP testnet is unstable and we regularly delete data from it.

To connect an `xmtpd` node to the XMTP testnet:

1. Do this...

2. Then this...

3. Then this...

4. Restart the node: `dev/down`, `dev/up`, and `dev/start`???

5. To confirm that your node is connected, xxxxxx.

## Monitor the node

The `xmtpd` node build provides two options for monitoring your node.

- To access your local Prometheus instance to explore node metrics, run:

    ```sh
    open http://localhost:9090
    ```

    To learn how to query node data in Prometheus, see [Metric Types in Prometheus and PromQL](https://promlabs.com/blog/2020/09/25/metric-types-in-prometheus-and-promql) and [The Anatomy of a PromQL Query](https://promlabs.com/blog/2020/06/18/the-anatomy-of-a-promql-query/).

- To access your local Grafana instance to explore and build node dashboards, run:

    ```sh
    open http://localhost:3000
    ```

    To learn how to visualize node data in Grafana, see [Prometheus Histograms with Grafana Heatmaps](https://towardsdatascience.com/prometheus-histograms-with-grafana-heatmaps-d556c28612c7) and [How to visualize Prometheus histograms in Grafana](https://grafana.com/blog/2020/06/23/how-to-visualize-prometheus-histograms-in-grafana/).

## Create a local `xmtpd` devnet

To learn how to create a cluster of `xmtpd` nodes to run a local devnet, see [Create an XMTP devnet](./dev/net/README.md).
