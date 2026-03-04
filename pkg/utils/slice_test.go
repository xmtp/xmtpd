package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChunkSlice_EmptySliceSizeZero(t *testing.T) {
	result := ChunkSlice([]int{}, 0)
	require.NotNil(t, result)
	assert.Empty(t, result)
}

func TestChunkSlice_EmptySliceSizeOne(t *testing.T) {
	result := ChunkSlice([]int{}, 1)
	require.NotNil(t, result)
	assert.Empty(t, result)
}

func TestChunkSlice_SingleElementSizeOne(t *testing.T) {
	result := ChunkSlice([]int{42}, 1)
	require.NotNil(t, result)
	assert.Equal(t, [][]int{{42}}, result)
}

func TestChunkSlice_EvenlyDivisible(t *testing.T) {
	result := ChunkSlice([]int{1, 2, 3, 4, 5, 6}, 2)
	require.NotNil(t, result)
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}}, result)
}

func TestChunkSlice_NotEvenlyDivisible(t *testing.T) {
	result := ChunkSlice([]int{1, 2, 3, 4, 5}, 2)
	require.NotNil(t, result)
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, result)
}

func TestChunkSlice_SizeLargerThanSlice(t *testing.T) {
	result := ChunkSlice([]int{1, 2, 3}, 10)
	require.NotNil(t, result)
	assert.Equal(t, [][]int{{1, 2, 3}}, result)
}

func TestChunkSlice_SizeZeroWithElements(t *testing.T) {
	result := ChunkSlice([]int{1, 2, 3}, 0)
	require.NotNil(t, result)
	assert.Equal(t, [][]int{{1, 2, 3}}, result)
}

func TestChunkSlice_NegativeSizeWithElements(t *testing.T) {
	result := ChunkSlice([]int{1, 2, 3}, -1)
	require.NotNil(t, result)
	assert.Equal(t, [][]int{{1, 2, 3}}, result)
}

func TestChunkSlice_NilSlice(t *testing.T) {
	result := ChunkSlice[int](nil, 2)
	require.NotNil(t, result)
	assert.Empty(t, result)
}
