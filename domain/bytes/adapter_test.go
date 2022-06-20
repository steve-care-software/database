package bytes

import (
	"reflect"
	"testing"
)

type testStruct struct {
	IsTrue    bool
	IsFalse   bool
	First     uint
	Second    float32
	Third     float64
	Fourth    uint8
	Fifth     uint16
	Sixth     uint32
	Seventh   uint64
	Height    int8
	Ninth     int16
	Tenth     int32
	Eleventh  int64
	Twelve    int
	Thirteen  *testSecondStruct
	Fourteen  testSecondStruct
	Fifteen   []uint8
	Sixteen   []testSecondStruct
	Seventeen string
	Heighteen []uint
}

type testSecondStruct struct {
	First  uint
	Second uint32
}

func TestAdapter_Success(t *testing.T) {
	ins := testStruct{
		IsTrue:   true,
		IsFalse:  false,
		First:    uint(^uint32(0)),
		Second:   float32(45.78),
		Third:    float64(567.765),
		Fourth:   uint8(22),
		Fifth:    uint16(45234),
		Sixth:    uint32(234523345),
		Seventh:  uint64(1234123412313),
		Height:   int8(22),
		Ninth:    int16(-4123),
		Tenth:    int32(-456),
		Eleventh: int64(23452345235),
		Twelve:   int(-234523452345234),
		Thirteen: &testSecondStruct{
			First:  uint(43),
			Second: uint32(242),
		},
		Fourteen: testSecondStruct{
			First:  uint(65),
			Second: uint32(3456),
		},
		Fifteen: []uint8{1, 2, 3},
		Sixteen: []testSecondStruct{
			testSecondStruct{
				First:  uint(1),
				Second: uint32(456),
			},
		},
		Seventeen: "voila",
		Heighteen: []uint{
			1,
			2,
			3,
		},
	}

	adapter, err := NewAdapterBuilder().Create().WithMapping(map[string]interface{}{
		"github.com/steve-care-software/database/domain/bytes/testStruct":       testStruct{},
		"github.com/steve-care-software/database/domain/bytes/testSecondStruct": testSecondStruct{},
		"[]bytes.testSecondStruct": testSecondStruct{},
		"[]uint8":                  uint8(0),
		"[]uint":                   uint(0),
	}).Now()
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	data, err := adapter.ToBytes(ins)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	retIns, _, err := adapter.ToInstance(data)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	if !reflect.DeepEqual(ins, retIns) {
		t.Errorf("the returned instance is invalid, \nexpected: %v, \nreturned: %v\n", ins, retIns)
		return
	}

}
