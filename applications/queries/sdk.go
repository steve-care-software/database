package queries

import (
	"github.com/steve-care-software/database/domain/commits"
	"github.com/steve-care-software/database/domain/pointers"
	"github.com/steve-care-software/database/domain/resources"
	"github.com/steve-care-software/database/domain/states"
	"github.com/steve-care-software/cryptography/domain/hash"
)

// NewApplication creates a new application instance
/*func NewApplication(
	resRepository resources.Repository,
	commitRepository commits.Repository,
	stateRepository states.Repository,
) Application {
	return createApplication(
		resRepository,
		commitRepository,
		stateRepository,
	)
}*/

// Builder represents an application builder
type Builder interface {
	Create() Builder
	WithResourceRepository(resRepository resources.Repository) Builder
	WithCommitRepository(commitRepository commits.Repository) Builder
	WithStateRepository(stateRepository states.Repository) Builder
	Now() (Application, error)
}

// Application represents a query application
type Application interface {
	Head() (states.State, error)
	State(hash hash.Hash) (states.State, error)
	Commits() ([]hash.Hash, error)
	Commit(hash hash.Hash) (commits.Commit, error)
	Resource(ptr pointers.Pointer) (resources.Resource, error)
}
