package deserializer

type ProtocolVersion uint16

type ContentType uint8

const (
	ContentTypeApplication ContentType = 1
	ContentTypeProposal    ContentType = 2
	ContentTypeCommit      ContentType = 3
)

type GroupID []byte

type SenderType uint8

const (
	SenderTypeReserved          SenderType = 0
	SenderTypeMember            SenderType = 1
	SenderTypeExternal          SenderType = 2
	SenderTypeNewMemberProposal SenderType = 3
	SenderTypeNewMemberCommit   SenderType = 4
)
