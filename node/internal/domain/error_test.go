package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadNodeLabelErr(t *testing.T) {
	err := NewBadNodeLabelErr(nil)
	assert.Equal(t, "label: *map[string]string is incorrect", err.Error(), "error message does not match")
}
