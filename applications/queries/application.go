package queries

import (
	"github.com/steve-care-software/database/domain/commits"
	"github.com/steve-care-software/database/domain/pointers"
	"github.com/steve-care-software/database/domain/resources"
	"github.com/steve-care-software/database/domain/states"
	"github.com/steve-care-software/cryptography/domain/hash"
)

type application struct {
	resRepository    resources.Repository
	commitRepository commits.Repository
	stateRepository  states.Repository
}

func createApplication(
	resRepository resources.Repository,
	commitRepository commits.Repository,
	stateRepository states.Repository,
) Application {
	out := application{
		resRepository:    resRepository,
		commitRepository: commitRepository,
		stateRepository:  stateRepository,
	}

	return &out
}

// Head returns the head state
func (app *application) Head() (states.State, error) {
	ins, _, err := app.stateRepository.Retrieve()
	if err != nil {
		return nil, err
	}

	return ins, nil
}

// State returns the state by hash
func (app *application) State(hash hash.Hash) (states.State, error) {
	head, err := app.Head()
	if err != nil {
		return nil, err
	}

	return head.Fetch(hash)
}

// Commits returns the commits list
func (app *application) Commits() ([]hash.Hash, error) {
	return app.commitRepository.List()
}

// Commit returns the commit by hash
func (app *application) Commit(hash hash.Hash) (commits.Commit, error) {
	return app.commitRepository.Retrieve(hash)
}

// Resource returns the resource by pointer
func (app *application) Resource(ptr pointers.Pointer) (resources.Resource, error) {
	return app.resRepository.Retrieve(ptr)
}
