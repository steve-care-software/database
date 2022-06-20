package resources

import (
	"github.com/steve-care-software/cryptography/domain/hash"
	"github.com/steve-care-software/database/domain/pointers"
)

// NewBuilder creates a new resource builder
func NewBuilder() Builder {
	pointerBuilder := pointers.NewPointerBuilder()
	return createBuilder(pointerBuilder)
}

// Builder represents a resource builder
type Builder interface {
	Create() Builder
	WithNamespace(namespace string) Builder
	WithKey(key hash.Hash) Builder
	WithData(data []byte) Builder
	WithIndex(index uint) Builder
	Now() (Resource, error)
}

// Resource represents a resource
type Resource interface {
	Pointer() pointers.Pointer
	Value() []byte
}

// Repository represents a resource repository
type Repository interface {
	NextIndex() (uint, error)
	Retrieve(ptr pointers.Pointer) (Resource, error)
}
