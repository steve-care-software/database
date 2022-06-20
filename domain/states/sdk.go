package states

import (
	"time"

	"github.com/steve-care-software/database/domain/commits"
	"github.com/steve-care-software/database/domain/pointers"
	"github.com/steve-care-software/cryptography/domain/hash"
)

const dataLengthErrorPattern = "the remaining data length was expected to be bigger than %d bytes, %d provided"

// SuccessCallBackFn represents a success func callback
type SuccessCallBackFn func(ctx commits.Commit) error

// FailCallBackFn represents a failed func callback
type FailCallBackFn func(ctx commits.Commit, err error) error

// NewMapping returns the pointers conversion mapping
func NewMapping() map[string]interface{} {
	pointersMapping := pointers.NewMapping()
	mp := map[string]interface{}{
		"github.com/steve-care-software/database/domain/states/state": new(state),
	}

	for keyname, value := range pointersMapping {
		mp[keyname] = value
	}

	return mp
}

// NewBuilder creates a new builder instance
func NewBuilder() Builder {
	hashAdapter := hash.NewAdapter()
	return createBuilder(hashAdapter)
}

// Builder represents the state builder
type Builder interface {
	Create() Builder
	WithPointers(ptrs pointers.Pointers) Builder
	WithPrevious(previous State) Builder
	CreatedOn(createdOn time.Time) Builder
	Now() (State, error)
}

// State represents a state
type State interface {
	Hash() hash.Hash
	Height() uint
	Root() State
	Fetch(state hash.Hash) (State, error)
	Pointer(namespace string, resource hash.Hash) (pointers.Pointer, error)
	Pointers() pointers.Pointers
	CreatedOn() time.Time
	HasPrevious() bool
	Previous() State
}

// Repository represents a state repository
type Repository interface {
	Retrieve() (State, uint, error)
}

// Service represents a pointer service
type Service interface {
	Insert(commit commits.Commit, worked SuccessCallBackFn, failed FailCallBackFn) error
}
