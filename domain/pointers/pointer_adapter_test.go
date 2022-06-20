package pointers

import (
	"testing"

	"github.com/steve-care-software/database/domain/bytes"
)

func TestPointerAdapter_Success(t *testing.T) {
	pointer := NewPointerForTests()
	adapter, err := bytes.NewAdapterBuilder().Create().WithMapping(NewPointerMapping()).Now()
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	data, err := adapter.ToBytes(pointer)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	retPointer, remaining, err := adapter.ToInstance(data)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	casted := retPointer.(Pointer)
	if !pointer.Hash().Compare(casted.Hash()) {
		t.Errorf("the pointer hash was expected to be %s, %s returned", pointer.Hash().String(), casted.Hash().String())
		return
	}

	if len(remaining) != 0 {
		t.Errorf("the remaining data was expected to be empty, %d bytes returned", len(remaining))
		return
	}
}
