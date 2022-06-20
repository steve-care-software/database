package bytes

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
)

type adapter struct {
	mapping map[string]interface{}
}

func createAdapter(
	mapping map[string]interface{},
) Adapter {
	out := adapter{
		mapping: mapping,
	}

	return &out
}

// ToBytes converts an instance to bytes
func (app *adapter) ToBytes(ins interface{}) ([]byte, error) {
	value := reflect.ValueOf(ins)
	return app.valueToBytes(value)
}

// ToInstance converts bytes to an instance
func (app *adapter) ToInstance(bytes []byte) (interface{}, []byte, error) {
	pVal, remaining, err := app.toInstance(bytes, false)
	if err != nil {
		return nil, nil, err
	}

	return pVal.Interface(), remaining, nil
}

func (app *adapter) toInstance(bytes []byte, isPtr bool) (*reflect.Value, []byte, error) {
	if len(bytes) < 1 {
		str := fmt.Sprintf(bytesLengthTooSmallErr, 1, len(bytes))
		return nil, nil, errors.New(str)
	}

	remaining := bytes[1:]
	if bytes[0]&Bool != 0 {
		if len(remaining) < 1 {
			str := fmt.Sprintf(bytesLengthTooSmallErr, 1, len(remaining))
			return nil, nil, errors.New(str)
		}

		value := remaining[0:1][0]
		if value == 0 {
			val := reflect.ValueOf(false)
			return &val, remaining[1:], nil
		}

		val := reflect.ValueOf(true)
		return &val, remaining[1:], nil
	}

	if bytes[0]&String != 0 {
		return app.bytesToString(remaining)
	}

	if bytes[0]&Ptr != 0 {
		if len(remaining) <= 0 {
			value := reflect.ValueOf(nil)
			return &value, nil, nil
		}

		return app.toInstance(remaining, true)
	}

	if bytes[0]&Struct != 0 {
		return app.bytesToStruct(remaining, isPtr)
	}

	if bytes[0]&Array != 0 {
		return app.bytesToArray(remaining)
	}

	if bytes[0]&Uint != 0 {
		return app.bytesToUint(remaining)
	}

	if bytes[0]&Int != 0 {
		return app.bytesToInt(remaining)
	}

	if bytes[0]&Float != 0 {
		return app.bytesToFloat(remaining)
	}

	str := fmt.Sprintf("the given data (%v) could not be converted to an instance", bytes)
	return nil, nil, errors.New(str)
}

func (app *adapter) bytesToFloat(data []byte) (*reflect.Value, []byte, error) {
	if len(data) < 1 {
		str := fmt.Sprintf(bytesLengthTooSmallErr, 1, len(data))
		return nil, nil, errors.New(str)
	}

	remaining := data[1:]
	if data[0]&ThirtyTwo != 0 {
		output, rem, err := app.bytesToFloat32(remaining)
		if err != nil {
			return nil, nil, err
		}

		value := reflect.ValueOf(output)
		return &value, rem, nil
	}

	if data[0]&SixtyFour != 0 {
		output, rem, err := app.bytesToFloat64(remaining)
		if err != nil {
			return nil, nil, err
		}

		value := reflect.ValueOf(output)
		return &value, rem, nil
	}

	str := fmt.Sprintf("the data (%v) could not be converted to a float", data)
	return nil, nil, errors.New(str)
}

func (app *adapter) bytesToFloat32(data []byte) (float32, []byte, error) {
	var value float32
	rem, err := app.bytesToCasted(data, &value)
	if err != nil {
		return 0, nil, err
	}

	return value, rem, nil
}

func (app *adapter) bytesToFloat64(data []byte) (float64, []byte, error) {
	var value float64
	rem, err := app.bytesToCasted(data, &value)
	if err != nil {
		return 0, nil, err
	}

	return value, rem, nil
}

