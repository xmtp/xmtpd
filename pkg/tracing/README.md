# XMTPD APM Tracing

This package provides Datadog APM distributed tracing for xmtpd, enabling end-to-end visibility into message processing, cross-node replication, and database operations.

## Configuration

Tracing is **disabled by default** and must be explicitly enabled via the application flag.
When enabled, tracing is **deterministic** — 100% of traces are collected, with no sampling.

### Enabling Tracing

Set the `XMTPD_TRACING_ENABLE` environment variable (or CLI flag `--tracing.enable`):

```bash
export XMTPD_TRACING_ENABLE=true
```

In Terraform, set the variable on the xmtpd server module:

```hcl
tracing_enable = true
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `XMTPD_TRACING_ENABLE` | `false` | Set to `true` to enable tracing. **This is the only control.** |
| `ENV` | `test` | Environment name (used as Datadog env tag) |
| `DD_AGENT_HOST` | `localhost` | Datadog agent host (standard DD env var) |
| `DD_TRACE_AGENT_PORT` | `8126` | Datadog agent port (standard DD env var) |

### Zero Overhead When Disabled

When `XMTPD_TRACING_ENABLE` is not set (the default), all span creation functions
return a shared no-op singleton. There is no allocation, no string formatting,
and no network I/O. The DB tracer falls back to Prometheus-only logging.

## Instrumented Paths

### Write Path (Message Publish)
```
node.publish_payer_envelopes
  └── node.stage_transaction
      └── pgx.query [db.role=writer] INSERT INTO staged_originator_envelopes...
          │
          ▼ (async - trace_linked=true/false)
      publish_worker.process
          ├── publish_worker.calculate_fees
          ├── publish_worker.sign_envelope
          ├── publish_worker.insert_gateway
          │   └── pgx.query INSERT INTO gateway_envelopes...
          └── publish_worker.delete_staged
              └── pgx.query DELETE FROM staged_originator_envelopes...
```

### Read Path (Message Query)
```
node.query_envelopes
  └── pgx.query [db.role=reader] SELECT * FROM gateway_envelopes...
```

### Subscribe Path (Client Streaming)
```
subscribe_worker.dispatch
  ├── [batch_size, envelopes_parsed, parse_errors]
  └── subscribe_worker.listener_closed [reason=channel_full|context_done]
```

### Cross-Node Replication
```
sync.connect_to_node [target_node, target_address]
  └── sync.setup_stream
      └── sync.subscribe_envelopes
          └── sync.receive_batch [source_node, num_envelopes]
              └── sync.validate_envelope [sequence_id, topic]
                  └── sync_worker.store_envelope
                      ├── sync_worker.verify_fees
                      └── sync_worker.insert_gateway
```

## Key Tags for Debugging

### Identifying Read-Replica Issues
```
trigger=timer_fallback     # Notification was missed, fell back to polling
notification_miss=true     # DBSubscription poll found nothing after notification
db.role=reader             # Query went to read replica
db.role=writer             # Query went to primary
```

### Tracing Async Processing
```
trace_linked=true          # Async span successfully linked to parent
trace_linked=false         # Parent context not found (TTL expired or timer fallback)
staged_id=12345            # Envelope ID for correlation
```

### Client Issues
```
reason=channel_full        # Client not consuming fast enough
reason=context_done        # Client disconnected
dropped_envelopes=5        # Number of envelopes client missed
```

## Example Datadog Queries

### Find all timer fallbacks (read-replica bug indicator)
```
service:xmtpd @trigger:timer_fallback
```

### Find empty query results
```
service:xmtpd operation_name:xmtpd.node.query_envelopes @zero_results:true
```

### Find slow database queries
```
service:xmtpd resource_name:pgx.query @duration:>100ms
```

### Find queries that hit the read replica
```
service:xmtpd resource_name:pgx.query @db.role:reader
```

### Find dropped client connections
```
service:xmtpd operation_name:xmtpd.subscribe_worker.listener_closed @reason:channel_full
```

### Find cross-node sync errors
```
service:xmtpd operation_name:xmtpd.sync.* @error:true
```

### Find out-of-order envelopes
```
service:xmtpd @out_of_order:true
```

## Production Safety Limits

The tracing package includes built-in limits to prevent runaway resource usage:

| Limit | Value | Purpose |
|-------|-------|---------|
| `MaxTagValueLength` | 1024 | String tags longer than this are truncated |
| `MaxStoreSize` | 10000 | Maximum entries in TraceContextStore |

## Troubleshooting

### Traces not appearing in Datadog

1. Verify `XMTPD_TRACING_ENABLE=true` is set
2. Verify agent is running: `curl http://localhost:8126/info`
3. Verify `ENV` is set correctly for your environment
4. Check `DD_AGENT_HOST` if agent is not on localhost

### Missing parent spans (trace_linked=false)

This indicates async context propagation failed. Causes:
- TraceContextStore TTL expired (5 minutes)
- Timer fallback was used instead of notification
- Envelope processed by different worker instance

## Architecture Notes

### Async Context Propagation

The `TraceContextStore` bridges async boundaries between staging and worker processing:

1. When `PublishPayerEnvelopes` stages an envelope, it stores the span context
2. When `publish_worker` processes the envelope, it retrieves the context
3. Child spans are created with `ChildOf(parentContext)` for trace linking
4. Contexts expire after 5 minutes to prevent memory leaks

### Composite Database Tracer

Database tracing uses a composite pattern to preserve existing functionality:
- `tracelog.TraceLog` - Prometheus metrics logging (existing, always active)
- `apmQueryTracer` - Datadog APM spans (only wired when tracing is enabled)
