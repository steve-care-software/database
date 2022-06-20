package pointers

import "github.com/steve-care-software/cryptography/domain/hash"

type pointer struct {
	Hsh      hash.Hash
	NmeSpace string
	Res      hash.Hash
	Idx      uint
	Lgth     uint
}

func createPointer(
	hash hash.Hash,
	namespace string,
	resource hash.Hash,
	index uint,
	length uint,
) Pointer {
	out := pointer{
		Hsh:      hash,
		NmeSpace: namespace,
		Res:      resource,
		Idx:      index,
		Lgth:     length,
	}

	return &out
}

// Hash returns the hash
func (obj *pointer) Hash() hash.Hash {
	return obj.Hsh
}

// Namespace returns the namespace
func (obj *pointer) Namespace() string {
	return obj.NmeSpace
}

// Resource returns the resource
func (obj *pointer) Resource() hash.Hash {
	return obj.Res
}

// Index returns the index
func (obj *pointer) Index() uint {
	return obj.Idx
}

// Length returns the length
func (obj *pointer) Length() uint {
	return obj.Lgth
}