func (app *adapter) bytesToInt(data []byte) (*reflect.Value, []byte, error) {
	if len(data) < 1 {
		str := fmt.Sprintf(bytesLengthTooSmallErr, 1, len(data))
		return nil, nil, errors.New(str)
	}

	remaining := data[1:]
	if data[0]&Height != 0 {
		output, rem, err := app.bytesToInt8(remaining)
		if err != nil {
			return nil, nil, err
		}

		value := reflect.ValueOf(output)
		return &value, rem, nil
	}

	if data[0]&Sixteen != 0 {
		output, rem, err := app.bytesToInt16(remaining)
		if err != nil {
			return nil, nil, err
		}

		value := reflect.ValueOf(output)
		return &value, rem, nil
	}

	if data[0]&ThirtyTwo != 0 {
		output, rem, err := app.bytesToInt32(remaining)
		if err != nil {
			return nil, nil, err
		}

		value := reflect.ValueOf(output)
		return &value, rem, nil
	}

	if data[0]&SixtyFour != 0 {
		output, rem, err := app.bytesToInt64(remaining)
		if err != nil {
			return nil, nil, err
		}

		value := reflect.ValueOf(output)
		return &value, rem, nil
	}

	if data[0] == 0 {
		output, rem, err := app.bytesToInt64(remaining)
		if err != nil {
			return nil, nil, err
		}

		value := reflect.ValueOf(int(output))
		return &value, rem, nil
	}

	str := fmt.Sprintf("the data (%v) could not be converted to an int", data)
	return nil, nil, errors.New(str)
}

func (app *adapter) bytesToInt8(data []byte) (int8, []byte, error) {
	var value int8
	rem, err := app.bytesToCasted(data, &value)
	if err != nil {
		return 0, nil, err
	}

	return value, rem, nil
}

func (app *adapter) bytesToInt16(data []byte) (int16, []byte, error) {
	var value int16
	rem, err := app.bytesToCasted(data, &value)
	if err != nil {
		return 0, nil, err
	}

	return value, rem, nil
}

func (app *adapter) bytesToInt32(data []byte) (int32, []byte, error) {
	var value int32
	rem, err := app.bytesToCasted(data, &value)
	if err != nil {
		return 0, nil, err
	}

	return value, rem, nil
}

func (app *adapter) bytesToInt64(data []byte) (int64, []byte, error) {
	var value int64
	rem, err := app.bytesToCasted(data, &value)
	if err != nil {
		return 0, nil, err
	}

	return value, rem, nil
}

func (app *adapter) bytesToUint(data []byte) (*reflect.Value, []byte, error) {
	if len(data) < 1 {
		str := fmt.Sprintf(bytesLengthTooSmallErr, 1, len(data))
		return nil, nil, errors.New(str)
	}

	remaining := data[1:]
	if data[0]&Height != 0 {
		output, rem, err := app.bytesToUint8(remaining)
		if err != nil {
			return nil, nil, err
		}

		value := reflect.ValueOf(output)
		return &value, rem, nil
	}

	if data[0]&Sixteen != 0 {
		output, rem, err := app.bytesToUint16(remaining)
		if err != nil {
			return nil, nil, err
		}

		value := reflect.ValueOf(output)
		return &value, rem, nil
	}

	if data[0]&ThirtyTwo != 0 {
		output, rem, err := app.bytesToUint32(remaining)
		if err != nil {
			return nil, nil, err
		}

		value := reflect.ValueOf(output)
		return &value, rem, nil
	}

	if data[0]&SixtyFour != 0 {
		output, rem, err := app.bytesToUint64(remaining)
		if err != nil {
			return nil, nil, err
		}

		value := reflect.ValueOf(output)
		return &value, rem, nil
	}

	if data[0] == 0 {
		output, rem, err := app.bytesToUint64(remaining)
		if err != nil {
			return nil, nil, err
		}

		value := reflect.ValueOf(uint(output))
		return &value, rem, nil
	}

	str := fmt.Sprintf("the data (%v) could not be converted to a uint", data)
	return nil, nil, errors.New(str)
}

