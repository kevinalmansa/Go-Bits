package gobits

//BitError custom error interface to include bit count of returned value
//accompanying error
type BitError interface {
	Error() string
	BitCount() byte
}

type customBitError struct {
	bitCount byte
	s        string
}

// NewBitError creates a new BitError, which is compatible with error interface
func NewBitError(s string, bitCount byte) BitError {
	err := customBitError{s: s, bitCount: bitCount}
	return err
}

func (err customBitError) Error() string {
	return err.s
}

func (err customBitError) BitCount() byte {
	return err.bitCount
}
