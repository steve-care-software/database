package pointers

import (
	"errors"
	"fmt"

	"github.com/steve-care-software/cryptography/domain/hash"
)

type pointers struct {
	Hsh hash.Hash
	Lst []Pointer
	mp  map[string]map[string]Pointer
}

func createPointers(
	hash hash.Hash,
	list []Pointer,
) Pointers {
	out := pointers{
		Hsh: hash,
		Lst: list,
		mp:  nil,
	}

	return &out
}

// Hash returns the hash
func (obj *pointers) Hash() hash.Hash {
	return obj.Hsh
}

// List returns the list of pointers
func (obj *pointers) List() []Pointer {
	return obj.Lst
}

func (obj *pointers) initMp() {
	if obj.mp != nil {
		return
	}

	obj.mp = map[string]map[string]Pointer{}
	for _, onePointer := range obj.Lst {
		namespace := onePointer.Namespace()
		if _, ok := obj.mp[namespace]; !ok {
			obj.mp[namespace] = map[string]Pointer{}
		}

		keyname := onePointer.Resource().String()
		obj.mp[namespace][keyname] = onePointer
	}
}

// Exists returns true if the pointer exists, false otherwise
func (obj *pointers) Exists(namespace string, resource hash.Hash) bool {
	obj.initMp()
	if resources, ok := obj.mp[namespace]; ok {
		keyname := resource.String()
		if _, ok := resources[keyname]; ok {
			return true
		}

		return false
	}

	return false
}

// Fetch fetches the pointer, if any
func (obj *pointers) Fetch(namespace string, resource hash.Hash) (Pointer, error) {
	obj.initMp()
	if resources, ok := obj.mp[namespace]; ok {
		keyname := resource.String()
		if ins, ok := resources[keyname]; ok {
			return ins, nil
		}

		str := fmt.Sprintf("the resource (namespace: %s, hash: %s) does not contain a matching pointer", namespace, keyname)
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the namespace (hash: %s) does not exists", namespace)
	return nil, errors.New(str)
}
