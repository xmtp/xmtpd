package storetest

import (
	"testing"

	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/store/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func TestStore_QueryEnvelopes(t *testing.T, storeMaker TestStoreMaker) {
	t.Helper()

	t.Run("all sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 20)
		test.RequireProtoEqual(t, envs, res.Envelopes)
	})

	t.Run("all sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 20)
		test.RequireProtoEqual(t, envs, res.Envelopes)
	})

	t.Run("all sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
			},
		})
		require.NoError(t, err)
		utils.Reverse(envs)
		require.Len(t, res.Envelopes, 20)
		test.RequireProtoEqual(t, envs, res.Envelopes)
	})

	t.Run("limit sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Limit: 5,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 5)
		test.RequireProtoEqual(t, envs[:5], res.Envelopes)
	})

	t.Run("limit sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
				Limit:     5,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 5)
		test.RequireProtoEqual(t, envs[:5], res.Envelopes)
	})

	t.Run("limit sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
				Limit:     5,
			},
		})
		require.NoError(t, err)
		utils.Reverse(envs)
		require.Len(t, res.Envelopes, 5)
		test.RequireProtoEqual(t, envs[:5], res.Envelopes)
	})

	t.Run("start time sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 10,
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 11)
		test.RequireProtoEqual(t, envs[9:], res.Envelopes)
	})

	t.Run("end time sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			EndTimeNs: 10,
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 10)
		test.RequireProtoEqual(t, envs[:10], res.Envelopes)
	})

	t.Run("time range sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 5,
			EndTimeNs:   15,
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 11)
		test.RequireProtoEqual(t, envs[4:15], res.Envelopes)
	})

	t.Run("start time sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 11)
		test.RequireProtoEqual(t, envs[9:], res.Envelopes)
	})

	t.Run("end time sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			EndTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 10)
		test.RequireProtoEqual(t, envs[:10], res.Envelopes)
	})

	t.Run("time range sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 5,
			EndTimeNs:   15,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 11)
		test.RequireProtoEqual(t, envs[4:15], res.Envelopes)
	})

	t.Run("start time sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
			},
		})
		require.NoError(t, err)
		envs = envs[9:]
		utils.Reverse(envs)
		require.Len(t, res.Envelopes, 11)
		test.RequireProtoEqual(t, envs, res.Envelopes)
	})

	t.Run("end time sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			EndTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
			},
		})
		require.NoError(t, err)
		envs = envs[:10]
		utils.Reverse(envs)
		require.Len(t, res.Envelopes, 10)
		test.RequireProtoEqual(t, envs, res.Envelopes)
	})

	t.Run("time range sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 5,
			EndTimeNs:   15,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
			},
		})
		require.NoError(t, err)
		envs = envs[4:15]
		utils.Reverse(envs)
		require.Len(t, res.Envelopes, 11)
		test.RequireProtoEqual(t, envs, res.Envelopes)
	})

	t.Run("limit start time sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Limit: 3,
			},
		})
		require.NoError(t, err)
		envs = envs[9:]
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, envs[:3], res.Envelopes)
	})

	t.Run("limit end time sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			EndTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Limit: 3,
			},
		})
		require.NoError(t, err)
		envs = envs[:10]
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, envs[:3], res.Envelopes)
	})

	t.Run("limit time range sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 5,
			EndTimeNs:   15,
			PagingInfo: &messagev1.PagingInfo{
				Limit: 3,
			},
		})
		require.NoError(t, err)
		envs = envs[4:15]
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, envs[:3], res.Envelopes)
	})

	t.Run("limit start time sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
				Limit:     3,
			},
		})
		require.NoError(t, err)
		envs = envs[9:]
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, envs[:3], res.Envelopes)
	})

	t.Run("limit end time sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			EndTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
				Limit:     3,
			},
		})
		require.NoError(t, err)
		envs = envs[:10]
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, envs[:3], res.Envelopes)
	})

	t.Run("limit time range sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 5,
			EndTimeNs:   15,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
				Limit:     3,
			},
		})
		require.NoError(t, err)
		envs = envs[4:15]
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, envs[:3], res.Envelopes)
	})

	t.Run("limit start time sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
				Limit:     3,
			},
		})
		require.NoError(t, err)
		envs = envs[9:]
		utils.Reverse(envs)
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, envs[:3], res.Envelopes)
	})

	t.Run("limit end time sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			EndTimeNs: 10,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
				Limit:     3,
			},
		})
		require.NoError(t, err)
		envs = envs[:10]
		utils.Reverse(envs)
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, envs[:3], res.Envelopes)
	})

	t.Run("limit time range sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			StartTimeNs: 5,
			EndTimeNs:   15,
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
				Limit:     3,
			},
		})
		require.NoError(t, err)
		envs = envs[5:15]
		utils.Reverse(envs)
		require.Len(t, res.Envelopes, 3)
		test.RequireProtoEqual(t, envs[:3], res.Envelopes)
	})

	t.Run("cursor sort default", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Cursor: &messagev1.Cursor{
					Cursor: &messagev1.Cursor_Index{
						Index: &messagev1.IndexCursor{
							SenderTimeNs: 10,
							Digest:       envCid(t, envs[9]),
						},
					},
				},
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 10)
		test.RequireProtoEqual(t, envs[10:], res.Envelopes)
	})

	t.Run("cursor sort ascending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_ASCENDING,
				Cursor: &messagev1.Cursor{
					Cursor: &messagev1.Cursor_Index{
						Index: &messagev1.IndexCursor{
							SenderTimeNs: 10,
							Digest:       envCid(t, envs[9]),
						},
					},
				},
			},
		})
		require.NoError(t, err)
		require.Len(t, res.Envelopes, 10)
		test.RequireProtoEqual(t, envs[10:], res.Envelopes)
	})

	t.Run("cursor sort descending", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		envs := s.seed(t, 20)
		res, err := s.query(t, &messagev1.QueryRequest{
			PagingInfo: &messagev1.PagingInfo{
				Direction: messagev1.SortDirection_SORT_DIRECTION_DESCENDING,
				Cursor: &messagev1.Cursor{
					Cursor: &messagev1.Cursor_Index{
						Index: &messagev1.IndexCursor{
							SenderTimeNs: 10,
							Digest:       envCid(t, envs[9]),
						},
					},
				},
			},
		})
		require.NoError(t, err)
		envs = envs[:9]
		utils.Reverse(envs)
		require.Len(t, res.Envelopes, 9)
		test.RequireProtoEqual(t, envs, res.Envelopes)
	})
}

func envCid(t *testing.T, env *messagev1.Envelope) multihash.Multihash {
	wrappedEnv, err := types.WrapEnvelope(env)
	require.NoError(t, err)
	return wrappedEnv.Cid
}
