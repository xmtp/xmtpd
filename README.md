# xmtpd

[![Github tag](https://img.shields.io/github/v/tag/xmtp/xmtpd.svg?sort=semver)](https://github.com/xmtp/xmtpd/releases/latest)
[![Test](https://github.com/xmtp/xmtpd/actions/workflows/test.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/test.yml)
[![Lint](https://github.com/xmtp/xmtpd/actions/workflows/lint-go.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/lint-go.yml)
[![Build](https://github.com/xmtp/xmtpd/actions/workflows/build-xmtpd.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/build-xmtpd.yml)
[![License](https://img.shields.io/github/license/xmtp/xmtpd)](https://github.com/xmtp/xmtpd/blob/main/LICENSE)
[![Twitter](https://img.shields.io/twitter/follow/xmtp_)](https://x.com/xmtp_)

**⚠️ Experimental:** This software is in early development. Expect frequent changes and unresolved issues.

`xmtpd` (XMTP daemon) is an experimental version of XMTP node software. It is **not** the node software that currently forms the XMTP network.

After `xmtpd` meets specific functional requirements, the plan is for it to become the node software that powers the XMTP network.

Some of these requirements include reaching functional parity with the current node software and reliably performing data replication without data loss.

To keep up with and provide feedback about `xmtpd` development, see the [Issues tab](https://github.com/xmtp/xmtpd/issues) in this repo.

> [!TIP]
> Looking for info on how to get started as a registered node operator? See the [xmtpd-infrastructure](https://github.com/xmtp/xmtpd-infrastructure) repository, which provides infrastructure-as-code examples and tooling to help node operators deploy and manage xmtpd nodes.

## Prerequisites

- [Go 1.25](https://go.dev/doc/install) -- On macOS install with `brew install go@1.25`
- [Docker](https://www.docker.com/get-started/)
- [Foundry](https://github.com/foundry-rs/foundry)

## Quickstart

### Install and build dependencies

To install dependencies and build background services needed by the node, run:

```sh
dev/up
```

This starts the following background services in Docker:

- A [Postgres](https://www.postgresql.org/) database that is used to store data needed to run the node

- A [Prometheus](https://prometheus.io/) instance that enables you to view metrics and data about the node, such as resource usage and errors

To stop the background services running in Docker, run:

```sh
dev/down
```

### Start the node

To start the `xmtpd` node, run:

```sh
dev/run
```

### Test the node

To run tests against the `xmtpd` node, run:

```sh
dev/test
```

These tests provide a full suite of unit and integration tests for the `xmtpd` repo to help ensure and maintain correctness of the code over time and to avoid regressions as the code evolves. You can explore the tests by taking a look at any files with the suffix `_test.go`.

### Monitor the node

To enable metrics you need to start the node with `--metrics.enable`.

By default, we expose Prometheus metrics on port `8008`.

```sh
curl http://localhost:8008
```

## About the XMTP testnet

The XMTP testnet environment lives in [XMTP Sepolia Chain](https://xmtp-testnet.explorer.alchemy.com/).

The XMTP testnet nodes operated by Ephemera include:

| DNS name                           | Location   | Public key                                                           |
| ---------------------------------- | ---------- | -------------------------------------------------------------------- |
| https://grpc.testnet.xmtp.network  | US-EAST-2  | 0x03e5442c5d1fe2f02b6b9a1a386383a7766860b40a6079a0223994ffa2ce10512c |
| https://grpc2.testnet.xmtp.network | EU-NORTH-1 | 0x02fc261d43a0153539a4c64c29763cb0e7e377c0eac2910c3d4bedb2235ac70371 |

## Deploy a local environment for developing with xmtpd

See [Deploy a local environment for developing with xmtpd](doc/deploy.md).

## Support a new XMTP testnet node operator

See [Support a new XMTP testnet node operator](doc/onboarding.md).

## Migrate to XMTP node registry V2

See [Node Registry Migration](doc/node-registry-migration-v2.md).

## xmtpd OpenMetrics catalog

See [xmptd OpenMetrics catalog](doc/metrics_catalog.md).

## Contributing

Please follow the [style guide](https://google.github.io/styleguide/go/decisions).

### Modifying the protobuf schema

Submit and land a PR to https://github.com/xmtp/proto. Then run:

```sh
dev/gen/all
```

Or wait for the [nightly automation](https://github.com/xmtp/xmtpd/actions/workflows/nightly.yml) to finish.

### Modify the database schema

Create a new migration by running:

```sh
dev/gen/migration {migration-name}
```

Fill in the migrations in the generated files. If you are unfamiliar with migrations, you may follow [this guide](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md). The database is PostgreSQL and the driver is PGX.

### Modify database queries

We use [sqlc](https://docs.sqlc.dev/en/latest/index.html) to generate the code for our DB queries. Modify the `queries.sql` file, and then run:

```sh
dev/gen/all
```

If needed, you may use `dev/psql` to test more complex SQL syntax. To determine the level of support in sqlc, there is a sqlc [playground](https://play.sqlc.dev/p/f6eebe941750560934cefa943c77f63497debc828c487e8d1771fb6d83773246) for experimenting with how the query syntax translates into Go code.
