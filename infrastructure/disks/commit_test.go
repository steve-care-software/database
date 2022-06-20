package disks

import (
	"os"
	"testing"

	"github.com/steve-care-software/database/domain/commits"
	"github.com/steve-care-software/cryptography/domain/hash"
)

func TestCommit_Success(t *testing.T) {
	baseDir := "./test_files"
	commitDirPath := "commits"
	dbFileName := "database.db"
	dbTmpExtension := ".tmp"
	defer func() {
		os.RemoveAll(baseDir)
	}()

	application, err := hash.NewAdapter().FromBytes([]byte("this is some data"))
	if err != nil {
		panic(err)
	}

	commit := commits.NewCommitForTests(map[string][][]byte{
		"my_namespace": [][]byte{
			[]byte("this is the first element"),
			[]byte("this is the second element"),
			[]byte("yes, this is the last element"),
		},
	})

	commitRepository, commitService, _, _, _, err := NewBuilder(baseDir, commitDirPath, dbFileName, dbTmpExtension).Create().WithApplication(*application).Now()
	if err != nil {
		panic(err)
	}

	err = commitService.Insert(
		commit,
		func(ctx commits.Commit) error {
			if !ctx.Hash().Compare(commit.Hash()) {
				t.Errorf("the commit is invalid")
				return nil
			}

			return nil
		},
		func(ctx commits.Commit, err error) error {
			t.Errorf("the execution was expected to work, error returned: %s", err.Error())
			return nil
		},
	)

	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	retList, err := commitRepository.List()
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	if len(retList) != 1 {
		t.Errorf("the list was expected to contain %d elements, %d returned", 1, len(retList))
		return
	}

	if !retList[0].Compare(commit.Hash()) {
		t.Errorf("the commit list element (hash: %s) was expected to contain the commit hash: %s", retList[0].String(), commit.Hash().String())
		return
	}

	retCommit, err := commitRepository.Retrieve(commit.Hash())
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	if !retCommit.Hash().Compare(commit.Hash()) {
		t.Errorf("the returned commit is invalid")
		return
	}

	err = commitService.Delete(
		commit,
		func(ctx commits.Commit) error {
			if !ctx.Hash().Compare(commit.Hash()) {
				t.Errorf("the commit is invalid")
				return nil
			}

			return nil
		},
		func(ctx commits.Commit, err error) error {
			t.Errorf("the execution was expected to work, error returned: %s", err.Error())
			return nil
		},
	)

	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}
}
