| Name | Type | Description | File |
|------|------|-------------|------|
| `xmtp_api_failed_grpc_requests_counter` | `Counter` | Number of failed GRPC requests by code | `pkg/metrics/api.go` |
| `xmtp_api_incoming_node_connection_by_version_gauge` | `Gauge` | Number of incoming node connections by version | `pkg/metrics/api.go` |
| `xmtp_api_node_connection_requests_by_version_counter` | `Counter` | Number of incoming node connections by version | `pkg/metrics/api.go` |
| `xmtp_api_open_connections_gauge` | `Gauge` | Number of open API connections | `pkg/metrics/api.go` |
| `xmtp_payer_failed_attempts_to_publish_to_node_via_banlist` | `Histogram` | Number of failed attempts to publish to a node via banlist | `pkg/metrics/payer.go` |
| `xmtp_payer_messages_originated` | `Counter` | Number of messages originated by the payer. | `pkg/metrics/payer.go` |
| `xmtp_payer_node_publish_duration_seconds` | `Histogram` | Duration of the node publish call | `pkg/metrics/payer.go` |
| `xmtp_payer_read_own_commit_in_time_seconds` | `Histogram` | Read your own commit duration in seconds | `pkg/metrics/payer.go` |
| `xmtp_sync_failed_outgoing_sync_connections_counter` | `Counter` | Counter of total number of failed outgoing sync connection attempts | `pkg/metrics/sync.go` |
| `xmtp_sync_messages_received_count` | `Counter` | Count of messages received from the originator | `pkg/metrics/sync.go` |
| `xmtp_sync_messages_received_error_count` | `Counter` | Count of failed/errored messages received from the originator | `pkg/metrics/sync.go` |
| `xmtp_sync_originator_sequence_id` | `Gauge` | Last synced sequence id of the originator | `pkg/metrics/sync.go` |
| `xmtpd_indexer_log_streamer_block_lag` | `Gauge` | Lag between current block and max block | `pkg/metrics/indexer.go` |
| `xmtpd_indexer_log_streamer_current_block` | `Gauge` | Current block being processed by the log streamer | `pkg/metrics/indexer.go` |
| `xmtpd_indexer_log_streamer_get_logs_duration` | `Histogram` | Duration of the get logs call | `pkg/metrics/indexer.go` |
| `xmtpd_indexer_log_streamer_get_logs_requests` | `Counter` | Number of get logs requests | `pkg/metrics/indexer.go` |
| `xmtpd_indexer_log_streamer_logs` | `Counter` | Number of logs found by the log streamer | `pkg/metrics/indexer.go` |
| `xmtpd_indexer_log_streamer_max_block` | `Gauge` | Max block on the chain to be processed by the log streamer | `pkg/metrics/indexer.go` |
| `xmtpd_indexer_retryable_storage_error_count` | `Counter` | Number of retryable storage errors | `pkg/metrics/indexer.go` |
