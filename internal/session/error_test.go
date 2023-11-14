package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidStoreKeyErr(t *testing.T) {
	err := NewInvalidStoreKeyErr("NodeProvisionTimeStamp.aks-nodepool1-37632515-vmss0000j4")
	assert.Equal(t, "StoreKey: NodeProvisionTimeStamp.aks-nodepool1-37632515-vmss0000j4 does not exist", err.Error(), "error message does not match")
	assert.Equal(t, true, err.Is(err), "should be equal error valuewise")

	anotherinvalidStoreKey := "NodeProvisionTimeStamp.aks-nodepool2-37632515-vmss0000j4"
	anotherErr := NewInvalidStoreKeyErr(anotherinvalidStoreKey)
	assert.Equal(t, false, err.Is(anotherErr), "should be equal error valuewise")
}
