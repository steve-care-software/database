package bytes

import "errors"

type adapterBuilder struct {
	mapping map[string]interface{}
}

func createAdapterBuilder() AdapterBuilder {
	out := adapterBuilder{
		mapping: map[string]interface{}{},
	}

	return &out
}

// Create initializes the builder
func (app *adapterBuilder) Create() AdapterBuilder {
	return createAdapterBuilder()
}

// WithMapping adds a mapping to the builder
func (app *adapterBuilder) WithMapping(mapping map[string]interface{}) AdapterBuilder {
	app.mapping = mapping
	return app
}

// Now builds a new adapter instance
func (app *adapterBuilder) Now() (Adapter, error) {
	if app.mapping == nil {
		return nil, errors.New("the mapping is mandatory in order to build an Adapter instance")
	}

	return createAdapter(app.mapping), nil
}
