package commits

import (
	"errors"

	"github.com/steve-care-software/cryptography/domain/hash"
)

type valuesBuilder struct {
	hashAdapter hash.Adapter
	list        []Value
}

func createValuesBuilder(
	hashAdapter hash.Adapter,
) ValuesBuilder {
	out := valuesBuilder{
		hashAdapter: hashAdapter,
		list:        nil,
	}

	return &out
}

// Create initializes the builder
func (app *valuesBuilder) Create() ValuesBuilder {
	return createValuesBuilder(app.hashAdapter)
}

// WithLists adds a list to the builder
func (app *valuesBuilder) WithList(list []Value) ValuesBuilder {
	app.list = list
	return app
}

// Now builds a new Values instance
func (app *valuesBuilder) Now() (Values, error) {
	if app.list == nil && len(app.list) <= 0 {
		app.list = nil
	}

	if app.list == nil {
		return nil, errors.New("there must be at least 1 Value instance in order to build a Values instance")
	}

	data := [][]byte{}
	for _, oneValue := range app.list {
		data = append(data, oneValue.Hash().Bytes())
	}

	hash, err := app.hashAdapter.FromMultiBytes(data)
	if err != nil {
		return nil, err
	}

	return createValues(*hash, app.list), nil
}
