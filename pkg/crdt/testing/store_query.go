package crdttest

import (
	"testing"

	"github.com/xmtp/xmtpd/pkg/api"
)

func RunStoreQueryTests(t *testing.T, topic string, storeMaker TestStoreMaker) {
	t.Helper()
	s := storeMaker(t)
	defer s.Close()
	// seed with events with timestamps 1, 2, ..., 20
	s.Seed(t, topic, 20)

	t.Run("all sort default", func(t *testing.T) {
		res := s.query(t, topic)
		requireResultEqual(t, res, 1, 20)
		requireNoCursor(t, res)
	})

	t.Run("all sort ascending", func(t *testing.T) {
		res := s.query(t, topic, api.Ascending())
		requireResultEqual(t, res, 1, 20)
		requireNoCursor(t, res)
	})

	t.Run("all sort descending", func(t *testing.T) {
		res := s.query(t, topic, api.Descending())
		requireResultEqual(t, res, 20, 1)
		requireNoCursor(t, res)
	})

	t.Run("limit sort default", func(t *testing.T) {
		res := s.query(t, topic, api.Limit(5))
		requireResultEqual(t, res, 1, 5)
		requireResultCursor(t, res, 5)
	})

	t.Run("limit sort ascending", func(t *testing.T) {
		res := s.query(t, topic, api.Limit(5), api.Ascending())
		requireResultEqual(t, res, 1, 5)
		requireResultCursor(t, res, 5)
	})

	t.Run("limit sort descending", func(t *testing.T) {
		res := s.query(t, topic, api.Limit(5), api.Descending())
		requireResultEqual(t, res, 20, 16)
		requireResultCursor(t, res, 16)
	})

	t.Run("start time sort default", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(10, 0))
		requireResultEqual(t, res, 10, 20)
		requireNoCursor(t, res)
	})

	t.Run("end time sort default", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(0, 10))
		requireResultEqual(t, res, 1, 10)
		requireNoCursor(t, res)
	})

	t.Run("time range sort default", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(5, 15))
		requireResultEqual(t, res, 5, 15)
		requireNoCursor(t, res)
	})

	t.Run("start time sort ascending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(10, 0), api.Ascending())
		requireResultEqual(t, res, 10, 20)
		requireNoCursor(t, res)
	})

	t.Run("end time sort ascending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(0, 10), api.Ascending())
		requireResultEqual(t, res, 1, 10)
		requireNoCursor(t, res)
	})

	t.Run("time range sort ascending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(5, 15), api.Ascending())
		requireResultEqual(t, res, 5, 15)

		requireNoCursor(t, res)
	})

	t.Run("start time sort descending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(10, 0), api.Descending())
		requireResultEqual(t, res, 20, 10)
		requireNoCursor(t, res)
	})

	t.Run("end time sort descending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(0, 10), api.Descending())
		requireResultEqual(t, res, 10, 1)

		requireNoCursor(t, res)
	})

	t.Run("time range sort descending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(5, 15), api.Descending())
		requireResultEqual(t, res, 15, 5)

		requireNoCursor(t, res)
	})

	t.Run("limit start time sort default", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(10, 0), api.Limit(3))
		requireResultEqual(t, res, 10, 12)
		requireResultCursor(t, res, 12)
	})

	t.Run("limit end time sort default", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(0, 10), api.Limit(3))
		requireResultEqual(t, res, 1, 3)
		requireResultCursor(t, res, 3)
	})

	t.Run("limit time range sort default", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(5, 15), api.Limit(3))
		requireResultEqual(t, res, 5, 7)
		requireResultCursor(t, res, 7)
	})

	t.Run("limit start time sort ascending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(10, 0), api.Limit(3), api.Ascending())
		requireResultEqual(t, res, 10, 12)
		requireResultCursor(t, res, 12)
	})

	t.Run("limit end time sort ascending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(0, 10), api.Limit(3), api.Ascending())
		requireResultEqual(t, res, 1, 3)
		requireResultCursor(t, res, 3)
	})

	t.Run("limit time range sort ascending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(5, 15), api.Limit(3), api.Ascending())
		requireResultEqual(t, res, 5, 7)
		requireResultCursor(t, res, 7)
	})

	t.Run("limit start time sort descending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(10, 0), api.Limit(3), api.Descending())
		requireResultEqual(t, res, 20, 18)
		requireResultCursor(t, res, 18)
	})

	t.Run("limit end time sort descending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(0, 10), api.Limit(3), api.Descending())
		requireResultEqual(t, res, 10, 8)
		requireResultCursor(t, res, 8)
	})

	t.Run("limit time range sort descending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(5, 15), api.Limit(3), api.Descending())
		requireResultEqual(t, res, 15, 13)
		requireResultCursor(t, res, 13)
	})

	t.Run("cursor sort default", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(5, 13), api.Limit(4))
		requireResultEqual(t, res, 5, 8)
		requireResultCursor(t, res, 8)
		res = s.query(t, topic, api.TimeRange(5, 13), api.Limit(4), api.Cursor(res))
		requireResultEqual(t, res, 9, 12)
		requireResultCursor(t, res, 12)
		res = s.query(t, topic, api.TimeRange(5, 13), api.Limit(4), api.Cursor(res))
		requireResultEqual(t, res, 13, 13)
		requireNoCursor(t, res)
	})

	t.Run("cursor sort ascending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(5, 13), api.Limit(4), api.Ascending())
		requireResultEqual(t, res, 5, 8)
		requireResultCursor(t, res, 8)
		res = s.query(t, topic, api.TimeRange(5, 13), api.Limit(4), api.Ascending(), api.Cursor(res))
		requireResultEqual(t, res, 9, 12)
		requireResultCursor(t, res, 12)
		res = s.query(t, topic, api.TimeRange(5, 13), api.Limit(4), api.Ascending(), api.Cursor(res))
		requireResultEqual(t, res, 13, 13)
		requireNoCursor(t, res)
	})

	t.Run("cursor sort descending", func(t *testing.T) {
		res := s.query(t, topic, api.TimeRange(7, 15), api.Limit(4), api.Descending())
		requireResultEqual(t, res, 15, 12)
		requireResultCursor(t, res, 12)
		res = s.query(t, topic, api.TimeRange(7, 15), api.Limit(4), api.Descending(), api.Cursor(res))
		requireResultEqual(t, res, 11, 8)
		requireResultCursor(t, res, 8)
		res = s.query(t, topic, api.TimeRange(7, 15), api.Limit(4), api.Descending(), api.Cursor(res))
		requireResultEqual(t, res, 7, 7)
		requireNoCursor(t, res)
	})
}
