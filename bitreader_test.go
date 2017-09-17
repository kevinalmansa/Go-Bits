package gobits

import (
	"bytes"
	"testing"
)

func TestBitReaderReadBit(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()
	b.Insert(tests[:], 8)

	r := NewBitStream(bytes.NewReader(b.store))
	for i := 0; i < len(tests); i++ {
		t.Logf("Byte: %d\n", i)
		for j := 0; j < 8; j++ {
			val, bErr := b.Read(uint64((i * 8) + j))
			test, err := r.ReadBit()
			if err != nil || bErr != nil {
				t.Errorf("Readbit Error: Error recieved from function")
			}
			if val != test {
				t.Errorf("ReadBit Error: Expected %d, got %d\n", val, test)
			}
			if j == 7 {
				t.Logf("All 8 Bits Retrieved")
			}
		}
	}

}

func TestBitReaderReadByte(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.Insert(tests[:], 8)

	r := NewBitStream(bytes.NewReader(b.store))
	for i := 0; i < len(tests); i++ {
		test, err := r.ReadByte()
		if test != tests[i] {
			t.Errorf("ReadByte Error: Expected %d, got %d\n", tests[i], test)
		}
		if err != nil {
			t.Errorf("ReadByte Error: Unexpected error.\n")
		}
	}
	if _, err := r.ReadBit(); err == nil {
		t.Errorf("ReadBit Error: Expected EoF\n")
	}
}

func TestBitReaderReadByteAlignment1(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.InsertBit(1)
	b.Insert(tests[:], 8)

	r := NewBitStream(bytes.NewReader(b.store))
	bit, _ := r.ReadBit()
	if bit != 1 {
		t.Errorf("ReadByte Error: Expected %d, got %d\n", 1, bit)
	}
	for i := 0; i < len(tests); i++ {
		test, err := r.ReadByte()
		if test != tests[i] {
			t.Errorf("ReadByte Error: Expected %d, got %d\n", tests[i], test)
		}
		if err != nil {
			t.Errorf("ReadByte Error: Unexpected error.\n")
		}
	}
}

func TestBitReaderReadByteAlignment2(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.InsertBit(1)
	b.InsertBit(0)
	b.Insert(tests[:], 8)

	r := NewBitStream(bytes.NewReader(b.store))
	bit, _ := r.ReadBit()
	if bit != 1 {
		t.Errorf("ReadByte Error: Expected %d, got %d\n", 1, bit)
	}
	bit, _ = r.ReadBit()
	if bit != 0 {
		t.Errorf("ReadByte Error: Expected %d, got %d\n", 0, bit)
	}

	for i := 0; i < len(tests); i++ {
		test, err := r.ReadByte()
		if test != tests[i] {
			t.Errorf("ReadByte Error: Expected %d, got %d\n", tests[i], test)
		}
		if err != nil {
			t.Errorf("ReadByte Error: Unexpected error.\n")
		}
	}
}

func TestBitReaderReadBitsBytes(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.Insert(tests[:], 8)

	r := NewBitStream(bytes.NewReader(b.store))
	test, _ := r.ReadBits(24)

	if len(test) != 3 {
		t.Errorf("ReadBits Error: Expected Length: %d, got %d\n", 3, len(test))
	}

	for i := 0; i < len(tests); i++ {
		if test[i] != tests[i] {
			t.Errorf("ReadBits Error: Expected %d, got %d\n", tests[i], test[i])
		}
	}
}

func TestBitReaderReadBitsBits(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.Insert(tests[0:1], 3)

	r := NewBitStream(bytes.NewReader(b.store))
	test, _ := r.ReadBits(3)

	if len(test) != 1 {
		t.Errorf("ReadBits Error: Expected Length: %d, got %d\n", 1, len(test))
	}
	for i := 0; i < len(test); i++ {
		if test[i] != tests[i] {
			t.Errorf("ReadBits Error: Expected %d, got %d\n", tests[i], test[i])
		}
	}
}

func TestBitReaderReadBitsError(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()
	b.Insert(tests[:], 8)

	r := NewBitStream(bytes.NewReader(b.store[0:0]))
	test, err := r.ReadBits(0)
	if test != nil || err != nil {
		t.Errorf("ReadBits Error: Expected (nil, nil).\n")
	}

	test, err = r.ReadBits(8)
	if test[0] != 0 || err == nil {
		t.Errorf("ReadBits Error: Expected ({0}, EoF), Got ({%d}, %s).\n", test[0], err.Error())
	}

	test, err = r.ReadBits(1)
	if test[0] != 0 || err == nil {
		t.Errorf("ReadBits Error: Expected ({0}, EoF), Got ({%d}, %s).\n", test[0], err.Error())
	}

	t.Logf("All Errors Properly Handled: bitcount %d, err %s", err.BitCount(), err.Error())
}

func TestBitReaderLen(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.Insert(tests[:], 8)
	size := len(b.store)
	r := NewBitStream(bytes.NewReader(b.store))

	if tbyte, tbit := r.Len(); tbyte != 3 || tbit != 0 {
		t.Errorf("Len Error: Expected (%d, %d), got (%d,%d)\nStore Size: %d\n",
			3, 0, tbyte, tbit, size)
	}
}

func TestBitReaderLenBit(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.Insert(tests[:], 8)

	r := NewBitStream(bytes.NewReader(b.store))
	_, _ = r.ReadBits(3)
	tbyte, tbit := r.Len()
	if tbyte != 2 || tbit != 5 {
		t.Errorf("Len Error: Expected (%d, %d), got (%d,%d)\n", 2, 5, tbyte, tbit)
	}
}

func TestBitLen(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.Insert(tests[:], 8)
	r := NewBitStream(bytes.NewReader(b.store))

	if r.BitLen() != uint(8*len(tests)) {
		t.Errorf("BitLen Error: Expected %d, got %d\n", uint(8*len(tests)),
			r.BitLen())
	}
}

func TestBitLenBits(t *testing.T) {
	tests := [3]byte{192, 39, 156}
	b := NewBitBuffer()

	b.Insert(tests[:], 8)
	r := NewBitStream(bytes.NewReader(b.store))

	_, _ = r.ReadBits(3)

	test := r.BitLen()
	if test != (uint(8*len(tests)) - (8 - 5)) {
		t.Errorf("BitLen Error: Expected %d, got %d\n", (uint(8*len(tests)) - (8 - 5)),
			test)
	}
}
