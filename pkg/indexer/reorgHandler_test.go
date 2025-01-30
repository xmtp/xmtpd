package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_blockRange(t *testing.T) {
	tests := []struct {
		name           string
		from           uint64
		wantStartBlock uint64
		wantEndBlock   uint64
	}{
		{
			name:           "block range with subtraction",
			from:           1001,
			wantStartBlock: 1,
			wantEndBlock:   1001,
		},
		{
			name:           "block range without subtraction",
			from:           500,
			wantStartBlock: 0,
			wantEndBlock:   500,
		},
		{
			name:           "block range zero",
			from:           0,
			wantStartBlock: 0,
			wantEndBlock:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startBlock, endBlock := blockRange(tt.from)
			assert.Equal(t, tt.wantStartBlock, startBlock)
			assert.Equal(t, tt.wantEndBlock, endBlock)
		})
	}
}
