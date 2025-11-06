# XMTP Node Communication APIs

This document describes the communication protocols and APIs implemented in the xmtpd server. The server provides both client-facing APIs and internal node-to-node synchronization mechanisms.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Protocol Stack](#protocol-stack)
- [Client-Facing APIs](#client-facing-apis)
  - [Replication API](#replication-api)
  - [Metadata API](#metadata-api)
  - [Health Check API](#health-check-api)
  - [Reflection API](#reflection-api)
- [Internal Communication](#internal-communication)
  - [Sync Service](#sync-service)
- [Configuration](#configuration)
- [Testing and Development](#testing-and-development)

## Architecture Overview

### TLDR

The xmtpd system uses a **dual-protocol architecture**:

- **External/Client-Facing APIs**: Implemented with **Connect-RPC handlers** that automatically support three protocols:

  - Connect (HTTP/1.1 or HTTP/2 with standard HTTP semantics)
  - gRPC (HTTP/2 with gRPC protocol)
  - gRPC-Web (browser-compatible gRPC)

- **Internal Node-to-Node Communication**: Uses **native gRPC clients** for:

  - Sync worker (node-to-node envelope replication)
  - MLS validation service
  - Gateway client manager (gateway-to-node communication)

This approach provides **maximum flexibility for clients** (they can use any protocol) while maintaining **high-performance gRPC** for internal communication where protocol flexibility isn't needed.

### System Diagram

```ascii
┌──────────────────────────────────────────────────────────────────────┐
│                      XMTP Client Applications                        │
│              (Connect, gRPC, or gRPC-Web protocols)                  │
└──────────────────────────────────────────────────────────────────────┘
                                  │
                    HTTP/1.1, HTTP/2, or gRPC
                                  │
        ┌─────────────────────────┴───────────────────────────┐
        │                                                     │
        ▼                                                     ▼
┌──────────────────────────┐                   ┌──────────────────────────┐
│   Gateway (Payer API)    │                   │   XMTPD Node (Port: 5050)│
│      Port: 5051          │                   │                          │
│ ┌──────────────────────┐ │                   │  ┌────────────────────┐  │
│ │ • PublishEnvelopes   │ │  gRPC (internal)  │  │ • Replication API  │  │
│ │ • GetNodes           │ │◄─────────────────►│  │ • Metadata API     │  │
│ └──────────────────────┘ │                   │  │ • Health Check     │  │
│  (Connect-RPC handlers)  │                   │  │ • Reflection       │  │
└──────────────────────────┘                   │  └────────────────────┘  │
                                               │   (Connect-RPC handlers) │
                                               └──────────────────────────┘
                                                            │
                                                            │        gRPC (native)
                  ┌─────────────────────────────────────────┼──────────────┐
                  │                                         │              │
                  ▼                                         ▼              ▼
       ┌────────────────────┐                    ┌──────────────┐  ┌─────────────┐
       │  Sync Service      │                    │              │  │             │
       │   (Background)     │                    ┤  Database    │  │ MLS Service │
       │                    │                    │              │  │  (gRPC)     │
       │ • Subscribes to    │                    └──────────────┘  └─────────────┘
       │   other nodes      │
       │ • Replicates       │
       │   envelopes        │
       └────────────────────┘
                │
                │ gRPC (native)
                ▼
       ┌──────────────────┐
       │  Other XMTP Nodes│
       └──────────────────┘
```

### Connect-RPC vs gRPC: What's the Difference?

**Connect-RPC** (used for client-facing APIs):

- Built on standard HTTP semantics (works with HTTP/1.1 and HTTP/2)
- Single handler implementation supports 3 protocols automatically
- Uses standard HTTP headers (easier to debug with curl, browser dev tools)
- Better browser support (Connect protocol and gRPC-Web)
- Works seamlessly with HTTP proxies, load balancers (ALB, NGINX)
- Compatible with existing gRPC clients (accepts gRPC protocol)
- Simpler middleware/interceptor model

**Native gRPC** (used for internal communication):

- Requires HTTP/2
- High-performance binary protocol
- Used where we control both client and server
- Existing ecosystem of tools (grpcurl, Buf CLI)
- More mature for service-to-service communication

### Why This Architecture?

1. **Client Flexibility**: External clients can use whatever protocol works best for them (mobile SDKs can use Connect, existing gRPC clients work unchanged, browsers can use Connect-Web or gRPC-Web)

2. **Simplified Deployment**: Connect-RPC works with standard HTTP infrastructure (no special h2c requirements in production with proper load balancers)

3. **Internal Performance**: Native gRPC for internal communication provides battle-tested performance and reliability

4. **Single Codebase**: Connect-RPC handlers are generated from the same .proto files as gRPC, maintaining type safety and consistency

## Protocol Stack

The xmtpd and gateway servers expose APIs using the connect-go implementation:

- **Connect-RPC**: A modern, protocol-agnostic RPC framework built on Protobuf and on top of Go standard lib HTTP
- **HTTP/2 Cleartext (h2c)**: Enables HTTP/2 over unencrypted connections for development/testing
- **gRPC**: Compatible through Connect-RPC's protocol support
- **Protocol Buffers**: For message serialization

### Why h2c?

The server wraps handlers with h2c (HTTP/2 Cleartext) support to:

- Enable HTTP/2 on plaintext (no TLS) connections in development
- Allow gRPC tooling like `grpcurl` to function without TLS locally
- Support gRPC reflection over HTTP/2 without terminating TLS at the server

When running in production, prefer HTTPS with HTTP/2 (no h2c) unless traffic is terminated by an ingress/proxy that speaks HTTP/2 to clients and plain HTTP/1.1 to backends.

### When to use which protocol

- Connect (default via connect-go)

  - Use for Go services and services that benefit from protocol flexibility
  - Works over HTTP/1.1 and HTTP/2; supports Protobuf and JSON
  - Best developer ergonomics and browser compatibility via Connect-Web

- gRPC (classic)

  - Use when interoperating with existing gRPC-only clients/infra
  - Requires HTTP/2; commonly secured with TLS in production
  - Strong ecosystem of tooling (`grpcurl`, language SDKs)

- gRPC-Web
  - Use from browsers (WASM) without a custom transport
  - Typically requires an Envoy/Ingress or the server to support gRPC-Web
  - Connect handlers support gRPC-Web out of the box

Rule of thumb:

- Local dev without TLS: h2c + Connect or gRPC
- Production: HTTPS + HTTP/2; terminate TLS at ingress (Envoy/NGINX) and forward HTTP/1.1 or h2c to backend
- Browser clients: Connect-Web or gRPC-Web

## Client-Facing APIs

### Replication API

The Replication API is the primary interface for XMTP clients to publish and query messages.

**Service Name**: `xmtp.xmtpv4.message_api.ReplicationApi`

#### Methods

##### 1. SubscribeEnvelopes (Server Streaming)

Subscribe to real-time envelope updates for specific topics or nodes.

**Endpoint**: `/xmtp.xmtpv4.message_api.ReplicationApi/SubscribeEnvelopes`

**Request**:

```protobuf
message SubscribeEnvelopesRequest {
  EnvelopesQuery query = 1;
}

message EnvelopesQuery {
  // Client queries by topics
  repeated bytes topics = 1;

  // Node queries by originator node IDs
  repeated uint32 originator_node_ids = 2;

  // Resume from last seen cursor
  Cursor last_seen = 3;
}
```

**Response Stream**:

```protobuf
message SubscribeEnvelopesResponse {
  repeated OriginatorEnvelope envelopes = 1;
}
```

**Usage**: Long-lived streaming connection that delivers new envelopes as they arrive.

##### 2. QueryEnvelopes (Unary)

Query historical envelopes with pagination support.

**Endpoint**: `/xmtp.xmtpv4.message_api.ReplicationApi/QueryEnvelopes`

**Request**:

```protobuf
message QueryEnvelopesRequest {
  EnvelopesQuery query = 1;
  uint32 limit = 2;  // Maximum number of envelopes to return
}
```

**Response**:

```protobuf
message QueryEnvelopesResponse {
  repeated OriginatorEnvelope envelopes = 1;
}
```

**Usage**: Fetch historical messages for specific topics or originator nodes.

##### 3. PublishPayerEnvelopes (Unary)

Publish new message envelopes to the XMTP network.

**Endpoint**: `/xmtp.xmtpv4.message_api.ReplicationApi/PublishPayerEnvelopes`

**Request**:

```protobuf
message PublishPayerEnvelopesRequest {
  repeated PayerEnvelope payer_envelopes = 1;
}
```

**Response**:

```protobuf
message PublishPayerEnvelopesResponse {
  repeated OriginatorEnvelope originator_envelopes = 1;
}
```

**Usage**: Submit signed envelopes with payment information for network propagation.

##### 4. GetInboxIds (Unary)

Resolve blockchain addresses or installation IDs to XMTP inbox IDs.

**Endpoint**: `/xmtp.xmtpv4.message_api.ReplicationApi/GetInboxIds`

**Request**:

```protobuf
message GetInboxIdsRequest {
  repeated Request requests = 1;

  message Request {
    string identifier = 1;
    IdentifierKind identifier_kind = 2;
  }
}
```

**Response**:

```protobuf
message GetInboxIdsResponse {
  repeated Response responses = 1;

  message Response {
    string identifier = 1;
    optional string inbox_id = 2;
  }
}
```

**Usage**: Map user identities to their XMTP inbox IDs.

##### 5. GetNewestEnvelope (Unary)

Get the newest envelope for each requested topic.

**Endpoint**: `/xmtp.xmtpv4.message_api.ReplicationApi/GetNewestEnvelope`

**Request**:

```protobuf
message GetNewestEnvelopeRequest {
  repeated bytes topics = 1;
}
```

**Response**:

```protobuf
message GetNewestEnvelopeResponse {
  map<string, OriginatorEnvelope> envelopes = 1;
}
```

**Usage**: Check for new messages without subscribing or querying full history.

---

### Metadata API

The Metadata API provides information about node state, sync status, and operational metrics.

**Service Name**: `xmtp.xmtpv4.metadata_api.MetadataApi`

#### Method

##### 1. GetSyncCursor (Unary)

Retrieve the current sync cursor (vector clock) for this node.

**Endpoint**: `/xmtp.xmtpv4.metadata_api.MetadataApi/GetSyncCursor`

**Request**:

```protobuf
message GetSyncCursorRequest {}
```

**Response**:

```protobuf
message GetSyncCursorResponse {
  Cursor cursor = 1;
}
```

**Usage**: Determine the node's current position in the distributed message log.

##### 2. SubscribeSyncCursor (Server Streaming)

Subscribe to real-time updates of the sync cursor.

**Endpoint**: `/xmtp.xmtpv4.metadata_api.MetadataApi/SubscribeSyncCursor`

**Request**:

```protobuf
message GetSyncCursorRequest {}
```

**Response Stream**:

```protobuf
message GetSyncCursorResponse {
  Cursor cursor = 1;
}
```

**Usage**: Monitor sync progress in real-time.

##### 3. GetVersion (Unary)

Get the node's software version.

**Endpoint**: `/xmtp.xmtpv4.metadata_api.MetadataApi/GetVersion`

**Request**:

```protobuf
message GetVersionRequest {}
```

**Response**:

```protobuf
message GetVersionResponse {
  string version = 1;
}
```

**Usage**: Verify node version for compatibility checks.

##### 4. GetPayerInfo (Unary)

Retrieve payer spending information and statistics.

**Endpoint**: `/xmtp.xmtpv4.metadata_api.MetadataApi/GetPayerInfo`

**Request**:

```protobuf
message GetPayerInfoRequest {
  string payer_address = 1;
  int64 start_time = 2;
  int64 end_time = 3;
  PayerInfoGranularity granularity = 4;  // HOUR or DAY
}
```

**Response**:

```protobuf
message GetPayerInfoResponse {
  repeated PayerSpend spends = 1;
}
```

**Usage**: Query fee spending history for specific payer addresses.

---

### Health Check API

gRPC Health Check API, implemented using [Connect-RPC grpchealth](connectrpc.com/grpchealth).

**Service Name**: `grpc.health.v1.Health`

**Endpoint**: `/grpc.health.v1.Health/Check`

**Usage**: Monitor service health status. Returns `SERVING_STATUS_SERVING` when healthy.

**Example**:

```bash
curl -X POST http://localhost:5050/grpc.health.v1.Health/Check \
  -H "Content-Type: application/json" \
  -d '{}'
```

---

### Reflection API

gRPC Server Reflection API for service discovery, implemented using [Connect-RPC grpcreflect](connectrpc.com/grpcreflect).

**Service Names**:

- `grpc.reflection.v1.ServerReflection` (v1)
- `grpc.reflection.v1alpha.ServerReflection` (v1alpha)

**Configuration**: Enabled via `XMTPD_REFLECTION_ENABLE=true`

**Usage**: Allows tools like `grpcurl` to discover available services and methods without proto files.

**Example**:

```bash
# List all services
grpcurl -plaintext localhost:5050 list

# List methods of a specific service
grpcurl -plaintext localhost:5050 list xmtp.xmtpv4.message_api.ReplicationApi

# Describe a method
grpcurl -plaintext localhost:5050 describe xmtp.xmtpv4.message_api.ReplicationApi.QueryEnvelopes
```

---

## Internal Communication

### Sync Service

The Sync Service runs in the background and is responsible for node-to-node message replication.

#### Purpose

- Maintains eventual consistency across the XMTP node network
- Subscribes to envelope streams from all registered nodes
- Handles network partitions and node failures with exponential backoff
- Tracks vector clocks to prevent message duplication

#### Architecture

```ascii
┌────────────────────────────────────────────────┐
│           Sync Worker (Background)             │
│                                                │
│  ┌──────────────────────────────────────────┐  │
│  │  Node Registry Watcher                   │  │
│  │  - Monitors new node registrations       │  │
│  │  - Auto-subscribes to new nodes          │  │
│  └──────────────────────────────────────────┘  │
│                                                │
│  ┌──────────────────────────────────────────┐  │
│  │  Per-Node Subscription Workers           │  │
│  │  - One goroutine per remote node         │  │
│  │  - Exponential backoff on failures       │  │
│  │  - Vector clock based sync               │  │
│  └──────────────────────────────────────────┘  │
│                                                │
│  ┌──────────────────────────────────────────┐  │
│  │  Envelope Write Queue                    │  │
│  │  - Batches incoming envelopes            │  │
│  │  - Writes to local database              │  │
│  │  - Validates and deduplicates            │  │
│  └──────────────────────────────────────────┘  │
└────────────────────────────────────────────────┘
```

#### How It Works

1. **Discovery**: The sync worker queries the node registry for all registered nodes
2. **Subscription**: For each node, it creates a streaming `SubscribeEnvelopes` connection
3. **Cursor Management**: Uses vector clocks to track the last seen message from each node
4. **Replication**: Receives envelopes from remote nodes and writes them to the local database
5. **Fault Tolerance**: Implements exponential backoff for failed connections
6. **Dynamic Updates**: Watches for new node registrations and automatically subscribes

#### Enable sync service

The sync service is enabled via `XMTPD_SYNC_ENABLE=true` in the server configuration.

#### Use Cases

- **Multi-region deployments**: Replicate messages across geographically distributed nodes
- **High availability**: Ensure message durability through redundant storage
- **Network resilience**: Automatically recover from network partitions
- **Load distribution**: Clients can query any node and receive consistent results

---

## Configuration

### Environment Variables

#### API Server

```bash
# Enable the API server
export XMTPD_API_ENABLE=true

# API server port (default: 5050)
export XMTPD_API_PORT=5050

# Send keep-alive interval
export XMTPD_API_SEND_KEEP_ALIVE_INTERVAL=30s
```

#### Reflection

```bash
# Enable gRPC reflection (for development)
export XMTPD_REFLECTION_ENABLE=true
```

#### Sync

```bash
# Enable sync service for node-to-node replication
export XMTPD_SYNC_ENABLE=true
```

#### Database

```bash
export XMTPD_DB_WRITER_CONNECTION_STRING="postgres://user:pass@localhost:5432/xmtp?sslmode=disable"
export XMTPD_DB_READER_CONNECTION_STRING="postgres://user:pass@localhost:5432/xmtp?sslmode=disable"
```

#### MLS Validation

```bash
# MLS validation service address (required for API and Indexer)
export XMTPD_MLS_VALIDATION_GRPC_ADDRESS="http://localhost:60051"
```

#### Node Identity

```bash
# Private key for node signing (required for API and Sync)
export XMTPD_SIGNER_PRIVATE_KEY="0x..."
```

### Server Startup

The server initializes services in the following order:

1. **Metrics Server** (if enabled)
2. **Registrant** (node identity management)
3. **MLS Validation Service** (if API or Indexer enabled)
4. **Indexer** (if enabled)
5. **Migrator** (if enabled)
6. **API Server** (if enabled)
7. **Sync Service** (if enabled)
8. **Payer Report Workers** (if enabled)

### Graceful Shutdown

The server handles graceful shutdown via OS signals (SIGINT, SIGTERM, SIGHUP, SIGQUIT):

1. API server stops accepting new connections (30s timeout)
2. Active streams are allowed to complete
3. Sync subscriptions are closed
4. Database connections are released
5. Metrics server is shut down

---

## Testing and Development

### Using curl

**Health Check**:

```bash
curl -X POST http://localhost:5050/grpc.health.v1.Health/Check \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Query Envelopes**:

```bash
curl -X POST http://localhost:5050/xmtp.xmtpv4.message_api.ReplicationApi/QueryEnvelopes \
  -H "Content-Type: application/json" \
  -d '{
    "query": {
      "topics": ["<base64-encoded-topic>"]
    },
    "limit": 10
  }'
```

**Get Version**:

```bash
curl -X POST http://localhost:5050/xmtp.xmtpv4.metadata_api.MetadataApi/GetVersion \
  -H "Content-Type: application/json" \
  -d '{}'
```

### Using grpcurl

**Prerequisites**: Reflection must be enabled (`XMTPD_REFLECTION_ENABLE=true`)

**List all services**:

```bash
grpcurl -plaintext localhost:5050 list
```

**List methods of a service**:

```bash
grpcurl -plaintext localhost:5050 list xmtp.xmtpv4.message_api.ReplicationApi
```

**Call a method**:

```bash
# Get newest envelope
grpcurl -plaintext -d '{"topics": []}' \
  localhost:5050 \
  xmtp.xmtpv4.message_api.ReplicationApi/GetNewestEnvelope

# Get version
grpcurl -plaintext -d '{}' \
  localhost:5050 \
  xmtp.xmtpv4.metadata_api.MetadataApi/GetVersion
```

**Subscribe to a stream**:

```bash
grpcurl -plaintext -d '{"query": {"topics": ["<topic>"]}}' \
  localhost:5050 \
  xmtp.xmtpv4.message_api.ReplicationApi/SubscribeEnvelopes
```

### Using buf CLI

The `buf` CLI provides native Connect-RPC support:

```bash
# Using Connect protocol

buf curl --protocol connect \
  --http2-prior-knowledge \
  http://localhost:5050/xmtp.xmtpv4.metadata_api.MetadataApi/GetVersion

# With reflection

buf curl --protocol connect \
  --http2-prior-knowledge \
  http://localhost:5050 \
  --list-methods

# Using grpc-health-probe

grpc-health-probe -addr localhost:5050 -plaintext
```

Note: requires HTTP/2; with plaintext backends use h2c or a proxy that upgrades to HTTP/2.

---

## Go Client Examples

Below are minimal examples demonstrating how to call the APIs using different protocols.

Replace `BASE_URL` with your server URL, e.g. `http://localhost:5050`.

### Connect-Go client (recommended)

```go
package main

import (
    "context"
    "fmt"
    "net/http"

    connect "connectrpc.com/connect"
    messageapiconnect "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
    messageapi "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
)

func main() {
    baseURL := "http://localhost:5050"
    client := messageapiconnect.NewReplicationApiClient(http.DefaultClient, baseURL)

    // Unary example: GetNewestEnvelope
    resp, err := client.GetNewestEnvelope(
        context.Background(),
        connect.NewRequest(&messageapi.GetNewestEnvelopeRequest{Topics: [][]byte{}}),
    )
    if err != nil {
        panic(err)
    }
    fmt.Println("OK", resp.Msg)
}
```

### gRPC client over TLS (classic gRPC)

For example: Rust clients using tonic, Go applications with classic grpc bindings.

```go
package main

import (
    "context"
    "crypto/tls"
    "log"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    pb "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
)

func main() {
    tlsCfg := &tls.Config{InsecureSkipVerify: true} // do not use in prod
    conn, err := grpc.Dial(
        "localhost:5050",
        grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)),
    )
    if err != nil { log.Fatal(err) }
    defer conn.Close()

    // Note: classic gRPC stubs are not generated in this repo; prefer Connect.
    _ = pb.File_xmtpv4_message_api_message_api_proto // placeholder to show package
    _ = conn
}
```

### gRPC client over h2c (plaintext HTTP/2)

```go
package main

import (
    "context"
    "log"
    "net"
    "net/http"

    "golang.org/x/net/http2"
    "golang.org/x/net/http2/h2c"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

// Dial h2c with grpc
func h2cDialer(ctx context.Context, addr string) (net.Conn, error) {
    d := &net.Dialer{}
    return d.DialContext(ctx, "tcp", addr)
}

func main() {
    // For grpc-go h2c, use insecure creds and custom dialer with http2
    h2Transport := &http.Transport{}
    http2.ConfigureTransport(h2Transport)
    _ = h2c.NewHandler(nil, &http2.Server{}) // ensures import is used

    conn, err := grpc.Dial(
        "localhost:5050",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithContextDialer(h2cDialer),
    )
    if err != nil { log.Fatal(err) }
    defer conn.Close()
}
```

### Connect-Go forcing gRPC protocol

```go
package main

import (
    "context"
    "fmt"
    "net/http"

    connect "connectrpc.com/connect"
    messageapiconnect "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
    messageapi "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
)

func main() {
    baseURL := "http://localhost:5050"
    client := messageapiconnect.NewReplicationApiClient(
        http.DefaultClient,
        baseURL,
        connect.WithGRPC(), // force gRPC protocol
    )

    resp, err := client.QueryEnvelopes(
        context.Background(),
        connect.NewRequest(&messageapi.QueryEnvelopesRequest{}),
    )
    if err != nil { panic(err) }
    fmt.Println("OK", resp.Msg)
}
```

### Connect-Web or gRPC-Web (browser)

Use Connect-Web from the browser and ensure your ingress supports gRPC-Web or Connect over HTTP/1.1.

---

## Health and Reflection Checks

### curl (works for all Connect/gRPC/gRPC-Web handlers)

```bash
# Health check
curl -s -X POST http://localhost:5050/grpc.health.v1.Health/Check \
  -H 'Content-Type: application/json' -d '{}'

# Arbitrary unary method (GetVersion)
curl -s -X POST http://localhost:5050/xmtp.xmtpv4.metadata_api.MetadataApi/GetVersion \
  -H 'Content-Type: application/json' -d '{}'

# Arbitrary unary method (GetNewestEnvelope)
curl -s -X POST http://localhost:5050/xmtp.xmtpv4.message_api.ReplicationApi/GetNewestEnvelope \
  -H 'Content-Type: application/json' -d '{"topics": []}'
```

### grpcurl (requires reflection and HTTP/2)

```bash
# List services (requires XMTPD_REFLECTION_ENABLE=true)
grpcurl -plaintext localhost:5050 list

# List methods of a service
grpcurl -plaintext localhost:5050 list xmtp.xmtpv4.message_api.ReplicationApi

# Describe a method
grpcurl -plaintext localhost:5050 \
  describe xmtp.xmtpv4.metadata_api.MetadataApi.GetVersion

# Call a method
grpcurl -plaintext -d '{}' \
  localhost:5050 \
  xmtp.xmtpv4.metadata_api.MetadataApi/GetVersion
```

If you see HTTP/2 issues with plaintext, ensure the server is running with h2c enabled (it is in development), or use a proxy/ingress that terminates TLS and forwards HTTP/2/h2c appropriately.

```bash
# Development Setup
#
# See `dev/local.env` for a complete development environment configuration.

# Source the environment
source dev/local.env

# Enable reflection for testing
export XMTPD_REFLECTION_ENABLE=true

# Start the services
./dev/up

# In another terminal, test the API
grpcurl -plaintext localhost:5050 list
```

---

## Protocol Details

### HTTP/2 Cleartext (h2c)

The API server wraps all handlers with h2c support to enable HTTP/2 over plain TCP connections. This is implemented in `pkg/api/server.go`:

```go
// Wrap the handler with h2c to support HTTP/2 Cleartext for gRPC reflection.
// This is required for gRPC reflection to work with HTTP/2, and tools such as grpcurl.
h2cHandler := h2c.NewHandler(mux, &http2.Server{})

svc.httpServer = &http.Server{
    Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
    Handler: h2cHandler,
}
```

### Connect-RPC Protocol

Connect-RPC provides three protocol modes:

- **Connect**: Optimized for web and mobile, uses standard HTTP headers
- **gRPC**: Compatible with existing gRPC tooling
- **gRPC-Web**: Browser-compatible gRPC

All three protocols are automatically supported by the server handlers.

### Authentication

The API server supports JWT-based authentication through the `RegistryVerifier`:

- JWTs are signed by the requesting node's private key
- The server verifies signatures using public keys from the node registry
- Authentication is enforced via interceptors on a per-method basis

---

## Troubleshooting

### API Server Not Responding

1. Check if the server is listening:

   ```bash
   lsof -i :5050
   ```

2. Verify reflection is enabled (for grpcurl):

   ```bash
   echo $XMTPD_REFLECTION_ENABLE
   ```

3. Test with curl (always works):

   ```bash
   curl -X POST http://localhost:5050/grpc.health.v1.Health/Check \
     -H "Content-Type: application/json" -d '{}'
   ```

### grpcurl Timeout

If grpcurl times out, it usually means:

- Server is not running with h2c support
- Reflection is disabled
- Port mismatch

Solution: Use curl with JSON or ensure reflection is enabled.

### Connection Refused

- Verify the server is running: `ps aux | grep xmtpd` or `ps aux | grep main`
- Check the port: `netstat -an | grep 5050`
- Review server logs for startup errors

### Sync Issues

If nodes aren't replicating:

1. Check sync is enabled: `echo $XMTPD_SYNC_ENABLE`
2. Verify node registry connectivity
3. Review sync worker logs for connection errors
4. Check vector clock state in the database

---

## Additional Resources

- [XMTP Protocol Documentation](https://github.com/xmtp/proto)
- [Connect-RPC Documentation](https://connectrpc.com/)
- [gRPC Reflection Protocol](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md)
- [Deployment Guide](./deploy.md)
- [Node Onboarding](./onboarding.md)
