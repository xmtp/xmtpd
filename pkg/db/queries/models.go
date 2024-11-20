// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

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

type GatewayEnvelope struct {
	GatewayTime          time.Time
	OriginatorNodeID     int32
	OriginatorSequenceID int64
	Topic                []byte
	OriginatorEnvelope   []byte
}

type LatestBlock struct {
	ContractAddress string
	BlockNumber     int64
}

type NodeInfo struct {
	NodeID      int32
	PublicKey   []byte
	SingletonID int16
}

type StagedOriginatorEnvelope struct {
	ID             int64
	OriginatorTime time.Time
	Topic          []byte
	PayerEnvelope  []byte
}