func (app *adapter) bytesToUint8(data []byte) (uint8, []byte, error) {
	var value uint8
	rem, err := app.bytesToCasted(data, &value)
	if err != nil {
		return 0, nil, err
	}

	return value, rem, nil
}

func (app *adapter) bytesToUint16(data []byte) (uint16, []byte, error) {
	var value uint16
	rem, err := app.bytesToCasted(data, &value)
	if err != nil {
		return 0, nil, err
	}

	return value, rem, nil
}

func (app *adapter) bytesToUint32(data []byte) (uint32, []byte, error) {
	var value uint32
	rem, err := app.bytesToCasted(data, &value)
	if err != nil {
		return 0, nil, err
	}

	return value, rem, nil
}

func (app *adapter) bytesToUint64(data []byte) (uint64, []byte, error) {
	var value uint64
	rem, err := app.bytesToCasted(data, &value)
	if err != nil {
		return 0, nil, err
	}

	return value, rem, nil
}

func (app *adapter) bytesToString(data []byte) (*reflect.Value, []byte, error) {
	if len(data) < 8 {
		str := fmt.Sprintf(bytesLengthTooSmallErr, 8, len(data))
		return nil, nil, errors.New(str)
	}

	var strLength uint64
	strLengthBuf := bytes.NewReader(data)
	err := binary.Read(strLengthBuf, binary.LittleEndian, &strLength)
	if err != nil {
		return nil, nil, err
	}

	castedLength := int(strLength)
	if len(data) < castedLength {
		str := fmt.Sprintf(bytesLengthTooSmallErr, castedLength, len(data))
		return nil, nil, errors.New(str)
	}

	data = data[8:]
	value := reflect.ValueOf(string(data[:castedLength]))
	return &value, data[castedLength:], nil
}

func (app *adapter) bytesToArray(data []byte) (*reflect.Value, []byte, error) {
	if len(data) < 8 {
		str := fmt.Sprintf(bytesLengthTooSmallErr, 8, len(data))
		return nil, nil, errors.New(str)
	}

	var nameLength uint64
	nameLengthBuf := bytes.NewReader(data)
	err := binary.Read(nameLengthBuf, binary.LittleEndian, &nameLength)
	if err != nil {
		return nil, nil, err
	}

	castedNameLength := int(nameLength)
	if len(data) < castedNameLength {
		str := fmt.Sprintf(bytesLengthTooSmallErr, castedNameLength, len(data))
		return nil, nil, errors.New(str)
	}

	data = data[8:]
	name := string(data[:castedNameLength])
	data = data[castedNameLength:]
	if len(data) < 8 {
		str := fmt.Sprintf(bytesLengthTooSmallErr, 8, len(data))
		return nil, nil, errors.New(str)
	}

	if ptr, ok := app.mapping[name]; ok {
		if len(data) < 8 {
			str := fmt.Sprintf(bytesLengthTooSmallErr, 8, len(data))
			return nil, nil, errors.New(str)
		}

		var length uint64
		lengthBuf := bytes.NewReader(data)
		err := binary.Read(lengthBuf, binary.LittleEndian, &length)
		if err != nil {
			return nil, nil, err
		}

		ptrType := reflect.TypeOf(ptr)
		if ptrType.Kind() == reflect.Ptr {
			ptrType = ptrType.Elem()
		}

		remaining := data[8:]
		castedLength := int(length)
		slice := reflect.MakeSlice(reflect.SliceOf(ptrType), 0, 0)
		for i := 0; i < castedLength; i++ {
			instance, rem, err := app.toInstance(remaining, false)
			if err != nil {
				return nil, nil, err
			}

			remaining = rem
			slice = reflect.Append(slice, *instance)
		}

		return &slice, remaining, nil
	}

	str := fmt.Sprintf("the array type (name: %s) could not be found in the mapping", name)
	return nil, nil, errors.New(str)
}

