package resources

import "github.com/steve-care-software/database/domain/pointers"

type resource struct {
	pointer pointers.Pointer
	value   []byte
}

func createResource(
	pointer pointers.Pointer,
	value []byte,
) Resource {
	out := resource{
		pointer: pointer,
		value:   value,
	}

	return &out
}

// Pointer returns the pointer
func (obj *resource) Pointer() pointers.Pointer {
	return obj.pointer
}

// Value returns the value
func (obj *resource) Value() []byte {
	return obj.value
}
