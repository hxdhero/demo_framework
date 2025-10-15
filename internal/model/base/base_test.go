package base

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNullIDFromInt(t *testing.T) {
	ni := NullIDFromInt(3)
	v, err := ni.Value()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), v)
}
