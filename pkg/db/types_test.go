package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xmtp/xmtpd/pkg/db"
)

func TestFillMissingOriginators_NoneAdded(t *testing.T) {
	vc := db.VectorClock{100: 10, 200: 20, 300: 30}
	db.FillMissingOriginators(vc, []int32{100, 200, 300})
	assert.Equal(t, db.VectorClock{100: 10, 200: 20, 300: 30}, vc)
}

func TestFillMissingOriginators_SomeMissing(t *testing.T) {
	vc := db.VectorClock{100: 10}
	db.FillMissingOriginators(vc, []int32{100, 200, 300})
	assert.Len(t, vc, 3)
	assert.Equal(t, uint64(10), vc[100])
	assert.Equal(t, uint64(0), vc[200])
	assert.Equal(t, uint64(0), vc[300])
}

func TestFillMissingOriginators_AllMissing(t *testing.T) {
	vc := db.VectorClock{}
	db.FillMissingOriginators(vc, []int32{100, 200})
	assert.Equal(t, db.VectorClock{100: 0, 200: 0}, vc)
}

func TestFillMissingOriginators_EmptyAll(t *testing.T) {
	vc := db.VectorClock{}
	db.FillMissingOriginators(vc, nil)
	assert.Empty(t, vc)
}
