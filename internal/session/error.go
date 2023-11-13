package session

import (
	"fmt"
)

type InvalidStoreKeyErr struct {
	key string
}

func (r *InvalidStoreKeyErr) Error() string {
	return fmt.Sprintf("StoreKey: %s is incorrect", r.key)
}

// assert style in golang
func (s *InvalidStoreKeyErr) Is(target error) bool {
	targetErr, ok := target.(*InvalidStoreKeyErr)
	if !ok {
		return false
	}
	return s.key == targetErr.key
}

func NewInvalidStoreKeyErr(key string) *InvalidStoreKeyErr {
	return &InvalidStoreKeyErr{key}
}
