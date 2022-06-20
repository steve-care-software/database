package disks

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/steve-care-software/database/domain/bytes"
	"github.com/steve-care-software/database/domain/commits"
)

type commitService struct {
	commitAdapter bytes.Adapter
	baseDirPath   string
}

func createCommitService(
	commitAdapter bytes.Adapter,
	baseDirPath string,
) commits.Service {
	out := commitService{
		commitAdapter: commitAdapter,
		baseDirPath:   baseDirPath,
	}

	return &out
}

// Insert inserts a commit instance
func (app *commitService) Insert(commit commits.Commit, worked commits.SuccessCallBackFn, failed commits.FailCallBackFn) error {
	// if the base dir is not created, create it:
	if _, err := os.Stat(app.baseDirPath); os.IsNotExist(err) {
		err := os.MkdirAll(app.baseDirPath, 0777)
		if err != nil {
			return failed(commit, err)
		}
	}

	bytes, err := app.commitAdapter.ToBytes(commit)
	if err != nil {
		return failed(commit, err)
	}

	path := filepath.Join(app.baseDirPath, commit.Hash().String())
	err = ioutil.WriteFile(path, bytes, 0777)
	if err != nil {
		return failed(commit, err)
	}

	err = worked(commit)
	if err != nil {
		return os.Remove(path)
	}

	return nil
}

// Delete deletes a commit instance
func (app *commitService) Delete(commit commits.Commit, worked commits.SuccessCallBackFn, failed commits.FailCallBackFn) error {
	path := filepath.Join(app.baseDirPath, commit.Hash().String())
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return failed(commit, err)
	}

	err = os.Remove(path)
	if err != nil {
		str := fmt.Sprintf("there was an error while deleting the commit file (path: %s): %s", path, err.Error())
		return failed(commit, errors.New(str))
	}

	err = worked(commit)
	if err != nil {
		return ioutil.WriteFile(path, bytes, 0777)
	}

	return nil
}
