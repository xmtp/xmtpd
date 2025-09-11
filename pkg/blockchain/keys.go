package blockchain

// All your well-known keys — unchanged, just moved here for discoverability.

const (
	APP_CHAIN_GATEWAY_PAUSED_KEY                     = "xmtp.appChainGateway.paused"
	DISTRIBUTION_MANAGER_PAUSED_KEY                  = "xmtp.distributionManager.paused"
	DISTRIBUTION_MANAGER_PROTOCOL_FEES_RECIPIENT_KEY = "xmtp.distributionManager.protocolFeesRecipient"

	GROUP_MESSAGE_BROADCASTER_MAX_PAYLOAD_SIZE_KEY = "xmtp.groupMessageBroadcaster.maxPayloadSize"
	GROUP_MESSAGE_BROADCASTER_MIN_PAYLOAD_SIZE_KEY = "xmtp.groupMessageBroadcaster.minPayloadSize"
	GROUP_MESSAGE_BROADCASTER_PAUSED_KEY           = "xmtp.groupMessageBroadcaster.paused"
	GROUP_MESSAGE_PAYLOAD_BOOTSTRAPPER_KEY         = "xmtp.groupMessageBroadcaster.payloadBootstrapper"

	IDENTITY_UPDATE_BROADCASTER_MAX_PAYLOAD_SIZE_KEY = "xmtp.identityUpdateBroadcaster.maxPayloadSize"
	IDENTITY_UPDATE_BROADCASTER_MIN_PAYLOAD_SIZE_KEY = "xmtp.identityUpdateBroadcaster.minPayloadSize"
	IDENTITY_UPDATE_BROADCASTER_PAUSED_KEY           = "xmtp.identityUpdateBroadcaster.paused"
	IDENTITY_UPDATE_PAYLOAD_BOOTSTRAPPER_KEY         = "xmtp.identityUpdateBroadcaster.payloadBootstrapper"

	NODE_REGISTRY_ADMIN_KEY               = "xmtp.nodeRegistry.admin"
	NODE_REGISTRY_MAX_CANONICAL_NODES_KEY = "xmtp.nodeRegistry.maxCanonicalNodes"

	PAYER_REGISTRY_MINIMUM_DEPOSIT_KEY      = "xmtp.payerRegistry.minimumDeposit"
	PAYER_REGISTRY_PAUSED_KEY               = "xmtp.payerRegistry.paused"
	PAYER_REGISTRY_WITHDRAW_LOCK_PERIOD_KEY = "xmtp.payerRegistry.withdrawLockPeriod"

	PAYER_REPORT_MANAGER_PROTOCOL_FEE_RATE_KEY = "xmtp.payerReportManager.protocolFeeRate"

	RATE_REGISTRY_CONGESTION_FEE_KEY         = "xmtp.rateRegistry.congestionFee"
	RATE_REGISTRY_MESSAGE_FEE_KEY            = "xmtp.rateRegistry.messageFee"
	RATE_REGISTRY_MIGRATOR_KEY               = "xmtp.rateRegistry.migrator"
	RATE_REGISTRY_STORAGE_FEE_KEY            = "xmtp.rateRegistry.storageFee"
	RATE_REGISTRY_TARGET_RATE_PER_MINUTE_KEY = "xmtp.rateRegistry.targetRatePerMinute"

	SETTLEMENT_CHAIN_GATEWAY_PAUSED_KEY = "xmtp.settlementChainGateway.paused"
)
