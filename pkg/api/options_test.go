package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAddr(t *testing.T) {
	assert.NoError(t, validateAddr("localhost", 0))
	assert.NoError(t, validateAddr("0.0.0.0", 0))
}