func (app *adapter) bytesToStruct(data []byte, isPtr bool) (*reflect.Value, []byte, error) {
	if len(data) < 8 {
		str := fmt.Sprintf(bytesLengthTooSmallErr, 8, len(data))
		return nil, nil, errors.New(str)
	}

	var nameLength uint64
	nameLengthBuf := bytes.NewReader(data)
	err := binary.Read(nameLengthBuf, binary.LittleEndian, &nameLength)
	if err != nil {
		return nil, nil, err
	}

	data = data[8:]
	castedNameLength := int(nameLength)
	if len(data) < castedNameLength {
		str := fmt.Sprintf(bytesLengthTooSmallErr, castedNameLength, len(data))
		return nil, nil, errors.New(str)
	}

	name := string(data[:castedNameLength])
	if ptr, ok := app.mapping[name]; ok {
		data = data[castedNameLength:]
		if len(data) < 8 {
			str := fmt.Sprintf(bytesLengthTooSmallErr, 8, len(data))
			return nil, nil, errors.New(str)
		}

		var amount uint64
		amountBuf := bytes.NewReader(data)
		err = binary.Read(amountBuf, binary.LittleEndian, &amount)
		if err != nil {
			return nil, nil, err
		}

		insVal := reflect.New(reflect.Indirect(reflect.ValueOf(ptr)).Type()).Elem()
		casted := int(amount)
		data = data[8:]
		for i := 0; i < casted; i++ {
			field := insVal.Field(i)
			if !field.CanInterface() {
				continue
			}

			if len(data) < 8 {
				str := fmt.Sprintf(bytesLengthTooSmallErr, 8, len(data))
				return nil, nil, errors.New(str)
			}

			var fieldLength uint64
			fieldBuf := bytes.NewReader(data)
			err := binary.Read(fieldBuf, binary.LittleEndian, &fieldLength)
			if err != nil {
				return nil, nil, err
			}

			data = data[8:]
			castedFieldLength := int(fieldLength)
			if len(data) < castedFieldLength {
				str := fmt.Sprintf(bytesLengthTooSmallErr, castedFieldLength, len(data))
				return nil, nil, errors.New(str)
			}

			elementBytes := data[:castedFieldLength]
			pValue, _, err := app.toInstance(elementBytes, false)
			if err != nil {
				return nil, nil, err
			}

			value := reflect.Zero(field.Type())
			if pValue.Kind() != reflect.Invalid {
				value = pValue.Convert(field.Type())
			}

			insVal.Field(i).Set(value)
			data = data[castedFieldLength:]
		}

		if isPtr {
			name := insVal.Addr().Type().String()
			if castTo, ok := app.mapping[name]; ok {
				castToType := reflect.TypeOf(castTo)
				if castToType.Kind() == reflect.Ptr {
					castToType = castToType.Elem()
				}

				val := insVal.Addr().Convert(castToType)
				return &val, data, nil
			}

			val := insVal.Addr()
			return &val, data, nil
		}

		name := insVal.Type().String()
		if castTo, ok := app.mapping[name]; ok {
			castToType := reflect.TypeOf(castTo)
			if castToType.Kind() == reflect.Ptr {
				castToType = castToType.Elem()
			}

			val := insVal.Convert(castToType)
			return &val, data, nil
		}

		return &insVal, data, nil
	}

	str := fmt.Sprintf("the struct type (%s) does not exists in the mapping", name)
	return nil, nil, errors.New(str)
}

func (app *adapter) ptrValueToBytes(value reflect.Value) ([]byte, error) {
	output := []byte{
		Ptr,
	}

	if value.IsNil() {
		return output, nil
	}

	elem := value.Elem()
	data, err := app.valueToBytes(elem)
	if err != nil {
		return nil, err
	}

	return append(output, data...), nil
}

