# xmptd OpenMetrics catalog

This document catalogs the [OpenMetrics](https://prometheus.io/docs/specs/om/open_metrics_spec/) instrumentation for xmptd.

| Name | Type | Description | File |
|------|------|-------------|------|
| `xmtp_api_failed_grpc_requests_counter` | `Counter` | Number of failed GRPC requests by code | `pkg/metrics/api.go` |
| `xmtp_api_incoming_node_connection_by_version_gauge` | `Gauge` | Number of incoming node connections by version | `pkg/metrics/api.go` |
| `xmtp_api_node_connection_requests_by_version_counter` | `Counter` | Number of incoming node connections by version | `pkg/metrics/api.go` |
| `xmtp_api_open_connections_gauge` | `Gauge` | Number of open API connections | `pkg/metrics/api.go` |
| `xmtp_blockchain_publish_payload_seconds` | `Histogram` | Time to publish a payload to the blockchain | `pkg/metrics/blockchain.go` |
| `xmtp_blockchain_wait_for_transaction_seconds` | `Histogram` | Time spent waiting for transaction receipt | `pkg/metrics/blockchain.go` |
| `xmtp_indexer_log_processing_time_seconds` | `Histogram` | Time to process a blockchain log | `pkg/metrics/indexer.go` |
| `xmtp_indexer_log_streamer_block_lag` | `Gauge` | Lag between current block and max block | `pkg/metrics/indexer.go` |
| `xmtp_indexer_log_streamer_current_block` | `Gauge` | Current block being processed by the log streamer | `pkg/metrics/indexer.go` |
| `xmtp_indexer_log_streamer_get_logs_duration` | `Histogram` | Duration of the get logs call | `pkg/metrics/indexer.go` |
| `xmtp_indexer_log_streamer_get_logs_requests` | `Counter` | Number of get logs requests | `pkg/metrics/indexer.go` |
| `xmtp_indexer_log_streamer_logs` | `Counter` | Number of logs found by the log streamer | `pkg/metrics/indexer.go` |
| `xmtp_indexer_log_streamer_max_block` | `Gauge` | Max block on the chain to be processed by the log streamer | `pkg/metrics/indexer.go` |
| `xmtp_indexer_retryable_storage_error_count` | `Counter` | Number of retryable storage errors | `pkg/metrics/indexer.go` |
| `xmtp_migrator_destination_blockchain_last_sequence_id` | `Gauge` | Last sequence ID published to blockchain | `pkg/metrics/migrator.go` |
| `xmtp_migrator_destination_database_last_sequence_id` | `Gauge` | Last sequence ID persisted in destination database | `pkg/metrics/migrator.go` |
| `xmtp_migrator_e2e_latency_seconds` | `Histogram` | Time spent migrating a message | `pkg/metrics/migrator.go` |
| `xmtp_migrator_reader_errors_total` | `Counter` | Total number of reader errors | `pkg/metrics/migrator.go` |
| `xmtp_migrator_reader_fetch_duration_seconds` | `Histogram` | Time spent fetching records from source database | `pkg/metrics/migrator.go` |
| `xmtp_migrator_reader_num_rows_found` | `Counter` | Number of rows fetched from source database | `pkg/metrics/migrator.go` |
| `xmtp_migrator_source_last_sequence_id` | `Gauge` | Last sequence ID pulled from source DB | `pkg/metrics/migrator.go` |
| `xmtp_migrator_transformer_errors_total` | `Counter` | Total number of transformation errors | `pkg/metrics/migrator.go` |
| `xmtp_migrator_writer_errors_total` | `Counter` | Total number of writer errors by destination and error type | `pkg/metrics/migrator.go` |
| `xmtp_migrator_writer_latency_seconds` | `Histogram` | Time spent writing to destination | `pkg/metrics/migrator.go` |
| `xmtp_migrator_writer_retry_attempts` | `Histogram` | Number of retry attempts before success or failure | `pkg/metrics/migrator.go` |
| `xmtp_migrator_writer_rows_migrated` | `Counter` | Total number of rows successfully migrated | `pkg/metrics/migrator.go` |
| `xmtp_payer_failed_attempts_to_publish_to_node_via_banlist` | `Histogram` | Number of failed attempts to publish to a node via banlist | `pkg/metrics/payer.go` |
| `xmtp_payer_get_nodes_available_nodes` | `Gauge` | Number of currently available nodes for reader selection | `pkg/metrics/payer.go` |
| `xmtp_payer_lru_nonce` | `Gauge` | Least recently used blockchain nonce of the payer (not guaranteed to be the highest nonce). | `pkg/metrics/payer.go` |
| `xmtp_payer_messages_originated` | `Counter` | Number of messages originated by the payer. | `pkg/metrics/payer.go` |
| `xmtp_payer_node_publish_duration_seconds` | `Histogram` | Duration of the node publish call | `pkg/metrics/payer.go` |
| `xmtp_payer_read_own_commit_in_time_seconds` | `Histogram` | Read your own commit duration in seconds | `pkg/metrics/payer.go` |
| `xmtp_sync_failed_outgoing_sync_connections` | `Gauge` | Gauge of current failed outgoing sync connections | `pkg/metrics/sync.go` |
| `xmtp_sync_failed_outgoing_sync_connections_counter` | `Counter` | Counter of total number of failed outgoing sync connection attempts | `pkg/metrics/sync.go` |
| `xmtp_sync_messages_received_count` | `Counter` | Count of messages received from the originator | `pkg/metrics/sync.go` |
| `xmtp_sync_messages_received_error_count` | `Counter` | Count of failed/errored messages received from the originator | `pkg/metrics/sync.go` |
| `xmtp_sync_originator_sequence_id` | `Gauge` | Last synced sequence id of the originator | `pkg/metrics/sync.go` |
| `xmtp_sync_outgoing_sync_connections` | `Gauge` | Gauge of open outgoing sync connections | `pkg/metrics/sync.go` |
