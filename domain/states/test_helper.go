package states

import (
	"time"

	"github.com/steve-care-software/database/domain/pointers"
)

// NewStateForTests creates a new state for tests
func NewStateForTests(hasPrevious bool) State {
	createdOn := time.Now().UTC()
	pointers, _ := pointers.NewPointersForTests()
	builder := NewBuilder().Create().CreatedOn(createdOn).WithPointers(pointers)
	if hasPrevious {
		prev := NewStateForTests(false)
		builder.WithPrevious(prev)
	}

	state, err := builder.Now()
	if err != nil {
		panic(err)
	}

	return state
}
