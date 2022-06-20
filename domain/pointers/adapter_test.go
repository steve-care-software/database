package pointers

import (
	"testing"

	"github.com/steve-care-software/database/domain/bytes"
)

func TestAdapter_Success(t *testing.T) {
	pointers, amount := NewPointersForTests()
	adapter, err := bytes.NewAdapterBuilder().Create().WithMapping(NewMapping()).Now()
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	data, err := adapter.ToBytes(pointers)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	retPointers, remaining, err := adapter.ToInstance(data)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	retCasted := retPointers.(Pointers)
	retList := retCasted.List()
	if len(retList) != amount {
		t.Errorf("%d pointer instances were expected, %d returned", amount, len(retList))
		return
	}

	if !pointers.Hash().Compare(retCasted.Hash()) {
		t.Errorf("the pointers hash was expected to be %s, %s returned", pointers.Hash().String(), retCasted.Hash().String())
		return
	}

	if len(remaining) != 0 {
		t.Errorf("the remaining data was expected to be empty, %d bytes returned", len(remaining))
		return
	}
}

func TestAdapter_withRemaining_Success(t *testing.T) {
	suffix := []byte("this is some data")
	pointers, amount := NewPointersForTests()
	adapter, err := bytes.NewAdapterBuilder().Create().WithMapping(NewMapping()).Now()
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	data, err := adapter.ToBytes(pointers)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	data = append(data, suffix...)
	retPointers, remaining, err := adapter.ToInstance(data)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	casted := retPointers.(Pointers)
	retList := casted.List()
	if len(retList) != amount {
		t.Errorf("%d pointer instances were expected, %d returned", amount, len(retList))
		return
	}

	if !pointers.Hash().Compare(casted.Hash()) {
		t.Errorf("the pointers hash was expected to be %s, %s returned", pointers.Hash().String(), casted.Hash().String())
		return
	}

	if len(remaining) != len(suffix) {
		t.Errorf("the remaining data was expected to contain %d bytes, %d bytes returned", len(suffix), len(remaining))
		return
	}
}
