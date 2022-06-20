package pointers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/steve-care-software/cryptography/domain/hash"
)

// NewPointersForTests creates a new pointers for tests
func NewPointersForTests() (Pointers, int) {
	list := []Pointer{}
	amount := 20
	for i := 0; i < amount; i++ {
		list = append(list, NewPointerForTests())
	}

	pointers, err := NewBuilder().Create().WithList(list).Now()
	if err != nil {
		panic(err)
	}

	return pointers, amount
}

// NewPointerForTests creates a new pointer for tests
func NewPointerForTests() Pointer {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	index := uint(r1.Intn(3452734845))
	length := uint(r1.Intn(3452734845)) + 1

	str := fmt.Sprintf("this is some resource, idx: %d, length: %d", index, length)
	resource, err := hash.NewAdapter().FromBytes([]byte(str))
	if err != nil {
		panic(err)
	}

	namespace := "my_namespace"
	pointer, err := NewPointerBuilder().Create().WithNamespace(namespace).WithResource(*resource).WithIndex(index).WithLength(length).Now()
	if err != nil {
		panic(err)
	}

	return pointer
}
