// Package types defines custom types for the db package.
package types

import "time"

// GatewayEnvelopeRow represents a single envelope to be inserted in a batch.
type GatewayEnvelopeRow struct {
	OriginatorNodeID     int32
	OriginatorSequenceID int64
	Topic                []byte
	PayerID              int32
	GatewayTime          time.Time
	Expiry               int64
	OriginatorEnvelope   []byte
	SpendPicodollars     int64
}
