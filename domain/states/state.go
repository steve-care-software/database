package states

import (
	"errors"
	"fmt"
	"time"

	"github.com/steve-care-software/cryptography/domain/hash"
	"github.com/steve-care-software/database/domain/pointers"
)

type state struct {
	Hsh  hash.Hash
	Ptrs pointers.Pointers
	CrOn int64
	Prev State
}

func createState(
	hash hash.Hash,
	ptrs pointers.Pointers,
	createdOn int64,
) State {
	return createStateInternally(hash, ptrs, createdOn, nil)
}

func createStateWithPrevious(
	hash hash.Hash,
	ptrs pointers.Pointers,
	createdOn int64,
	previous State,
) State {
	return createStateInternally(hash, ptrs, createdOn, previous)
}

func createStateInternally(
	hash hash.Hash,
	ptrs pointers.Pointers,
	createdOn int64,
	previous State,
) State {
	out := state{
		Hsh:  hash,
		Ptrs: ptrs,
		CrOn: createdOn,
		Prev: previous,
	}

	return &out
}

// Hash returns the hash
func (obj *state) Hash() hash.Hash {
	return obj.Hsh
}

// Height returns the state height
func (obj *state) Height() uint {
	if obj.HasPrevious() {
		return obj.Previous().Height() + 1
	}

	return 1
}

// Root returns the root state
func (obj *state) Root() State {
	if !obj.HasPrevious() {
		return obj
	}

	return obj.Previous().Root()
}

// Fetch fetches a state by hash
func (obj *state) Fetch(state hash.Hash) (State, error) {
	if state.Compare(obj.Hash()) {
		return obj, nil
	}

	if obj.HasPrevious() {
		return obj.Previous().Fetch(state)
	}

	str := fmt.Sprintf("the previous state (hash: %s) could not be found", state.String())
	return nil, errors.New(str)
}

// Pointer fetches a pointer by hash
func (obj *state) Pointer(namespace string, resource hash.Hash) (pointers.Pointer, error) {
	if obj.Ptrs.Exists(namespace, resource) {
		return obj.Ptrs.Fetch(namespace, resource)
	}

	if obj.HasPrevious() {
		return obj.Previous().Pointer(namespace, resource)
	}

	str := fmt.Sprintf("the resource (namespace: %s, hash: %s) does not contain a matching pointer", namespace, resource.String())
	return nil, errors.New(str)
}

// Pointers returns the pointers
func (obj *state) Pointers() pointers.Pointers {
	return obj.Ptrs
}

// CreatedOn returns the creation time
func (obj *state) CreatedOn() time.Time {
	return time.Unix(0, obj.CrOn)
}

// HasPrevious returns true if there is a previous state, false otherwise
func (obj *state) HasPrevious() bool {
	return obj.Prev != nil
}

// Previous returns the previous state, if any
func (obj *state) Previous() State {
	return obj.Prev
}
