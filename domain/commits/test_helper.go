package commits

import (
	"time"

	"github.com/steve-care-software/cryptography/domain/hash"
)

// NewCommitForTests creates a new commit for tests
func NewCommitForTests(raw map[string][][]byte) Commit {
	values := map[string]map[string][]byte{}
	hashAdapter := hash.NewAdapter()
	for namespace, oneData := range raw {
		values[namespace] = map[string][]byte{}
		for _, oneValue := range oneData {
			hash, err := hashAdapter.FromBytes(oneValue)
			if err != nil {
				panic(err)
			}

			values[namespace][hash.String()] = oneValue
		}

	}

	createdOn := time.Now().UTC()
	context, err := NewBuilder().Create().CreatedOn(createdOn).WithValues(values).Now()
	if err != nil {
		panic(err)
	}

	return context
}
