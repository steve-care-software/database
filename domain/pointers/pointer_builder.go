package pointers

import (
	"errors"
	"strconv"

	"github.com/steve-care-software/cryptography/domain/hash"
)

type pointerBuilder struct {
	hashAdapter hash.Adapter
	namespace   string
	resource    *hash.Hash
	index       *uint
	length      uint
}

func createPointerBuilder(
	hashAdapter hash.Adapter,
) PointerBuilder {
	out := pointerBuilder{
		hashAdapter: hashAdapter,
		namespace:   "",
		resource:    nil,
		index:       nil,
		length:      0,
	}

	return &out
}

// Create initializesthe builder
func (app *pointerBuilder) Create() PointerBuilder {
	return createPointerBuilder(
		app.hashAdapter,
	)
}

// WithNamespace adds a namespace to the builder
func (app *pointerBuilder) WithNamespace(namespace string) PointerBuilder {
	app.namespace = namespace
	return app
}

// WithResource adds a resource to the builder
func (app *pointerBuilder) WithResource(resource hash.Hash) PointerBuilder {
	app.resource = &resource
	return app
}

// WithIndex adds an index to the builder
func (app *pointerBuilder) WithIndex(index uint) PointerBuilder {
	app.index = &index
	return app
}

// WithLength adds a length to the builder
func (app *pointerBuilder) WithLength(length uint) PointerBuilder {
	app.length = length
	return app
}

// Now builds a new Pointer instance
func (app *pointerBuilder) Now() (Pointer, error) {
	if app.namespace == "" {
		return nil, errors.New("the namespace is mandatory in order to build a Pointer instance")
	}

	if app.resource == nil {
		return nil, errors.New("the resource is mandatory in order to build a Pointer instance")
	}

	if app.index == nil {
		return nil, errors.New("the index is mandatory in order to build a Pointer instance")
	}

	if app.length <= 0 {
		return nil, errors.New("the length must be greater than zero (0)in order to build a Pointer instance")
	}

	hash, err := app.hashAdapter.FromMultiBytes([][]byte{
		app.resource.Bytes(),
		[]byte(app.namespace),
		[]byte(strconv.Itoa(int(*app.index))),
		[]byte(strconv.Itoa(int(app.length))),
	})

	if err != nil {
		return nil, err
	}

	return createPointer(
		*hash,
		app.namespace,
		*app.resource,
		*app.index,
		app.length,
	), nil
}
