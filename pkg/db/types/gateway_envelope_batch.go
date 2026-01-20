// Package types defines custom types for the db package.
package types

import (
	"slices"
	"time"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

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

type GatewayEnvelopeBatch struct {
	Envelopes []GatewayEnvelopeRow
}

func NewGatewayEnvelopeBatch() *GatewayEnvelopeBatch {
	return &GatewayEnvelopeBatch{
		Envelopes: make([]GatewayEnvelopeRow, 0),
	}
}

func (b *GatewayEnvelopeBatch) Add(envelope GatewayEnvelopeRow) {
	b.Envelopes = append(b.Envelopes, envelope)
}

func (b *GatewayEnvelopeBatch) All() []GatewayEnvelopeRow {
	b.ensureOrdered()
	return slices.Clone(b.Envelopes)
}

func (b *GatewayEnvelopeBatch) LastSequenceID() int64 {
	b.ensureOrdered()
	if len(b.Envelopes) == 0 {
		return 0
	}
	return b.Envelopes[len(b.Envelopes)-1].OriginatorSequenceID
}

func (b *GatewayEnvelopeBatch) Len() int {
	return len(b.Envelopes)
}

func (b *GatewayEnvelopeBatch) Reset() {
	b.Envelopes = make([]GatewayEnvelopeRow, 0)
}

func (b *GatewayEnvelopeBatch) ToParams() queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams {
	n := b.Len()

	b.ensureOrdered()

	params := queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams{
		OriginatorNodeIds:     make([]int32, n),
		OriginatorSequenceIds: make([]int64, n),
		Topics:                make([][]byte, n),
		PayerIds:              make([]int32, n),
		GatewayTimes:          make([]time.Time, n),
		Expiries:              make([]int64, n),
		OriginatorEnvelopes:   make([][]byte, n),
		SpendPicodollars:      make([]int64, n),
	}

	for i, row := range b.Envelopes {
		params.OriginatorNodeIds[i] = row.OriginatorNodeID
		params.OriginatorSequenceIds[i] = row.OriginatorSequenceID
		params.Topics[i] = row.Topic
		params.PayerIds[i] = row.PayerID
		params.GatewayTimes[i] = row.GatewayTime
		params.Expiries[i] = row.Expiry
		params.OriginatorEnvelopes[i] = row.OriginatorEnvelope
		params.SpendPicodollars[i] = row.SpendPicodollars
	}

	return params
}

func (b *GatewayEnvelopeBatch) ensureOrdered() {
	slices.SortFunc(b.Envelopes, func(a, b GatewayEnvelopeRow) int {
		if a.OriginatorSequenceID < b.OriginatorSequenceID {
			return -1
		}
		if a.OriginatorSequenceID > b.OriginatorSequenceID {
			return 1
		}
		return 0
	})
}
