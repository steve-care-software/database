package states

import (
	"errors"
	"fmt"
	"time"

	"github.com/steve-care-software/database/domain/pointers"
	"github.com/steve-care-software/cryptography/domain/hash"
)

type builder struct {
	hashAdapter hash.Adapter
	ptrs        pointers.Pointers
	createdOn   *time.Time
	previous    State
}

func createBuilder(
	hashAdapter hash.Adapter,
) Builder {
	out := builder{
		hashAdapter: hashAdapter,
		ptrs:        nil,
		createdOn:   nil,
		previous:    nil,
	}

	return &out
}

// Create initializes the builder
func (app *builder) Create() Builder {
	return createBuilder(app.hashAdapter)
}

// WithPointers add pointers to the builder
func (app *builder) WithPointers(ptrs pointers.Pointers) Builder {
	app.ptrs = ptrs
	return app
}

// WithPrevious adds a previous state to the builder
func (app *builder) WithPrevious(previous State) Builder {
	app.previous = previous
	return app
}

// CreatedOn adds a creation time to the builder
func (app *builder) CreatedOn(createdOn time.Time) Builder {
	app.createdOn = &createdOn
	return app
}

// Now builds a new State instance
func (app *builder) Now() (State, error) {
	if app.ptrs == nil {
		return nil, errors.New("the pointers are mandatory in order to build a State instance")
	}

	if app.createdOn == nil {
		return nil, errors.New("the creation time is mandatory in order to build a State instance")
	}

	data := [][]byte{
		app.ptrs.Hash().Bytes(),
		[]byte(fmt.Sprintf("%d", app.createdOn.UnixNano())),
	}

	if app.previous != nil {
		data = append(data, app.previous.Hash().Bytes())
	}

	hash, err := app.hashAdapter.FromMultiBytes(data)
	if err != nil {
		return nil, err
	}

	if app.previous != nil {
		return createStateWithPrevious(*hash, app.ptrs, app.createdOn.UnixNano(), app.previous), nil
	}

	return createState(*hash, app.ptrs, app.createdOn.UnixNano()), nil
}
