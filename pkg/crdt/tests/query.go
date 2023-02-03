package tests

import messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"

// queryModifiers are handy for building more complex queries.

type queryModifier func(*messagev1.QueryRequest)

func timeRange(start, end uint64) queryModifier {
	return func(q *messagev1.QueryRequest) {
		q.StartTimeNs = start
		q.EndTimeNs = end
	}
}

func withPagingInfo(q *messagev1.QueryRequest, f func(pi *messagev1.PagingInfo)) {
	if q.PagingInfo == nil {
		q.PagingInfo = new(messagev1.PagingInfo)
	}
	f(q.PagingInfo)
}

func limit(l uint32) queryModifier {
	return func(q *messagev1.QueryRequest) {
		withPagingInfo(q, func(pi *messagev1.PagingInfo) {
			pi.Limit = l
		})
	}
}

func descending() queryModifier {
	return func(q *messagev1.QueryRequest) {
		withPagingInfo(q, func(pi *messagev1.PagingInfo) {
			pi.Direction = messagev1.SortDirection_SORT_DIRECTION_DESCENDING
		})
	}
}

func cursor(cursor *messagev1.Cursor) queryModifier {
	return func(q *messagev1.QueryRequest) {
		withPagingInfo(q, func(pi *messagev1.PagingInfo) {
			pi.Cursor = cursor
		})
	}
}
