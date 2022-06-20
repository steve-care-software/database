package disks

import (
	"github.com/steve-care-software/database/domain/commits"
	"github.com/steve-care-software/database/domain/pointers"
	"github.com/steve-care-software/database/domain/resources"
	"github.com/steve-care-software/database/domain/states"
	"github.com/steve-care-software/database/domain/bytes"
	"github.com/steve-care-software/cryptography/domain/hash"
)

const dataLengthErrorPattern = "the remaining data length was expected to be bigger than %d bytes, %d provided"

// NewBuilder creates anewdisk builder
func NewBuilder(
	baseDirPath string,
	commitDirPath string,
	dbFileName string,
	dbTmpExtension string,
) Builder {
	hashAdapter := hash.NewAdapter()
	pointersBuilder := pointers.NewBuilder()
	resourceBuilder := resources.NewBuilder()
	statesBuilder := states.NewBuilder()
	commitAdapter, err := bytes.NewAdapterBuilder().Create().WithMapping(commits.NewMapping()).Now()
	if err != nil {
		panic(err)
	}

	stateAdapter, err := bytes.NewAdapterBuilder().Create().WithMapping(states.NewMapping()).Now()
	if err != nil {
		panic(err)
	}

	return createBuilder(
		hashAdapter,
		commitAdapter,
		stateAdapter,
		pointersBuilder,
		resourceBuilder,
		statesBuilder,
		baseDirPath,
		commitDirPath,
		dbFileName,
		dbTmpExtension,
	)
}

// Builder represents the disk builder
type Builder interface {
	Create() Builder
	WithApplication(application hash.Hash) Builder
	Now() (commits.Repository, commits.Service, resources.Repository, states.Repository, states.Service, error)
}
