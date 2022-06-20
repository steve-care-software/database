package disks

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/steve-care-software/database/domain/commits"
	"github.com/steve-care-software/database/domain/bytes"
	"github.com/steve-care-software/cryptography/domain/hash"
)

type commitRepository struct {
	hashAdapter   hash.Adapter
	commitAdapter bytes.Adapter
	baseDirPath   string
}

func createCommitRepository(
	hashAdapter hash.Adapter,
	commitAdapter bytes.Adapter,
	baseDirPath string,
) commits.Repository {
	out := commitRepository{
		hashAdapter:   hashAdapter,
		commitAdapter: commitAdapter,
		baseDirPath:   baseDirPath,
	}

	return &out
}

// List lists the commits
func (app *commitRepository) List() ([]hash.Hash, error) {
	// if the base dir is not created, return an empty list:
	if _, err := os.Stat(app.baseDirPath); os.IsNotExist(err) {
		return []hash.Hash{}, nil
	}

	// read the dir content:
	files, err := ioutil.ReadDir(app.baseDirPath)
	if err != nil {
		return nil, err
	}

	list := []hash.Hash{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		str := file.Name()
		hash, err := app.hashAdapter.FromString(str)
		if err != nil {
			return nil, err
		}

		list = append(list, *hash)
	}

	return list, nil
}

// Retrieve retrieves a commit by hash
func (app *commitRepository) Retrieve(hash hash.Hash) (commits.Commit, error) {
	path := filepath.Join(app.baseDirPath, hash.String())
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		str := fmt.Sprintf("there is no commit for the given hash: %s", hash.String())
		return nil, errors.New(str)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ins, _, err := app.commitAdapter.ToInstance(bytes)
	if err != nil {
		return nil, err
	}

	if casted, ok := ins.(commits.Commit); ok {
		return casted, nil
	}

	str := fmt.Sprintf("the retrieved commit (hash: %s) could not be casted properly", hash.String())
	return nil, errors.New(str)
}
