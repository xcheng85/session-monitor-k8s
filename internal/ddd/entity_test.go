package ddd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	entity := NewEntity("id", "name")
	assert.Equal(t, "id", entity.ID())
	assert.Equal(t, "name", entity.EntityName())
	assert.Equal(t, true, entity.Equals(NewEntity("id", "name")))
}
