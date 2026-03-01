# CLAUDE.md — xmtpd

## Project Overview

Decentralized messaging node written in Go. Replicates messages across a network of nodes using blockchain-based coordination (Ethereum smart contracts for node registry, settlement).

**Module:** `github.com/xmtp/xmtpd` | **Go:** 1.25

## Local Development

```sh
dev/up       # install deps (brew, go mod tidy), start docker (postgres, prometheus), generate .env
dev/down     # stop docker
dev/run      # start replication server with local env vars
dev/test     # run tests (excludes abis, config, proto, mock, testing packages)
dev/lint-fix # run golangci-lint with --fix to auto-fix and format code
dev/psql     # connect to local postgres (localhost:8765, user=postgres, pass=xmtp)
```

## Code Generation

```sh
dev/gen/all              # run all generators (protos, mocks, sqlc, go generate)
dev/gen/protos           # pull from xmtp/proto repo, generate via buf
dev/gen/mocks            # generate mocks via mockery
dev/gen/sqlc             # generate Go from SQL queries
dev/gen/abi              # download contract ABIs, generate Go bindings via abigen
dev/gen/migration {name} # create new up/down migration files via golang-migrate
```

## Tools (tools/go.mod)

- **buf** — protobuf schema management
- **abigen** — Ethereum ABI → Go bindings
- **migrate** — database migrations (golang-migrate)
- **sqlc** — SQL queries → type-safe Go code
- **mockery** — interface mock generation

## Go Conventions

**Style:** [Google Go Style Guide](https://google.github.io/styleguide/go/decisions)

**File naming:** `snake_case.go`

**Logging** (zap):

- Hierarchical logger names: `"xmtpd.api.publish-worker"` (defined in `pkg/utils/log.go`)
- snake_case field names: `zap.String("app_chain_id", ...)`, `zap.Int64("block_number", ...)`
- Lowercase log messages
- Guard hot-path debug logs: `if logger.Core().Enabled(zap.DebugLevel)`

**Errors:**

- Lowercase messages: `fmt.Errorf("failed to parse config: %w", err)`
- Wrap with `%w` for context; return bare `err` when no context needed

**Testing:**

- testify: `require` for fatal assertions, `assert` for non-fatal
- Test utilities in `pkg/testutils/`
- Integration tests use testcontainers for postgres

**Config:**

- Environment variables with `XMTPD_` prefix (e.g. `XMTPD_DB_WRITER_CONNECTION_STRING`)
- CLI flags via `github.com/jessevdk/go-flags` with struct tags
- Options pattern for function configuration

**Linters** (`.golangci.yaml`):

- Linters: errcheck, govet, ineffassign, staticcheck, unused, nilerr, nakedret
- Formatters: gofmt, golines, gofumpt
- Before committing, run `dev/lint-fix` to auto-fix lint issues and format code

## Key Directories

```text
cmd/
  replication/   # main node binary
  gateway/       # gateway service
  xmtpd-cli/    # CLI tool
  prune/         # pruning tool
pkg/
  api/           # gRPC API handlers
  db/            # database layer (sqlc queries, handler, transactions)
  migrations/    # SQL migration files
  blockchain/    # Ethereum chain interaction
  sync/          # node-to-node replication
  indexer/       # blockchain event indexer
  envelopes/     # envelope processing
  config/        # configuration structs
  testutils/     # shared test helpers
    mocks/       # generated mocks
  registry/      # node registry client
  proto/         # generated protobuf code
  abi/           # generated contract bindings
contracts/       # Foundry smart contract project
dev/             # development scripts
tools/           # tool dependencies (go.mod)
```
