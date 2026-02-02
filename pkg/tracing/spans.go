// Package tracing provides span operation name constants for APM instrumentation.
// Using constants ensures consistency and makes refactoring easier.
package tracing

// Span operation names follow the pattern: xmtpd.{component}.{operation}
// This provides clear hierarchy in Datadog APM service maps.
const (
	// Node API spans - incoming request handling
	SpanNodePublishPayerEnvelopes = "xmtpd.node.publish_payer_envelopes"
	SpanNodeQueryEnvelopes        = "xmtpd.node.query_envelopes"
	SpanNodeStageTransaction      = "xmtpd.node.stage_transaction"
	SpanNodeWaitGatewayPublish    = "xmtpd.node.wait_gateway_publish"

	// Publish worker spans - async envelope processing
	SpanPublishWorkerProcess       = "xmtpd.publish_worker.process"
	SpanPublishWorkerCalculateFees = "xmtpd.publish_worker.calculate_fees"
	SpanPublishWorkerSignEnvelope  = "xmtpd.publish_worker.sign_envelope"
	SpanPublishWorkerInsertGateway = "xmtpd.publish_worker.insert_gateway"
	SpanPublishWorkerDeleteStaged  = "xmtpd.publish_worker.delete_staged"

	// Subscribe worker spans - client streaming
	SpanSubscribeWorkerDispatch       = "xmtpd.subscribe_worker.dispatch"
	SpanSubscribeWorkerBroadcast      = "xmtpd.subscribe_worker.broadcast"
	SpanSubscribeWorkerListenerClosed = "xmtpd.subscribe_worker.listener_closed"

	// DB subscription spans - polling mechanism
	SpanDBSubscriptionPoll = "xmtpd.db_subscription.poll"

	// Sync worker spans - cross-node replication (receiving side)
	SpanSyncConnectToNode   = "xmtpd.sync.connect_to_node"
	SpanSyncSetupStream     = "xmtpd.sync.setup_stream"
	SpanSyncSubscribe       = "xmtpd.sync.subscribe_envelopes"
	SpanSyncReceiveBatch    = "xmtpd.sync.receive_batch"
	SpanSyncValidateEnvelope = "xmtpd.sync.validate_envelope"

	// Envelope sink spans - storing synced envelopes
	SpanSyncWorkerStoreEnvelope         = "xmtpd.sync_worker.store_envelope"
	SpanSyncWorkerVerifyFees            = "xmtpd.sync_worker.verify_fees"
	SpanSyncWorkerInsertGateway         = "xmtpd.sync_worker.insert_gateway"
	SpanSyncWorkerStoreReservedEnvelope = "xmtpd.sync_worker.store_reserved_envelope"
	SpanSyncWorkerStorePayerReport      = "xmtpd.sync_worker.store_payer_report"
	SpanSyncWorkerStoreAttestation      = "xmtpd.sync_worker.store_attestation"

	// Database spans
	SpanDBQuery = "xmtpd.db.query"
)

// Span tag keys - use these for consistency
const (
	TagTrigger          = "trigger"
	TagStagedID         = "staged_id"
	TagOriginatorNode   = "originator_node"
	TagSourceNode       = "source_node"
	TagTargetNode       = "target_node"
	TagTopic            = "topic"
	TagSequenceID       = "sequence_id"
	TagNumEnvelopes     = "num_envelopes"
	TagNumResults       = "num_results"
	TagZeroResults      = "zero_results"
	TagNotificationMiss = "notification_miss"
	TagTraceLinked      = "trace_linked"
	TagOutOfOrder       = "out_of_order"
	TagDBRole           = "db.role"
	TagDBStatement      = "db.statement"
	TagDBRowsAffected   = "db.rows_affected"
	TagDBSystem         = "db.system"
	TagDBService        = "db.service"
)

// Trigger values for the trigger tag
const (
	TriggerNotification  = "notification"
	TriggerTimerFallback = "timer_fallback"
)

// DB role values
const (
	DBRoleReader = "reader"
	DBRoleWriter = "writer"
)
