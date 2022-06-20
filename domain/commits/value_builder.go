package commits

import (
	"errors"

	"github.com/steve-care-software/cryptography/domain/hash"
)

type valueBuilder struct {
	hashAdapter hash.Adapter
	namespace   string
	resource    *hash.Hash
	data        []byte
}

func createValueBuilder(
	hashAdapter hash.Adapter,
) ValueBuilder {
	out := valueBuilder{
		hashAdapter: hashAdapter,
		namespace:   "",
		resource:    nil,
		data:        nil,
	}

	return &out
}

// Create initializes the builder
func (app *valueBuilder) Create() ValueBuilder {
	return createValueBuilder(app.hashAdapter)
}

// WithNamespace adds a namespace to the builder
func (app *valueBuilder) WithNamespace(namespace string) ValueBuilder {
	app.namespace = namespace
	return app
}

// WithResource adds a resource to the builder
func (app *valueBuilder) WithResource(resource hash.Hash) ValueBuilder {
	app.resource = &resource
	return app
}

// WithData adds data to the builder
func (app *valueBuilder) WithData(data []byte) ValueBuilder {
	app.data = data
	return app
}

// Now builds a new Value instance
func (app *valueBuilder) Now() (Value, error) {
	if app.data == nil && len(app.data) <= 0 {
		app.data = nil
	}

	if app.data == nil {
		return nil, errors.New("the data is mandatory in order to build a Value instance")
	}

	if app.namespace == "" {
		return nil, errors.New("the namespace is mandatory in order to build a Value instance")
	}

	if app.resource == nil {
		return nil, errors.New("the resource hash is mandatory in order to build a Value instance")
	}

	hash, err := app.hashAdapter.FromMultiBytes([][]byte{
		[]byte(app.namespace),
		app.resource.Bytes(),
		app.data,
	})

	if err != nil {
		return nil, err
	}

	return createValue(
		*hash,
		app.namespace,
		*app.resource,
		app.data,
	), nil
}
