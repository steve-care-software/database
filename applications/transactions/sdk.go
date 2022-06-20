package transactions

import (
	"github.com/steve-care-software/database/domain/commits"
	"github.com/steve-care-software/database/domain/states"
	"github.com/steve-care-software/cryptography/domain/hash"
)

// NewApplication creates a new application instance
/*func NewApplication(
	commitRepository commits.Repository,
	commitService commits.Service,
	stateService states.Service,
) Application {
	hashAdapter := hash.NewAdapter()
	commitBuilder := commits.NewBuilder()
	lexerApp := lexers.NewApplication()
	return createApplication(
		hashAdapter,
		commitBuilder,
		commitRepository,
		commitService,
		stateService,
		lexerApp,
	)
}*/

// Builder represents an application builder
type Builder interface {
	Create() Builder
	WithCommitRepository(commitRepository commits.Repository) Builder
	WithCommitService(commitService commits.Service) Builder
	WithStateService(stateService states.Service) Builder
	Now() (Application, error)
}

// Application represents a transaction application
type Application interface {
	Begin() (*hash.Hash, error)
	Insert(context hash.Hash, namespace string, resource hash.Hash, value []byte) error
	Commit(context hash.Hash) error
	Queue(context hash.Hash) (map[string]map[string][]byte, error)
	RollBack(context hash.Hash) error
	Push(context hash.Hash) error
}
