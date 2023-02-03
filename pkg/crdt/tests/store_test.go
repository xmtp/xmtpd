package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Query(t *testing.T) {
	// create a topic with 20 messages
	net := RandomMsgTest(t, 1, 1, 20)
	defer net.Close()

	t.Run("all", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0)
		require.NoError(t, err)
		assert.Nil(t, pi, "paging info")
		net.AssertQueryResult(t, res, intRange(1, 20)...)
	})
	t.Run("descending", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, Descending())
		require.NoError(t, err)
		assert.NotNil(t, pi, "paging info")
		net.AssertQueryResult(t, res, intRange(20, 1)...)
	})
	t.Run("limit", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, Limit(5))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		net.AssertQueryCursor(t, 5, pi.Cursor)
		net.AssertQueryResult(t, res, intRange(1, 5)...)
	})
	t.Run("limit descending", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, Limit(5), Descending())
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		net.AssertQueryCursor(t, 16, pi.Cursor)
		net.AssertQueryResult(t, res, intRange(20, 16)...)
	})
	t.Run("range", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, TimeRange(5, 13))
		require.NoError(t, err)
		assert.Nil(t, pi, "paging info")
		net.AssertQueryResult(t, res, 5, 6, 7, 8, 9, 10, 11, 12, 13)

	})
	t.Run("range descending", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, TimeRange(5, 9), Descending())
		require.NoError(t, err)
		assert.NotNil(t, pi, "paging info")
		net.AssertQueryResult(t, res, 9, 8, 7, 6, 5)

	})
	t.Run("range limit", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, TimeRange(5, 15), Limit(4))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		net.AssertQueryCursor(t, 8, pi.Cursor)
		net.AssertQueryResult(t, res, 5, 6, 7, 8)

	})
	t.Run("range limit descending", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, TimeRange(5, 15), Limit(4), Descending())
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		net.AssertQueryCursor(t, 12, pi.Cursor)
		net.AssertQueryResult(t, res, 15, 14, 13, 12)

	})
	t.Run("cursor", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, TimeRange(5, 13), Limit(5))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		net.AssertQueryCursor(t, 9, pi.Cursor)
		net.AssertQueryResult(t, res, 5, 6, 7, 8, 9)

		res, pi, err = net.Query(t, 0, t0, TimeRange(5, 13), Limit(5), Cursor(pi.Cursor))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		net.AssertQueryCursor(t, 13, pi.Cursor)
		net.AssertQueryResult(t, res, 10, 11, 12, 13)

		res, pi, err = net.Query(t, 0, t0, TimeRange(5, 13), Limit(5), Cursor(pi.Cursor))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		assert.Nil(t, pi.Cursor)
		net.AssertQueryResult(t, res)

	})
	t.Run("cursor descending", func(t *testing.T) {
		res, pi, err := net.Query(t, 0, t0, TimeRange(7, 15), Limit(5), Descending())
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		assert.NotNil(t, pi.Cursor)
		net.AssertQueryResult(t, res, 15, 14, 13, 12, 11)

		res, pi, err = net.Query(t, 0, t0, TimeRange(7, 15), Limit(5), Descending(), Cursor(pi.Cursor))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		assert.NotNil(t, pi.Cursor)
		net.AssertQueryResult(t, res, 10, 9, 8, 7)

		res, pi, err = net.Query(t, 0, t0, TimeRange(7, 15), Limit(5), Descending(), Cursor(pi.Cursor))
		require.NoError(t, err)
		require.NotNil(t, pi, "paging info")
		assert.Nil(t, pi.Cursor)
		net.AssertQueryResult(t, res)
	})
}
