package states

import (
	"testing"

	"github.com/steve-care-software/database/domain/bytes"
)

func TestAdapter_Success(t *testing.T) {
	state := NewStateForTests(true)
	adapter, err := bytes.NewAdapterBuilder().Create().WithMapping(NewMapping()).Now()
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	data, err := adapter.ToBytes(state)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	retState, remaining, err := adapter.ToInstance(data)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	casted := retState.(State)
	if !state.Hash().Compare(casted.Hash()) {
		t.Errorf("the state hash was expected to be %s, %s returned", state.Hash().String(), casted.Hash().String())
		return
	}

	if len(remaining) > 0 {
		t.Errorf("the remaining []byte were expected to empty")
		return
	}
}
