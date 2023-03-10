package api

import proto "github.com/xmtp/proto/v3/go/message_api/v1"

func NewQuery(topic string, modifiers ...QueryModifier) *proto.QueryRequest {
	q := &proto.QueryRequest{}
	if topic != "" {
		q.ContentTopics = []string{topic}
	}
	for _, m := range modifiers {
		m(q)
	}
	return q
}

// QueryModifiers are handy for building more complex queries.
type QueryModifier func(*proto.QueryRequest)

func TimeRange(start, end uint64) QueryModifier {
	return func(q *proto.QueryRequest) {
		q.StartTimeNs = start
		q.EndTimeNs = end
	}
}

func withPagingInfo(q *proto.QueryRequest, f func(pi *proto.PagingInfo)) {
	if q.PagingInfo == nil {
		q.PagingInfo = new(proto.PagingInfo)
	}
	f(q.PagingInfo)
}

func Limit(l uint32) QueryModifier {
	return func(q *proto.QueryRequest) {
		withPagingInfo(q, func(pi *proto.PagingInfo) {
			pi.Limit = l
		})
	}
}

func Descending() QueryModifier {
	return func(q *proto.QueryRequest) {
		withPagingInfo(q, func(pi *proto.PagingInfo) {
			pi.Direction = proto.SortDirection_SORT_DIRECTION_DESCENDING
		})
	}
}

func Ascending() QueryModifier {
	return func(q *proto.QueryRequest) {
		withPagingInfo(q, func(pi *proto.PagingInfo) {
			pi.Direction = proto.SortDirection_SORT_DIRECTION_ASCENDING
		})
	}
}

// Set cursor from previous response if present
func Cursor(resp *proto.QueryResponse) QueryModifier {
	return func(q *proto.QueryRequest) {
		if resp == nil || resp.PagingInfo == nil || resp.PagingInfo.Cursor == nil {
			return
		}
		withPagingInfo(q, func(pi *proto.PagingInfo) {
			pi.Cursor = resp.PagingInfo.Cursor
		})
	}
}
