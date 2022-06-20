package resources

import (
	"fmt"
	"math/rand"
	"time"
)

// NewResourcesForTests creates a new resources for tests
func NewResourcesForTests() ([]Resource, int) {
	list := []Resource{}
	amount := 20
	for i := 0; i < amount; i++ {
		list = append(list, NewResourceForTests())
	}

	return list, amount
}

// NewResourceForTests creates a new resource for tests
func NewResourceForTests() Resource {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	index := uint(r1.Intn(3452734845))
	data := []byte(fmt.Sprintf("this is some data: %d", index))
	res, err := NewBuilder().Create().WithData(data).WithIndex(index).Now()
	if err != nil {
		panic(err)
	}

	return res
}
