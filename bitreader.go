package gobits

import (
	"bytes"
	"math"
)

// BitReader stream bits from bytes.Reader
type BitReader struct {
	reader *bytes.Reader // Underlying source of bytes
	buff   byte
	count  byte // Number of bits in buff
}

// loadbyte Read byte from reader into buff and reset count. This will drop the
// bits in buffer
func loadbyte(b *BitReader) error {
	var err error

	b.buff, err = b.reader.ReadByte()
	if err == nil {
		b.count = 8
	}
	return err
}

// NewBitStream Allocate new Bitstream, set reader, return pointer
func NewBitStream(r *bytes.Reader) *BitReader {
	ret := new(BitReader)
	ret.reader = r
	return ret
}

// Len Length, first return is bytes, second return is bits (if not byte rounded)
func (self *BitReader) Len() (int, byte) {
	return self.reader.Len(), self.count
}

// BitLen length in bits
func (self *BitReader) BitLen() uint {
	return (uint(self.reader.Len()) * 8) + uint(self.count)
}

// ReadBit Read bit from Bitstream
func (self *BitReader) ReadBit() (byte, BitError) {
	if self.count == 0 {
		if err := loadbyte(self); err != nil {
			return 0, NewBitError(err.Error(), self.count) // most graceful option for EoF
		}
	}
	ret := (self.buff & 128) >> 7
	self.buff = self.buff << 1
	self.count--
	return ret, nil
}

// ReadByte bit-aligned read
// TODO: Use custom error to return the number of valid bits
func (self *BitReader) ReadByte() (byte, BitError) {
	ret := self.buff
	bitmask := byte(math.Pow(2, float64(8-self.count))-1) << self.count
	count := self.count
	if err := loadbyte(self); err != nil {
		return ret, NewBitError(err.Error(), count) // gracefully handle EoF -> how do i get count in here?
	}

	tmp := (self.buff & bitmask) >> count
	ret += tmp
	self.buff = self.buff << (8 - count)
	self.count = count
	return ret, nil
}

// ReadBits read m bits from bitstream. returned in []byte
func (self *BitReader) ReadBits(m uint) ([]byte, BitError) {
	var tmp, bitsNeeded byte
	var err BitError
	var retSize, i uint
	var ret []byte

	if m == 0 {
		return nil, nil
	}
	retSize = uint(math.Ceil(float64(m) / 8.0))
	ret = make([]byte, retSize)
	//Read bytes
	for i = 0; i < (m / 8); i++ {
		ret[i], err = self.ReadByte()
		if err != nil {
			return ret, err // Handle EoF as gracefully as possible
		}
	}
	//Read bits
	bitsNeeded = byte(m % 8)
	for i = 0; byte(i) < bitsNeeded; i++ {
		tmp, err = self.ReadBit()
		if err != nil {
			return ret, err // Handle EoF as gracefully as possible
		}
		ret[retSize-1] += (tmp << (7 - byte(i)))
	}
	return ret, nil
}

// Equal tests wether two Bitsets are equal
// func (self *BitReader) Equal(b *BitReader) bool {
// 	return false
// }
