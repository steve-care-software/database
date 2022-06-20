package disks

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	domain_bytes "github.com/steve-care-software/database/domain/bytes"
	"github.com/steve-care-software/database/domain/commits"
	"github.com/steve-care-software/cryptography/domain/hash"
	"github.com/steve-care-software/database/domain/pointers"
	"github.com/steve-care-software/database/domain/resources"
	"github.com/steve-care-software/database/domain/states"
)

type stateService struct {
	hashAdapter        hash.Adapter
	pointersBuilder    pointers.Builder
	resourceBuilder    resources.Builder
	resourceRepository resources.Repository
	builder            states.Builder
	adapter            domain_bytes.Adapter
	repository         states.Repository
	databaseFilePath   string
	tmpExtension       string
	mutex              sync.Mutex
}

func createStateService(
	hashAdapter hash.Adapter,
	pointersBuilder pointers.Builder,
	resourceBuilder resources.Builder,
	resourceRepository resources.Repository,
	builder states.Builder,
	adapter domain_bytes.Adapter,
	repository states.Repository,
	databaseFilePath string,
	tmpExtension string,
) states.Service {
	out := stateService{
		hashAdapter:        hashAdapter,
		pointersBuilder:    pointersBuilder,
		resourceBuilder:    resourceBuilder,
		resourceRepository: resourceRepository,
		builder:            builder,
		adapter:            adapter,
		repository:         repository,
		databaseFilePath:   databaseFilePath,
		tmpExtension:       tmpExtension,
	}

	return &out
}

// Insert inserts a state instance from the passed commit
func (app *stateService) Insert(commit commits.Commit, worked states.SuccessCallBackFn, failed states.FailCallBackFn) error {
	// if the database directory does not exists, create it:
	resDir := filepath.Dir(app.databaseFilePath)
	if _, err := os.Stat(resDir); os.IsNotExist(err) {
		err := os.MkdirAll(resDir, 0777)
		if err != nil {
			return failed(commit, err)
		}
	}

	// if the database file does not exists, create it:
	if _, err := os.Stat(app.databaseFilePath); errors.Is(err, os.ErrNotExist) {
		err = ioutil.WriteFile(app.databaseFilePath, []byte{}, 0777)
		if err != nil {
			return failed(commit, err)
		}
	}

	// create the state instance:
	state, resources, prevStateSizeInBytes, err := app.createStateInstance(commit)
	if err != nil {
		return failed(commit, err)
	}

	// convert the state to bytes:
	stateBytes, err := app.adapter.ToBytes(state)
	if err != nil {
		return failed(commit, err)
	}

	var stateSize uint64 = uint64(len(stateBytes))
	stateSizeBuf := new(bytes.Buffer)
	err = binary.Write(stateSizeBuf, binary.LittleEndian, stateSize)
	if err != nil {
		return failed(commit, err)
	}

	// open the tmp database file:
	fin, err := os.Open(app.databaseFilePath)
	if err != nil {
		return failed(commit, err)
	}
	defer fin.Close()

	// open the output tmp file:
	resTmpPath := fmt.Sprintf("%s.%s", app.databaseFilePath, app.tmpExtension)
	fout, err := os.OpenFile(resTmpPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return failed(commit, err)
	}

	defer fout.Close()
	defer func() {
		os.Remove(resTmpPath)
	}()

	// add the state bytes at the beginning of the file:
	allStateBytes := stateSizeBuf.Bytes()
	allStateBytes = append(allStateBytes, stateBytes...)
	_, err = fout.Write(allStateBytes)
	if err != nil {
		return failed(commit, err)
	}

	// if there is a previous state:
	if prevStateSizeInBytes > 0 {
		// offset the original state data from the database file:
		_, err = fin.Seek(int64(prevStateSizeInBytes), io.SeekStart)
		if err != nil {
			return failed(commit, err)
		}

		// copy the original resource data to the tmp file:
		_, err = io.Copy(fout, fin)
		if err != nil {
			return failed(commit, err)
		}
	}

	// combine the resources:
	resBytes := []byte{}
	for _, oneResource := range resources {
		resBytes = append(resBytes, oneResource.Pointer().Resource().Bytes()...)
		resBytes = append(resBytes, oneResource.Value()...)
	}

	// append the new resource to the tmp database file:
	if _, err = fout.Write(resBytes); err != nil {
		return failed(commit, err)
	}

	// lock the mutex during the rename file operations and unlock when we exit the fn:
	app.mutex.Lock()
	defer app.mutex.Unlock()

	// execute the worked callback:
	err = worked(commit)
	if err != nil {
		return err
	}

	// rename and replace the tmp database file for the real resource file:
	return os.Rename(resTmpPath, app.databaseFilePath)
}

func (app *stateService) createStateInstance(commit commits.Commit) (states.State, []resources.Resource, uint, error) {
	nextIndex, err := app.resourceRepository.NextIndex()
	if err != nil {
		return nil, nil, 0, err
	}

	ptrList := []pointers.Pointer{}
	resources := []resources.Resource{}
	values := commit.Values().List()
	for _, oneValue := range values {
		namespace := oneValue.Namespace()
		resource := oneValue.Resource()
		data := oneValue.Data()
		res, err := app.resourceBuilder.Create().WithNamespace(namespace).WithKey(resource).WithData(data).WithIndex(nextIndex).Now()
		if err != nil {
			return nil, nil, 0, err
		}

		ptr := res.Pointer()
		nextIndex = ptr.Index() + ptr.Length()
		resources = append(resources, res)
		ptrList = append(ptrList, ptr)
	}

	ptrs, err := app.pointersBuilder.Create().WithList(ptrList).Now()
	if err != nil {
		return nil, nil, 0, err
	}

	createdOn := time.Now().UTC()
	builder := app.builder.Create().WithPointers(ptrs).CreatedOn(createdOn)
	prev, sizeInBytes, _ := app.repository.Retrieve()
	if prev != nil {
		builder.WithPrevious(prev)
	}

	ins, err := builder.Now()
	if err != nil {
		return nil, nil, 0, err
	}

	return ins, resources, sizeInBytes, nil
}
