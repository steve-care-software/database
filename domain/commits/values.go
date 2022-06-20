package commits

import (
	"errors"
	"fmt"

	"github.com/steve-care-software/cryptography/domain/hash"
)

type values struct {
	Hsh     hash.Hash
	Lst     []Value
	mpByRes map[string]Value
}

func createValues(
	hash hash.Hash,
	list []Value,
) Values {
	out := values{
		Hsh:     hash,
		Lst:     list,
		mpByRes: nil,
	}

	return &out
}

// Hash returns the hash
func (obj *values) Hash() hash.Hash {
	return obj.Hsh
}

// List returns the list of values
func (obj *values) List() []Value {
	return obj.Lst
}

// FetchByResource fetches a value by hash
func (obj *values) FetchByResource(res hash.Hash) (Value, error) {
	if obj.mpByRes == nil {
		obj.mpByRes = map[string]Value{}
		for _, oneValue := range obj.Lst {
			obj.mpByRes[oneValue.Resource().String()] = oneValue
		}
	}

	keyname := res.String()
	if ins, ok := obj.mpByRes[keyname]; ok {
		return ins, nil
	}

	str := fmt.Sprintf("there is no value assigned to the given resource (hash: %s)", keyname)
	return nil, errors.New(str)
}
