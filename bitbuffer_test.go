package gobits

import "testing"

func TestBitBufferInsertBitOneByte(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	for i := 0; i < 8; i++ {
		b.InsertBit((tests[0] & (1 << uint(7-i))) >> uint(7-i))
	}
	if b.store[0] != tests[0] {
		t.Errorf("InsertBit error: Invalid byte. Expected %d, got %d\n", tests[0],
			b.store[0])
	}
}

func TestBitBufferInsertBitTwoBytes(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	for i := 0; i < 8; i++ {
		b.InsertBit((tests[0] & (1 << uint(7-i))) >> uint(7-i))
	}
	if b.store[0] != tests[0] {
		t.Errorf("InsertBit error: Invalid byte. Expected %d, got %d\n", tests[0],
			b.store[0])
	}

	for i := 0; i < 8; i++ {
		b.InsertBit((tests[1] & (1 << uint(7-i))) >> uint(7-i))
	}
	if (b.store[0] != tests[0]) && (b.store[1] != tests[1]) {
		t.Errorf("InsertBit Invalid byte: Expected %d, got %d\nExpected %d got %d\n",
			tests[0], b.store[0], tests[1], b.store[1])
	}
}

func TestBitBufferInsertBitError(t *testing.T) {
	b := NewBitBuffer()
	if err := b.InsertBit(2); err == nil {
		t.Errorf("Error: InsertBit should only handle values 0 & 1\n")
	}
}

func TestBitBufferInsertByte(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.InsertByte(tests[0])
	if b.store[0] != tests[0] {
		t.Errorf("InsertByte Invalid byte: Expected %d, got %d\n", tests[0], b.store[0])
	}
}

func TestBitBufferInsertByteTwoBytes(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.InsertByte(tests[0])
	b.InsertByte(tests[1])
	if (b.store[0] != tests[0]) && (b.store[1] != tests[1]) {
		t.Errorf("InsertByte Invalid byte: Expected %d, got %d\nExpected %d, got %d\n",
			tests[0], b.store[0], tests[1], b.store[1])
	}
}

func TestBitBufferInsertByteAlignment(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.InsertBit(0)
	b.InsertBit(1)
	b.InsertBit(1)
	b.InsertBit(0)
	b.InsertByte(tests[0])
	//0110 1100 = 108
	if b.store[0] != 108 && b.store[1] != 0 {
		t.Errorf("InsertByte Invalid byte: Expected %d, got %d\nExpected %d, got %d\n",
			108, b.store[0], 0, b.store[1])
	}
}

func TestBitBufferInsertByteAlignmentGrowth(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	for i := 0; i < 8; i++ {
		b.InsertBit((tests[0] & (1 << uint(7-i))) >> uint(7-i))
	}
	b.InsertByte(tests[2])
	if b.store[0] != tests[0] && b.store[1] != tests[2] {
		t.Errorf("InsertByte Invalid byte: Expected %d, got %d\nExpected %d, got %d\n",
			tests[0], b.store[0], tests[2], b.store[1])
	}
}

func TestBitBufferInsertOneByte(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.Insert(tests[:1], 8)
	if b.store[0] != tests[0] {
		t.Errorf("Insert Invalid byte: Expected %d, got %d\n",
			tests[0], b.store[0])
	}
}

func TestBitBufferInsertOneByteFourBits(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.Insert(tests[:], 4)
	if b.store[0] != tests[0] && b.store[1] != tests[1] &&
		b.store[2] != 144 {
		t.Errorf("Insert Invalid byte: Expected %d, got %d\nExpected %d, got %d\nExpected %d, got %d\n",
			tests[0], b.store[0], tests[1], b.store[1], 144, b.store[2])
	}
}

func TestBitBufferInsertInvalidBits(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	err := b.Insert(tests[:], 9)
	if err == nil {
		t.Errorf("Invalid number of bits should result in error\n")
	}
}

func TestBitBufferInsertSliceError(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	err := b.Insert(nil, 0)
	if err == nil {
		t.Errorf("Nil byte array should be handled")
	}
	err = b.Insert(tests[:0], 0)
	if err == nil {
		t.Errorf("Empty byte array should be handled")
	}
}

func TestBitBufferRead(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.Insert(tests[:], 8)
	for i := 0; i < len(tests); i++ {
		for j := 0; j < 8; j++ {
			bitmask := byte(1 << byte(7-j))
			store, err := b.Read(uint64(j + (8 * i)))
			test := byte(tests[i]&bitmask) >> byte(7-j)
			if store != test || err != nil {
				t.Errorf("Read Error: Expected %d, Got %d\n", test, store)
			}
		}
	}
}

func TestBitBufferBitLen(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.InsertBit(1)
	if test := b.BitLen(); test != 1 {
		t.Errorf("BitLen error: Expected %d, Got %d\n", 1, test)
	}

	b.InsertByte(192)
	if test := b.BitLen(); test != 9 {
		t.Errorf("BitLen error: Expected %d, Got %d\n", 9, test)
	}

	b.Insert(tests[1:2], 4)
	if test := b.BitLen(); test != 13 {
		t.Errorf("BitLen error: Expected %d, Got %d\n", 13, test)
	}
}

func TestBitBufferFlush(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.Insert(tests[:], 8)
	b.Flush()
	if len(b.store) != 1 && b.store[0] != 0 && b.bitCount != 0 {
		t.Errorf("Flush error: BitBuffer not properly flushed.")
	}
}

func TestBitBufferLen(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	bytes, bits := b.Len()
	if bytes != 0 || bits != 0 {
		t.Errorf("Len error: Expected (%d,%d) got (%d,%d)\n", 0, 0, bytes, bits)
	}

	b.Insert(tests[:], 4)
	bytes, bits = b.Len()
	if bytes != 2 || bits != 4 {
		t.Errorf("Len error: Expected (%d,%d) got (%d,%d)\n", 2, 4, bytes, bits)
	}

	b.Flush()
	bytes, bits = b.Len()
	if bytes != 0 && bits != 0 {
		t.Errorf("Len error: Expected (%d,%d) got (%d,%d)\n", 0, 0, bytes, bits)
	}

	for i := 0; i < 8; i++ {
		b.InsertBit(1)
	}
	bytes, bits = b.Len()
	if bytes != 1 && bits != 0 {
		t.Errorf("Len error: Expected (%d,%d) got (%d,%d)\n", 1, 0, bytes, bits)
	}
}