func (app *adapter) structToBytes(strIns interface{}) ([]byte, error) {
	strType := reflect.TypeOf(strIns)
	amount := strType.NumField()
	amountBuf := new(bytes.Buffer)
	err := binary.Write(amountBuf, binary.LittleEndian, uint64(amount))
	if err != nil {
		return nil, err
	}

	name := strType.Name()
	if strType.PkgPath() != "" {
		name = filepath.Join(strType.PkgPath(), name)
	}

	if _, ok := app.mapping[name]; !ok {
		str := fmt.Sprintf("the struct type (%s) does not exists in the mapping", name)
		return nil, errors.New(str)
	}

	nameLength := uint64(len(name))
	nameLengthBuf := new(bytes.Buffer)
	err = binary.Write(nameLengthBuf, binary.LittleEndian, nameLength)
	if err != nil {
		return nil, err
	}

	output := []byte{
		Struct,
	}

	output = append(output, nameLengthBuf.Bytes()...)
	output = append(output, []byte(name)...)
	output = append(output, amountBuf.Bytes()...)
	for i := 0; i < amount; i++ {
		value := reflect.ValueOf(strIns).Field(i)
		if !value.CanInterface() {
			continue
		}

		fieldData, err := app.valueToBytes(value)
		if err != nil {
			return nil, err
		}

		fieldLength := uint64(len(fieldData))
		lengthBuf := new(bytes.Buffer)
		err = binary.Write(lengthBuf, binary.LittleEndian, fieldLength)
		if err != nil {
			return nil, err
		}

		output = append(output, lengthBuf.Bytes()...)
		output = append(output, fieldData...)
	}

	return output, nil
}

func (app *adapter) valueToBytes(value reflect.Value) ([]byte, error) {
	kind := value.Kind()

	switch kind {
	case reflect.Bool:
		return app.boolToBytes(value.Bool()), nil
	case reflect.String:
		return app.stringToBytes(value.String())
	case reflect.Int:
		return app.intToBytes(int(value.Int()))
	case reflect.Int8:
		return app.int8ToBytes(int8(value.Int()))
	case reflect.Int16:
		return app.int16ToBytes(int16(value.Int()))
	case reflect.Int32:
		return app.int32ToBytes(int32(value.Int()))
	case reflect.Int64:
		return app.int64ToBytes(int64(value.Int()))
	case reflect.Uint:
		return app.uintToBytes(uint(value.Uint()))
	case reflect.Uint8:
		return app.uint8ToBytes(uint8(value.Uint()))
	case reflect.Uint16:
		return app.uint16ToBytes(uint16(value.Uint()))
	case reflect.Uint32:
		return app.uint32ToBytes(uint32(value.Uint()))
	case reflect.Uint64:
		return app.uint64ToBytes(uint64(value.Uint()))
	case reflect.Float32:
		return app.float32ToBytes(float32(value.Float()))
	case reflect.Float64:
		return app.float64ToBytes(float64(value.Float()))
	case reflect.Struct:
		return app.structToBytes(value.Interface())
	case reflect.Interface:
		return app.ptrValueToBytes(value)
	case reflect.Ptr:
		return app.ptrValueToBytes(value)
	case reflect.Array:
		return app.valueArrayToBytes(value)
	case reflect.Slice:
		return app.valueArrayToBytes(value)
	case reflect.Uintptr:
		return nil, errors.New("the type uintptr cannot be converted to []byte")
	case reflect.Chan:
		return nil, errors.New("the chan cannot be converted to []byte")
	case reflect.Func:
		return nil, errors.New("the func cannot be converted to []byte")
	case reflect.Complex64:
		return nil, errors.New("floating numbers with imaginary parts (complex numbers) cannot be converted to []byte")
	case reflect.Complex128:
		return nil, errors.New("floating numbers with imaginary parts (complex numbers) cannot be converted to []byte")
	}

	str := fmt.Sprintf("the given value kind (%s) is not supported", kind.String())
	return nil, errors.New(str)
}

func (app *adapter) boolToBytes(ins bool) []byte {
	if ins {
		return []byte{
			Bool, 1,
		}
	}

	return []byte{
		Bool, 0,
	}
}

func (app *adapter) stringToBytes(ins string) ([]byte, error) {
	stringBytes := []byte(ins)
	castedLength := uint64(len(stringBytes))
	lengthBuf := new(bytes.Buffer)
	err := binary.Write(lengthBuf, binary.LittleEndian, castedLength)
	if err != nil {
		return nil, err
	}

	data := []byte{
		String,
	}

	data = append(data, lengthBuf.Bytes()...)
	data = append(data, stringBytes...)
	return data, nil
}

