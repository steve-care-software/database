package transactions

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/steve-care-software/database/domain/commits"
	"github.com/steve-care-software/database/domain/states"
	"github.com/steve-care-software/cryptography/domain/hash"
)

type application struct {
	hashAdapter      hash.Adapter
	commitBuilder    commits.Builder
	commitRepository commits.Repository
	commitService    commits.Service
	stateService     states.Service
	queue            map[string]map[string]map[string][]byte
	commits          map[string]hash.Hash
}

func createApplication(
	hashAdapter hash.Adapter,
	commitBuilder commits.Builder,
	commitRepository commits.Repository,
	commitService commits.Service,
	stateService states.Service,
) Application {
	out := application{
		hashAdapter:      hashAdapter,
		commitBuilder:    commitBuilder,
		commitRepository: commitRepository,
		commitService:    commitService,
		stateService:     stateService,
		queue:            map[string]map[string]map[string][]byte{},
		commits:          map[string]hash.Hash{},
	}

	return &out
}

// Begin creates a context
func (app *application) Begin() (*hash.Hash, error) {
	now := time.Now().UTC().UnixNano()
	str := fmt.Sprintf("%d", now)
	hash, err := app.hashAdapter.FromBytes([]byte(str))
	if err != nil {
		return nil, err
	}

	keyname := hash.String()
	app.queue[keyname] = map[string]map[string][]byte{}
	return hash, nil
}

// Insert inserts a resource to a context
func (app *application) Insert(ctx hash.Hash, namespace string, resource hash.Hash, value []byte) error {
	resCommit := ctx.String()
	if _, ok := app.queue[resCommit]; !ok {
		str := fmt.Sprintf("the commit (hash: %s) does not exists", resCommit)
		return errors.New(str)
	}

	if _, ok := app.queue[resCommit][namespace]; !ok {
		app.queue[resCommit][namespace] = map[string][]byte{}
	}

	app.queue[resCommit][namespace][resource.String()] = value
	return nil
}

// Commit commits a context
func (app *application) Commit(ctx hash.Hash) error {
	resCommit := ctx.String()
	if values, ok := app.queue[resCommit]; ok {
		createdOn := time.Now().UTC()
		commitIns, err := app.commitBuilder.Create().WithValues(values).CreatedOn(createdOn).Now()
		if err != nil {
			return err
		}

		err = app.commitService.Insert(
			commitIns,
			func(ctx commits.Commit) error {
				delete(app.queue, resCommit)
				app.commits[resCommit] = commitIns.Hash()
				return nil
			},
			func(ctx commits.Commit, err error) error {
				log.Printf("there was a problem while inserting commit (hash: %s): %s", ctx.Hash().String(), err.Error())
				return nil
			},
		)

		if err != nil {
			return err
		}

		return nil
	}

	str := fmt.Sprintf("the commit (hash: %s) does not exists", resCommit)
	return errors.New(str)
}

// Queue returns the queue
func (app *application) Queue(ctx hash.Hash) (map[string]map[string][]byte, error) {
	resCommit := ctx.String()
	if values, ok := app.queue[resCommit]; ok {
		return values, nil
	}

	str := fmt.Sprintf("the commit (hash: %s) does not exists", resCommit)
	return nil, errors.New(str)
}

// RollBack rollbacks a commit
func (app *application) RollBack(commit hash.Hash) error {
	retCtx, err := app.commitRepository.Retrieve(commit)
	if err != nil {
		return err
	}

	return app.commitService.Delete(
		retCtx,
		func(ctx commits.Commit) error {
			delete(app.commits, commit.String())
			log.Printf("the rollback was successfully executed on commit (hash: %s)", ctx.Hash().String())
			return nil
		},
		func(ctx commits.Commit, err error) error {
			delete(app.commits, commit.String())
			log.Printf("the rollback failed on commit (hash: %s): %s", ctx.Hash().String(), err.Error())
			return nil
		},
	)
}

// Push pushes a commit to the database
func (app *application) Push(ctx hash.Hash) error {
	keyname := ctx.String()
	if ctxHash, ok := app.commits[keyname]; ok {
		retCtx, err := app.commitRepository.Retrieve(ctxHash)
		if err != nil {
			return err
		}

		return app.stateService.Insert(
			retCtx,
			func(workedCtx commits.Commit) error {
				delete(app.commits, ctx.String())
				return app.commitService.Delete(
					workedCtx,
					func(ctx commits.Commit) error {
						log.Printf("the delete commit (hash: %s) was successful after creating a new state", ctx.Hash().String())
						return nil
					},
					func(ctx commits.Commit, err error) error {
						log.Printf("the rollback failed on commit (hash: %s): %s", ctx.Hash().String(), err.Error())
						return nil
					},
				)
			},
			func(failedCtx commits.Commit, err error) error {
				log.Printf("the state from commit (hash: %s) was expected to be successful but failed: %s", failedCtx.Hash().String(), err.Error())
				return nil
			},
		)
	}

	str := fmt.Sprintf("the commit (hash: %s) does not point to a valid commit", keyname)
	return errors.New(str)
}
