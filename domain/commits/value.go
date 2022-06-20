package commits

import "github.com/steve-care-software/cryptography/domain/hash"

type value struct {
	Hsh      hash.Hash
	NmeSpace string
	Res      hash.Hash
	Dat      []byte
}

func createValue(
	hash hash.Hash,
	namespace string,
	resource hash.Hash,
	data []byte,
) Value {
	out := value{
		Hsh:      hash,
		NmeSpace: namespace,
		Res:      resource,
		Dat:      data,
	}

	return &out
}

// Hash returns the resource hash
func (obj *value) Hash() hash.Hash {
	return obj.Hsh
}

// Namespace returns the namespace
func (obj *value) Namespace() string {
	return obj.NmeSpace
}

// Resource returns the resource hash
func (obj *value) Resource() hash.Hash {
	return obj.Res
}

// Data returns the data
func (obj *value) Data() []byte {
	return obj.Dat
}
