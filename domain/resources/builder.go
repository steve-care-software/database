package resources

import (
	"errors"

	"github.com/steve-care-software/database/domain/pointers"
	"github.com/steve-care-software/cryptography/domain/hash"
)

type builder struct {
	pointerBuilder pointers.PointerBuilder
	namespace      string
	key            *hash.Hash
	data           []byte
	index          *uint
}

func createBuilder(
	pointerBuilder pointers.PointerBuilder,
) Builder {
	out := builder{
		pointerBuilder: pointerBuilder,
		namespace:      "",
		key:            nil,
		data:           nil,
		index:          nil,
	}

	return &out
}

// Create initializes the builder
func (app *builder) Create() Builder {
	return createBuilder(
		app.pointerBuilder,
	)
}

// WithNamespace adds a namespace to the builder
func (app *builder) WithNamespace(namespace string) Builder {
	app.namespace = namespace
	return app
}

// WithKey adds a key to the builder
func (app *builder) WithKey(key hash.Hash) Builder {
	app.key = &key
	return app
}

// WithData adds data to the builder
func (app *builder) WithData(data []byte) Builder {
	app.data = data
	return app
}

// WithIndex adds an index to the builder
func (app *builder) WithIndex(index uint) Builder {
	app.index = &index
	return app
}

// Now builds a new Resource instance
func (app *builder) Now() (Resource, error) {
	if app.data == nil {
		return nil, errors.New("the data is mandatory in order to build a Resource instance")
	}

	if app.index == nil {
		return nil, errors.New("the index is mandatory in order to build a Resource instance")
	}

	if app.key == nil {
		return nil, errors.New("the key is mandatory in order to build a Resource instance")
	}

	if app.namespace == "" {
		return nil, errors.New("the namespace is mandatory in order to build a Resource instance")
	}

	length := uint(len(app.key.Bytes()) + len(app.data))
	pointer, err := app.pointerBuilder.Create().WithNamespace(app.namespace).WithResource(*app.key).WithIndex(*app.index).WithLength(length).Now()
	if err != nil {
		return nil, err
	}

	return createResource(pointer, app.data), nil
}
