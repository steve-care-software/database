package disks

import (
	"errors"
	"path/filepath"

	"github.com/steve-care-software/database/domain/bytes"
	"github.com/steve-care-software/database/domain/commits"
	"github.com/steve-care-software/cryptography/domain/hash"
	"github.com/steve-care-software/database/domain/pointers"
	"github.com/steve-care-software/database/domain/resources"
	"github.com/steve-care-software/database/domain/states"
)

type builder struct {
	hashAdapter     hash.Adapter
	commitAdapter   bytes.Adapter
	stateAdapter    bytes.Adapter
	pointersBuilder pointers.Builder
	resourceBuilder resources.Builder
	statesBuilder   states.Builder
	baseDir         string
	commitDirPath   string
	dbFileName      string
	dbTmpExtension  string
	application     *hash.Hash
}

func createBuilder(
	hashAdapter hash.Adapter,
	commitAdapter bytes.Adapter,
	stateAdapter bytes.Adapter,
	pointersBuilder pointers.Builder,
	resourceBuilder resources.Builder,
	statesBuilder states.Builder,
	baseDir string,
	commitDirPath string,
	dbFileName string,
	dbTmpExtension string,
) Builder {
	out := builder{
		hashAdapter:     hashAdapter,
		commitAdapter:   commitAdapter,
		stateAdapter:    stateAdapter,
		pointersBuilder: pointersBuilder,
		resourceBuilder: resourceBuilder,
		statesBuilder:   statesBuilder,
		baseDir:         baseDir,
		commitDirPath:   commitDirPath,
		dbFileName:      dbFileName,
		dbTmpExtension:  dbTmpExtension,
	}

	return &out
}

// Create initializes  the builder
func (app *builder) Create() Builder {
	return createBuilder(
		app.hashAdapter,
		app.commitAdapter,
		app.stateAdapter,
		app.pointersBuilder,
		app.resourceBuilder,
		app.statesBuilder,
		app.baseDir,
		app.commitDirPath,
		app.dbFileName,
		app.dbTmpExtension,
	)
}

// WithApplication adds an application hash to the builder
func (app *builder) WithApplication(application hash.Hash) Builder {
	app.application = &application
	return app
}

// Now builds a new Application instance
func (app *builder) Now() (commits.Repository, commits.Service, resources.Repository, states.Repository, states.Service, error) {
	if app.application == nil {
		return nil, nil, nil, nil, nil, errors.New("the application hash is mandatory in order to build an Application instance")
	}

	applicationDir := app.application.String()
	commitDirPath := filepath.Join(app.baseDir, applicationDir, app.commitDirPath)
	dbFilePath := filepath.Join(app.baseDir, applicationDir, app.dbFileName)

	// disk repositories:
	stateRepository := createStateRepository(app.stateAdapter, dbFilePath)
	resourceRepository := createResourceRepository(app.hashAdapter, app.resourceBuilder, stateRepository, dbFilePath)
	commitRepository := createCommitRepository(app.hashAdapter, app.commitAdapter, commitDirPath)

	// disk services:
	commitService := createCommitService(app.commitAdapter, commitDirPath)
	stateService := createStateService(app.hashAdapter, app.pointersBuilder, app.resourceBuilder, resourceRepository, app.statesBuilder, app.stateAdapter, stateRepository, dbFilePath, app.dbTmpExtension)

	// return the repositories and services:
	return commitRepository, commitService, resourceRepository, stateRepository, stateService, nil
}
