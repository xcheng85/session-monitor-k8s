package domain

import (
	"fmt"
	"reflect"
)

type BadNodeLabelErr struct {
	label *map[string]string
}

func (r *BadNodeLabelErr) Error() string {
	return fmt.Sprintf("label: %T is incorrect", r.label)
}

// assert style in golang
func (s *BadNodeLabelErr) Is(target error) bool {
	targetErr, ok := target.(*BadNodeLabelErr)
	if !ok {
		return false
	}
	return reflect.DeepEqual(*s.label, *targetErr.label)

}

func NewBadNodeLabelErr(label *map[string]string) *BadNodeLabelErr {
	return &BadNodeLabelErr{label}
}