func (app *adapter) valueArrayToBytes(value reflect.Value) ([]byte, error) {
	length := value.Len()
	castedLength := uint64(length)
	lengthBuf := new(bytes.Buffer)
	err := binary.Write(lengthBuf, binary.LittleEndian, castedLength)
	if err != nil {
		return nil, err
	}

	name := value.Type().String()
	if _, ok := app.mapping[name]; !ok {
		str := fmt.Sprintf("the array type (%s) does not exists in the mapping", name)
		return nil, errors.New(str)
	}

	nameLength := uint64(len(name))
	nameLengthBuf := new(bytes.Buffer)
	err = binary.Write(nameLengthBuf, binary.LittleEndian, nameLength)
	if err != nil {
		return nil, err
	}

	output := []byte{
		Array,
	}

	output = append(output, nameLengthBuf.Bytes()...)
	output = append(output, []byte(name)...)
	output = append(output, lengthBuf.Bytes()...)
	for i := 0; i < length; i++ {
		val := value.Index(i)
		element, err := app.valueToBytes(val)
		if err != nil {
			return nil, err
		}

		output = append(output, element...)
	}

	return output, nil
}

func (app *adapter) intToBytes(ins int) ([]byte, error) {
	return app.castedToBytes(int64(ins), []byte{
		Int,
		0,
	})
}

func (app *adapter) int8ToBytes(ins int8) ([]byte, error) {
	return app.castedToBytes(ins, []byte{
		Int,
		Height,
	})
}

func (app *adapter) int16ToBytes(ins int16) ([]byte, error) {
	return app.castedToBytes(ins, []byte{
		Int,
		Sixteen,
	})
}

func (app *adapter) int32ToBytes(ins int32) ([]byte, error) {
	return app.castedToBytes(ins, []byte{
		Int,
		ThirtyTwo,
	})
}

func (app *adapter) int64ToBytes(ins int64) ([]byte, error) {
	return app.castedToBytes(ins, []byte{
		Int,
		SixtyFour,
	})
}

func (app *adapter) uintToBytes(ins uint) ([]byte, error) {
	return app.castedToBytes(uint64(ins), []byte{
		Uint,
		0,
	})
}

func (app *adapter) uint8ToBytes(ins uint8) ([]byte, error) {
	return app.castedToBytes(ins, []byte{
		Uint,
		Height,
	})
}

func (app *adapter) uint16ToBytes(ins uint16) ([]byte, error) {
	return app.castedToBytes(ins, []byte{
		Uint,
		Sixteen,
	})
}

func (app *adapter) uint32ToBytes(ins uint32) ([]byte, error) {
	return app.castedToBytes(ins, []byte{
		Uint,
		ThirtyTwo,
	})
}

func (app *adapter) uint64ToBytes(ins uint64) ([]byte, error) {
	return app.castedToBytes(ins, []byte{
		Uint,
		SixtyFour,
	})
}

func (app *adapter) float32ToBytes(ins float32) ([]byte, error) {
	return app.castedToBytes(ins, []byte{
		Float,
		ThirtyTwo,
	})
}

func (app *adapter) float64ToBytes(ins float64) ([]byte, error) {
	return app.castedToBytes(ins, []byte{
		Float,
		SixtyFour,
	})
}

func (app *adapter) castedToBytes(ins interface{}, flags []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, ins)
	if err != nil {
		return nil, err
	}

	data := flags
	data = append(data, buf.Bytes()...)
	return data, nil
}

func (app *adapter) bytesToCasted(data []byte, ptr interface{}) ([]byte, error) {
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, ptr)
	if err != nil {
		return nil, err
	}

	unreadLength := buf.Len()
	if unreadLength <= 0 {
		return []byte{}, nil
	}

	index := len(data) - unreadLength
	return data[index:], nil
}
