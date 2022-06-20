package bytes

const (
	// Bool represents the bool flag
	Bool uint8 = 1 << iota

	// Int represents the int flag
	Int

	// Uint represents the uint flag
	Uint

	// Float represents the float flag
	Float

	// Struct represents the struct flag
	Struct

	// Array represents the array flag
	Array

	// Ptr represents the ptr flag
	Ptr

	// String represents the string
	String
)

const (
	// Height represents the 8 flag
	Height uint8 = 1 << iota

	// Sixteen represents the 16 flag
	Sixteen

	// ThirtyTwo represents the 32 flag
	ThirtyTwo

	// SixtyFour represents the 64 flag
	SixtyFour
)

const bytesLengthTooSmallErr = "the []byte was expected to contain at least %d elements, %d returned"

// NewAdapterBuilder creates a new adapter builder
func NewAdapterBuilder() AdapterBuilder {
	return createAdapterBuilder()
}

// AdapterBuilder represents an adapter builder
type AdapterBuilder interface {
	Create() AdapterBuilder
	WithMapping(mapping map[string]interface{}) AdapterBuilder
	Now() (Adapter, error)
}

// Adapter represents the bytes adapter
type Adapter interface {
	ToBytes(ins interface{}) ([]byte, error)
	ToInstance(bytes []byte) (interface{}, []byte, error)
}
