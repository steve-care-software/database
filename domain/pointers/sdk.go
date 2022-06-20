package pointers

import (
	"github.com/steve-care-software/cryptography/domain/hash"
)

const dataLengthErrorPattern = "the remaining data length was expected to be bigger than %d bytes, %d provided"

// NewMapping returns the pointers conversion mapping
func NewMapping() map[string]interface{} {
	pointerMapping := NewPointerMapping()
	mp := map[string]interface{}{
		"github.com/steve-care-software/database/domain/pointers/pointers": new(pointers),
		"[]pointers.Pointer": new(Pointer),
	}

	for keyname, value := range pointerMapping {
		mp[keyname] = value
	}

	return mp
}

// NewPointerMapping returns the pointer conversion mapping
func NewPointerMapping() map[string]interface{} {
	mp := map[string]interface{}{
		"github.com/steve-care-software/database/domain/pointers/pointer": new(pointer),
		"hash.Hash": uint8(0),
	}

	return mp
}

// NewBuilder creates a new builder instance
func NewBuilder() Builder {
	hashAdapter := hash.NewAdapter()
	return createBuilder(hashAdapter)
}

// NewPointerBuilder creates a new pointer builder
func NewPointerBuilder() PointerBuilder {
	hashAdapter := hash.NewAdapter()
	return createPointerBuilder(hashAdapter)
}

// Builder represents a pointers builder
type Builder interface {
	Create() Builder
	WithList(list []Pointer) Builder
	Now() (Pointers, error)
}

// Pointers represents pointers
type Pointers interface {
	Hash() hash.Hash
	List() []Pointer
	Fetch(namespace string, resource hash.Hash) (Pointer, error)
	Exists(namespace string, resource hash.Hash) bool
}

// PointerBuilder represents a pointer builder
type PointerBuilder interface {
	Create() PointerBuilder
	WithNamespace(namespace string) PointerBuilder
	WithResource(resource hash.Hash) PointerBuilder
	WithIndex(index uint) PointerBuilder
	WithLength(length uint) PointerBuilder
	Now() (Pointer, error)
}

// Pointer represents a pointer
type Pointer interface {
	Hash() hash.Hash
	Namespace() string
	Resource() hash.Hash
	Index() uint
	Length() uint
}
