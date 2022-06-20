package commits

import (
	"time"

	"github.com/steve-care-software/cryptography/domain/hash"
)

const dataLengthErrorPattern = "the remaining data length was expected to be bigger than %d bytes, %d provided"

// SuccessCallBackFn represents a success func callback
type SuccessCallBackFn func(ctx Commit) error

// FailCallBackFn represents a failed func callback
type FailCallBackFn func(ctx Commit, err error) error

// NewMapping returns the conversion mapping
func NewMapping() map[string]interface{} {
	mp := map[string]interface{}{
		"github.com/steve-care-software/database/domain/commits/commit": new(commit),
		"github.com/steve-care-software/database/domain/commits/values": new(values),
		"github.com/steve-care-software/database/domain/commits/value":  new(value),
		"[]commits.Value": new(Value),
		"[]uint8":         uint8(0),
		"hash.Hash":       uint8(0),
	}

	return mp
}

// NewBuilder creates a new builder instance
func NewBuilder() Builder {
	hashAdapter := hash.NewAdapter()
	valueBuilder := NewValueBuilder()
	valuesBuilder := NewValuesBuilder()
	return createBuilder(hashAdapter, valueBuilder, valuesBuilder)
}

// NewValuesBuilder creates a new values builder
func NewValuesBuilder() ValuesBuilder {
	hashAdapter := hash.NewAdapter()
	return createValuesBuilder(hashAdapter)
}

// NewValueBuilder creates a new value builder
func NewValueBuilder() ValueBuilder {
	hashAdapter := hash.NewAdapter()
	return createValueBuilder(hashAdapter)
}

// Builder represents a commit builder
type Builder interface {
	Create() Builder
	WithValues(values map[string]map[string][]byte) Builder
	CreatedOn(createdOn time.Time) Builder
	Now() (Commit, error)
}

// Commit represents a commit
type Commit interface {
	Hash() hash.Hash
	Values() Values
	CreatedOn() time.Time
}

// ValuesBuilder represents the values builder
type ValuesBuilder interface {
	Create() ValuesBuilder
	WithList(list []Value) ValuesBuilder
	Now() (Values, error)
}

// Values represents the values
type Values interface {
	Hash() hash.Hash
	List() []Value
	FetchByResource(res hash.Hash) (Value, error)
}

// ValueBuilder represents the value builder
type ValueBuilder interface {
	Create() ValueBuilder
	WithNamespace(namespace string) ValueBuilder
	WithResource(resource hash.Hash) ValueBuilder
	WithData(data []byte) ValueBuilder
	Now() (Value, error)
}

// Value represents a value
type Value interface {
	Hash() hash.Hash
	Namespace() string
	Resource() hash.Hash
	Data() []byte
}

// Repository represents a context repository
type Repository interface {
	List() ([]hash.Hash, error)
	Retrieve(hash hash.Hash) (Commit, error)
}

// Service represents a context service
type Service interface {
	Insert(ctx Commit, worked SuccessCallBackFn, failed FailCallBackFn) error
	Delete(ctx Commit, worked SuccessCallBackFn, failed FailCallBackFn) error
}
