package disks

import (
	"encoding/binary"
	"errors"
	"os"

	"github.com/steve-care-software/database/domain/bytes"
	"github.com/steve-care-software/database/domain/states"
)

type stateRepository struct {
	stateAdapter     bytes.Adapter
	databaseFilePath string
}

func createStateRepository(
	stateAdapter bytes.Adapter,
	databaseFilePath string,
) states.Repository {
	out := stateRepository{
		stateAdapter:     stateAdapter,
		databaseFilePath: databaseFilePath,
	}

	return &out
}

// Retrieve returns the head state
func (app *stateRepository) Retrieve() (states.State, uint, error) {
	// if the database file does not exists, return nil:
	if _, err := os.Stat(app.databaseFilePath); errors.Is(err, os.ErrNotExist) {
		return nil, 0, nil
	}

	// open the file:
	ptr, err := os.Open(app.databaseFilePath)
	if err != nil {
		return nil, 0, err
	}

	defer ptr.Close()

	// read the first 8 bytes:
	stateSizeLength := 8
	stateSizeInBytes := make([]byte, stateSizeLength, stateSizeLength)
	_, err = ptr.Read(stateSizeInBytes)
	if err != nil {
		return nil, 0, err
	}

	stateSize := binary.LittleEndian.Uint64(stateSizeInBytes)

	// read the state:
	stateBytes := make([]byte, stateSize, stateSize)
	_, err = ptr.Read(stateBytes)
	if err != nil {
		return nil, 0, err
	}

	// converts the bytes to a state instance:
	state, _, err := app.stateAdapter.ToInstance(stateBytes)
	if err != nil {
		return nil, 0, nil
	}

	if casted, ok := state.(states.State); ok {
		return casted, uint(stateSize) + uint(stateSizeLength), nil
	}

	return nil, 0, errors.New("the State []byte could not be casted properly")
}
