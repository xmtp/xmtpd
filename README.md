# xmtpd

[![Test](https://github.com/xmtp/xmtpd/actions/workflows/test.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/test.yml)
[![Build](https://github.com/xmtp/xmtpd/actions/workflows/build.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/build.yml)
[![Publish](https://github.com/xmtp/xmtpd/actions/workflows/publish-node.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/publish-node.yml)

**⚠️ Experimental:** This software is in early development. Expect frequent changes and unresolved issues.

`xmtpd` (XMTP daemon) is an experimental version of XMTP node software. It is **not** the node software that currently forms the XMTP network.

After `xmtpd` meets specific functional requirements, the plan is for it to become the node software that powers the XMTP network.

Some of these requirements include reaching functional parity with the current node software and reliably performing data replication without data loss.

To keep up with and provide feedback about `xmtpd` development, see the [Issues tab](https://github.com/xmtp/xmtpd/issues) in this repo.

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
dev/run
```

## Test the node

To run tests against the `xmtpd` node, run:

```sh
dev/test
```

These tests provide a full suite of unit and integration tests for the `xmtpd` repo to help ensure and maintain correctness of the code over time and to avoid regressions as the code evolves. You can explore the tests by taking a look at any files with the suffix `_test.go`.

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

# Contributing

Please follow the [style guide](https://google.github.io/styleguide/go/decisions).

## Modifying the protobuf schema

Submit and land a PR to https://github.com/xmtp/proto. Then run:

```sh
dev/generate
```

## Modifying the database schema

Create a new migration by running:

```sh
dev/gen-migration
```

Fill in the migrations in the generated files. If you are unfamiliar with migrations, you may follow [this guide](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md). The database is PostgreSQL and the driver is PGX.

## Modifying database queries

We use [sqlc](https://docs.sqlc.dev/en/latest/index.html) to generate the code for our DB queries. Modify the `queries.sql` file, and then run:

```sh
sqlc generate
```
