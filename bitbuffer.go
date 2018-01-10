package gobits

import (
	"bytes"
	"errors"
	"math"
)

//BitBuffer structure allowing storage and retrieval of individal bits.
type BitBuffer struct {
	store    []byte
	bitCount byte //bits in last byte of store
}

// NewBitBuffer returns a newly allocated BitBuffer.
func NewBitBuffer() *BitBuffer {
	ret := new(BitBuffer)
	ret.store = make([]byte, 1)
	return ret
}

// Flush resets the BitBuffer by flushing the internal store and bitcount.
func (self *BitBuffer) Flush() {
	self.store = make([]byte, 1)
	self.bitCount = 0
}

// Len returns the length.
// First value is the number of bytes, the last is the number of bits in last byte.
func (self *BitBuffer) Len() (int, byte) {
	if self.bitCount == 8 {
		return len(self.store), 0
	}
	return len(self.store) - 1, self.bitCount
}

// BitLen returns the length of the BitBuffer in bits.
// Low risk of overflow due to 64 bit unsigned int (2.3e18)
func (self *BitBuffer) BitLen() uint64 {
	//self.store should never have a len less than 1
	return uint64((len(self.store)*8)-8) + uint64(self.bitCount)
}

// grow increases the size of the internal buffer.
// func (self *BitBuffer) oldgrow() {
// 	l := len(self.store)
// 	if l == cap(self.store) {
// 		// This can be optimized by allocating in chunks through capacity
// 		// which would then use the else...
// 		tmp := make([]byte, l+1, l+1+8)
// 		copy(tmp, self.store)
// 		self.store = tmp
// 	} else {
// 		self.store = self.store[:l+1]
// 	}
// 	self.bitCount = 0
// }

func (self *BitBuffer) grow() {
	self.store = append(self.store, 0)
	self.bitCount = 0
}

// InsertBit appends bit to BitBuffer
func (self *BitBuffer) InsertBit(b byte) error {
	if b != 0 && b != 1 {
		return errors.New("Expecting bit: 0 or 1")
	}
	storePos := len(self.store) - 1
	if self.bitCount == 8 {
		self.grow()
		storePos++
	}
	bitmask := b << ((8 - self.bitCount) - 1)
	self.store[storePos] = self.store[storePos] | bitmask
	self.bitCount++
	return nil
}

// InsertByte appends byte to BitBuffer
func (self *BitBuffer) InsertByte(b byte) {
	storePos := len(self.store) - 1
	if self.bitCount == 8 {
		self.grow()
		storePos++
	}
	bitCount := self.bitCount //grow resets bitcount to 0
	if self.bitCount != 0 {
		self.grow() //we're inserting a byte, so in any case should grow
	}
	bitmask := byte(math.Pow(2, float64(bitCount)) - 1)

	//write the bits that fit
	tmp := b >> bitCount
	self.store[storePos] = self.store[storePos] | tmp

	if bitCount != 0 {
		storePos++
		//write the remaining bits
		tmp = byte((b & bitmask) << (8 - bitCount))
		self.store[storePos] = self.store[storePos] | tmp
		self.bitCount = bitCount
	} else {
		self.bitCount = 8
	}
}

// Insert multiple bytes & bits into BitBuffer
// b - bytes to insert
// bits - number of bits in the last byte
func (self *BitBuffer) Insert(b []byte, bits byte) error {
	if bits > 8 {
		return errors.New("Invalid number of bits")
	}
	if len(b) == 0 {
		return errors.New("Empty data to insert")
	}

	byteCount := len(b)
	if bits > 0 {
		byteCount--
	}
	for i := 0; i < byteCount; i++ {
		self.InsertByte(b[i])
	}
	lastPos := len(b) - 1
	for i := byte(0); i < bits; i++ {
		bitmask := byte(1 << (7 - i)) //1 << 8 = 0
		self.InsertBit((b[lastPos] & bitmask) >> (7 - i))
	}
	return nil
}

// Read read bit at position, like an array accessor
// BitBuffer stores bits from left to right.
func (self *BitBuffer) Read(position uint64) (byte, error) {
	bytePos := position / 8
	bitPos := position % 8
	bitmask := byte(1 << (7 - bitPos)) //bitpos = 0-7; 1 << 7 = 128

	if bytePos > uint64(len(self.store)-1) {
		return 0, errors.New("Position exceeds buffer length")
	}
	return (self.store[bytePos] & bitmask) >> (7 - bitPos), nil
}

func (self *BitBuffer) NewBitReader() (*BitReader, uint64) {
	ret := NewBitStream(bytes.NewReader(self.store))
	return ret, self.BitLen()
}
