package commits

import (
	"testing"

	"github.com/steve-care-software/database/domain/bytes"
)

func TestAdapter_Success(t *testing.T) {
	commit := NewCommitForTests(map[string][][]byte{
		"my_namespace": [][]byte{
			[]byte("this is the first element"),
			[]byte("this is the second element"),
			[]byte("yes, this is the last element"),
		},
	})

	adapter, err := bytes.NewAdapterBuilder().Create().WithMapping(NewMapping()).Now()
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	bytes, err := adapter.ToBytes(commit)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	retCommit, _, err := adapter.ToInstance(bytes)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	casted := retCommit.(Commit)
	if !commit.Hash().Compare(casted.Hash()) {
		t.Errorf("the commit has was expected to be %s, %s returned", commit.Hash().String(), casted.Hash().String())
		return
	}
}
