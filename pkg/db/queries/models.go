// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package queries

import (
	"database/sql"
	"time"
)

type AddressLog struct {
	Address               string
	InboxID               []byte
	AssociationSequenceID sql.NullInt64
	RevocationSequenceID  sql.NullInt64
}

type BlockchainMessage struct {
	BlockNumber          uint64
	BlockHash            []byte
	OriginatorNodeID     int32
	OriginatorSequenceID int64
	IsCanonical          bool
}

type GatewayEnvelope struct {
	GatewayTime          time.Time
	OriginatorNodeID     int32
	OriginatorSequenceID int64
	Topic                []byte
	OriginatorEnvelope   []byte
	PayerID              sql.NullInt32
}

type LatestBlock struct {
	ContractAddress string
	BlockNumber     int64
	BlockHash       []byte
}

type NodeInfo struct {
	NodeID      int32
	PublicKey   []byte
	SingletonID int16
}

type Payer struct {
	ID      int32
	Address string
}

type PayerSequence struct {
	ID        int32
	Available sql.NullBool
}

type StagedOriginatorEnvelope struct {
	ID             int64
	OriginatorTime time.Time
	Topic          []byte
	PayerEnvelope  []byte
}

type UnsettledUsage struct {
	PayerID           int32
	OriginatorID      int32
	MinutesSinceEpoch int32
	SpendPicodollars  int64
}
