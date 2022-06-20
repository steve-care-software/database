package commits

import (
	"time"

	"github.com/steve-care-software/cryptography/domain/hash"
)

type commit struct {
	Hsh  hash.Hash
	Vals Values
	CrOn int64
}

func createCommit(
	hash hash.Hash,
	values Values,
	createdOn int64,
) Commit {
	out := commit{
		Hsh:  hash,
		Vals: values,
		CrOn: createdOn,
	}

	return &out
}

// Hash returns the hash
func (obj *commit) Hash() hash.Hash {
	return obj.Hsh
}

// Values returns the values
func (obj *commit) Values() Values {
	return obj.Vals
}

// CreatedOn returns the creation time
func (obj *commit) CreatedOn() time.Time {
	return time.Unix(0, obj.CrOn)
}
