package disks

import (
	"errors"
	"fmt"
	"os"

	"github.com/steve-care-software/cryptography/domain/hash"
	"github.com/steve-care-software/database/domain/pointers"
	"github.com/steve-care-software/database/domain/resources"
	"github.com/steve-care-software/database/domain/states"
)

type resourceRepository struct {
	hashAdapter      hash.Adapter
	resourceBuilder  resources.Builder
	stateRepository  states.Repository
	databaseFilePath string
}

func createResourceRepository(
	hashAdapter hash.Adapter,
	resourceBuilder resources.Builder,
	stateRepository states.Repository,
	databaseFilePath string,
) resources.Repository {
	out := resourceRepository{
		hashAdapter:      hashAdapter,
		resourceBuilder:  resourceBuilder,
		stateRepository:  stateRepository,
		databaseFilePath: databaseFilePath,
	}

	return &out
}

// NextIndex returns the pointer next index
func (app *resourceRepository) NextIndex() (uint, error) {
	// if the database file does not exists, return 0:
	if _, err := os.Stat(app.databaseFilePath); errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}

	// if we can't retrieve the state, it means the file is empty, so the nextIndex is 0:
	_, stateSize, err := app.stateRepository.Retrieve()
	if err != nil {
		return 0, nil
	}

	// open the file:
	filePtr, err := os.Open(app.databaseFilePath)
	if err != nil {
		return 0, err
	}

	defer filePtr.Close()

	file, err := filePtr.Stat()
	if err != nil {
		return 0, err
	}

	fileSize := uint(file.Size())
	if fileSize < stateSize {
		str := fmt.Sprintf("the file size (%d bytes) cannot be smaller than the stateSize (%d bytes)", fileSize, stateSize)
		return 0, errors.New(str)
	}

	return fileSize - stateSize, nil
}

// Retrieve retrieves a resource from a pointer
func (app *resourceRepository) Retrieve(ptr pointers.Pointer) (resources.Resource, error) {
	_, stateSize, err := app.stateRepository.Retrieve()
	if err != nil {
		return nil, err
	}

	// open the file:
	filePtr, err := os.Open(app.databaseFilePath)
	if err != nil {
		return nil, err
	}

	defer filePtr.Close()

	offset := stateSize + ptr.Index()
	length := ptr.Length()
	resData := make([]byte, length, length)
	_, err = filePtr.ReadAt(resData, int64(offset))
	if err != nil {
		return nil, err
	}

	if len(resData) <= hash.Size {
		str := fmt.Sprintf(dataLengthErrorPattern, hash.Size, len(resData))
		return nil, errors.New(str)
	}

	keyBytes := resData[:hash.Size]
	key, err := app.hashAdapter.FromBytes(keyBytes)
	if err != nil {
		return nil, err
	}

	ptrIndex := ptr.Index()
	namespace := ptr.Namespace()
	return app.resourceBuilder.Create().WithNamespace(namespace).WithKey(*key).WithData(resData[hash.Size:]).WithIndex(ptrIndex).Now()
}
