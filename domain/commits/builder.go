package commits

import (
	"errors"
	"fmt"
	"time"

	"github.com/steve-care-software/cryptography/domain/hash"
)

type builder struct {
	hashAdapter   hash.Adapter
	valueBuilder  ValueBuilder
	valuesBuilder ValuesBuilder
	values        map[string]map[string][]byte
	createdOn     *time.Time
}

func createBuilder(
	hashAdapter hash.Adapter,
	valueBuilder ValueBuilder,
	valuesBuilder ValuesBuilder,
) Builder {
	out := builder{
		hashAdapter:   hashAdapter,
		valueBuilder:  valueBuilder,
		valuesBuilder: valuesBuilder,
		values:        nil,
		createdOn:     nil,
	}

	return &out
}

// Create initializes the builder
func (app *builder) Create() Builder {
	return createBuilder(
		app.hashAdapter,
		app.valueBuilder,
		app.valuesBuilder,
	)
}

// WithValues add values to the builder
func (app *builder) WithValues(values map[string]map[string][]byte) Builder {
	app.values = values
	return app
}

// CreatedOn adds a creation time to the builder
func (app *builder) CreatedOn(createdOn time.Time) Builder {
	app.createdOn = &createdOn
	return app
}

// Now builds a new Commit instance
func (app *builder) Now() (Commit, error) {
	if app.values == nil {
		return nil, errors.New("the values are mandatory in order to build a Commit instance")
	}

	if app.createdOn == nil {
		return nil, errors.New("the creation time is mandatory in order to build a Commit instance")
	}

	list := []Value{}
	for namespace, valueMap := range app.values {
		for resStr, data := range valueMap {
			resource, err := app.hashAdapter.FromString(resStr)
			if err != nil {
				return nil, err
			}

			value, err := app.valueBuilder.Create().WithNamespace(namespace).WithResource(*resource).WithData(data).Now()
			if err != nil {
				return nil, err
			}

			list = append(list, value)
		}
	}

	values, err := app.valuesBuilder.Create().WithList(list).Now()
	if err != nil {
		return nil, err
	}

	hash, err := app.hashAdapter.FromMultiBytes([][]byte{
		values.Hash().Bytes(),
		[]byte(fmt.Sprintf("%d", app.createdOn.UnixNano())),
	})

	if err != nil {
		return nil, err
	}

	return createCommit(*hash, values, app.createdOn.UnixNano()), nil
}
