# xmtpd

[![Test](https://github.com/xmtp/xmtpd/actions/workflows/test.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/test.yml)
[![Build](https://github.com/xmtp/xmtpd/actions/workflows/build-xmtpd.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/build-xmtpd.yml)

**⚠️ Experimental:** This software is in early development. Expect frequent changes and unresolved issues.

`xmtpd` (XMTP daemon) is an experimental version of XMTP node software. It is **not** the node software that currently forms the XMTP network.

After `xmtpd` meets specific functional requirements, the plan is for it to become the node software that powers the XMTP network.

Some of these requirements include reaching functional parity with the current node software and reliably performing data replication without data loss.

To keep up with and provide feedback about `xmtpd` development, see the [Issues tab](https://github.com/xmtp/xmtpd/issues) in this repo.

## Prerequisites

- [Go 1.24](https://go.dev/doc/install) -- On macOS install with `brew install go@1.24`
- [Docker](https://www.docker.com/get-started/)
- [Foundry](https://github.com/foundry-rs/foundry)

Note that a dependency of this repo will not compile with version of go newer than 1.22.

## Install and build dependencies

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

To enable metrics you need to start the node with `--metrics.enable`.
By default, we expose Prometheus metrics on port `8008`.

  ```sh
  curl http://localhost:8008
  ```

# Contributing

Please follow the [style guide](https://google.github.io/styleguide/go/decisions).

## Modifying the protobuf schema

Submit and land a PR to https://github.com/xmtp/proto. Then run:

```sh
dev/generate
```

Or wait for the [nightly automation](https://github.com/xmtp/xmtpd/actions/workflows/nightly.yml) to finish.

## Modifying the database schema

Create a new migration by running:

```sh
dev/gen-migration {migration-name}
```

Fill in the migrations in the generated files. If you are unfamiliar with migrations, you may follow [this guide](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md). The database is PostgreSQL and the driver is PGX.

## Modifying database queries

We use [sqlc](https://docs.sqlc.dev/en/latest/index.html) to generate the code for our DB queries. Modify the `queries.sql` file, and then run:

```sh
dev/generate
```

If needed, you may use `dev/psql` to test more complex SQL syntax. To determine the level of support in sqlc, there is a sqlc [playground](https://play.sqlc.dev/p/f6eebe941750560934cefa943c77f63497debc828c487e8d1771fb6d83773246) for experimenting with how the query syntax translates into Go code.
