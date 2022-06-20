package disks

import (
	"bytes"
	"os"
	"testing"

	"github.com/steve-care-software/database/domain/commits"
	"github.com/steve-care-software/cryptography/domain/hash"
)

func TestState_Success(t *testing.T) {
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
			[]byte("1) this is the first element"),
			[]byte("1) this is the second element"),
			[]byte("1) yes, this is the last element"),
		},
	})

	_, _, resourceRepository, stateRepository, stateService, err := NewBuilder(baseDir, commitDirPath, dbFileName, dbTmpExtension).Create().WithApplication(*application).Now()
	if err != nil {
		panic(err)
	}

	err = stateService.Insert(
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

	// retrieve the state head:
	stateHead, _, err := stateRepository.Retrieve()
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	if stateHead.HasPrevious() {
		t.Errorf("the state was NOT expecting a previous state")
		return
	}

	ptrs := stateHead.Pointers().List()
	values := commit.Values().List()
	if len(values) != len(ptrs) {
		t.Errorf("%d pointers were expected in the state, %d returned", ptrs, len(values))
		return
	}

	if stateHead.Height() != 1 {
		t.Errorf("the state head was expected to be %d, %d returned", 1, stateHead.Height())
		return
	}

	res, err := resourceRepository.Retrieve(ptrs[1])
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	fetchedVal, err := commit.Values().FetchByResource(res.Pointer().Resource())
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	if bytes.Compare(res.Value(), fetchedVal.Data()) != 0 {
		t.Errorf("the resource bytes do not match")
		return
	}

	secondCommit := commits.NewCommitForTests(map[string][][]byte{
		"my_namespace": [][]byte{
			[]byte("2) this is the first element"),
			[]byte("2) this is the second element"),
			[]byte("2) yes, this is the last element"),
			[]byte("2) yet again!!!"),
		},
	})

	// insert some more commits:
	err = stateService.Insert(
		secondCommit,
		func(ctx commits.Commit) error {
			if !ctx.Hash().Compare(secondCommit.Hash()) {
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

	// retrieve the state head:
	secondStateHead, _, err := stateRepository.Retrieve()
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	if secondStateHead.Height() != 2 {
		t.Errorf("the state head was expected to be %d, %d returned", 2, secondStateHead.Height())
		return
	}

	if !secondStateHead.HasPrevious() {
		t.Errorf("the state was expected a previous state, none provided")
		return
	}

	prevPointer, err := secondStateHead.Pointer("my_namespace", ptrs[2].Resource())
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	resAgain, err := resourceRepository.Retrieve(prevPointer)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	fetchedValAgain, err := commit.Values().FetchByResource(resAgain.Pointer().Resource())
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	if bytes.Compare(resAgain.Value(), fetchedValAgain.Data()) != 0 {
		t.Errorf("the resource bytes do not match")
		return
	}

	secondPtrs := secondStateHead.Pointers().List()
	secondValues := secondCommit.Values().List()
	if len(secondValues) != len(secondPtrs) {
		t.Errorf("%d pointers were expected in the state, %d returned", secondPtrs, len(secondValues))
		return
	}

	lastPointer, err := secondStateHead.Pointer("my_namespace", secondPtrs[3].Resource())
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	resLast, err := resourceRepository.Retrieve(lastPointer)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	fetchedSecondVal, err := secondCommit.Values().FetchByResource(resLast.Pointer().Resource())
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	if bytes.Compare(resLast.Value(), fetchedSecondVal.Data()) != 0 {
		t.Errorf("the resource bytes do not match")
		return
	}
}
