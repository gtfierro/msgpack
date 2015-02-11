package msgpack

import (
	"testing"
	"time"
)

func TestEncodeBool(t *testing.T) {
	val := true
	bytes := Encode(val)
	if len(bytes) != 1 {
		t.Errorf("Encoded length should be 1 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xc3) {
		t.Errorf("True should be 0xc3 but is %v", bytes[0])
	}

	val = false
	bytes = Encode(val)
	if len(bytes) != 1 {
		t.Errorf("Encoded length should be 1 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xc2) {
		t.Errorf("True should be 0xc2 but is %v", bytes[0])
	}
}

func BenchmarkEncodeBool(b *testing.B) {
	val := true
	for i := 0; i < b.N; i++ {
		Encode(val)
	}
}

func TestEncodeFixInt(t *testing.T) {
	var val int
	var bytes []byte

	// valid pos fixint
	val = 120
	bytes = Encode(val)
	if len(bytes) != 1 {
		t.Errorf("Encoded length should be 1 but is %v", len(bytes))
	}
	if bytes[0] != byte(0x78) {
		t.Errorf("Int should be 0x78 but is %v", bytes[0])
	}
	_, dec := Decode(&bytes, 0)
	if dec.(int64) != 120 {
		t.Errorf("Should have been decoded as 120 but was %v", dec)
	}

	// valid neg fixint
	val = -20
	bytes = Encode(val)
	if len(bytes) != 1 {
		t.Errorf("Encoded length should be 1 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xf4) {
		t.Errorf("Int should be 0xf4 but is %v", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(int64) != -20 {
		t.Errorf("Should have been decoded as -20 but was %v", dec)
	}
}

func BenchmarkEncodeFixInt(b *testing.B) {
	val := 120
	for i := 0; i < b.N; i++ {
		Encode(val)
	}
}

func TestEncodeInt(t *testing.T) {
	var val int
	var bytes []byte
	var dec interface{}

	val = 32123 // int16
	bytes = Encode(val)
	if len(bytes) != 3 {
		t.Errorf("Encoded length should be 3 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xd1) {
		t.Errorf("Should be encoded as int16 0xd1 but is %v", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(int64) != 32123 {
		t.Errorf("Decode should be 32123 but was %v", dec)
	}

	val = -1234 // int16
	bytes = Encode(val)
	if len(bytes) != 3 {
		t.Errorf("Encoded length should be 3 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xd1) {
		t.Errorf("Should be encoded as int16 0xd1 but is %v", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(int64) != -1234 {
		t.Errorf("Decode should be -1234 but was %v", dec)
	}
}

func TestEncodeIntBigger(t *testing.T) {
	var val int
	var bytes []byte
	var dec interface{}

	val = 2147483647 // int32
	bytes = Encode(val)
	if len(bytes) != 5 {
		t.Errorf("Encoded length should be 5 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xd2) {
		t.Errorf("Should be encoded as int32 0xd2 but is %v", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(int64) != 2147483647 {
		t.Errorf("Decode should be 2147483647 but was %v", dec)
	}

	val = 4194957296 // int64
	bytes = Encode(val)
	if len(bytes) != 9 {
		t.Errorf("Encoded length should be 9 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xd3) {
		t.Errorf("Should be encoded as int64 0xd3 but is %v", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(int64) != 4194957296 {
		t.Errorf("Decode should be 4194957296 but was %v", dec)
	}
}

func BenchmarkEncodeInt16(b *testing.B) {
	val := 65432
	for i := 0; i < b.N; i++ {
		Encode(val)
	}
}

func BenchmarkEncodeInt32(b *testing.B) {
	val := 2147483647
	for i := 0; i < b.N; i++ {
		Encode(val)
	}
}

func BenchmarkEncodeInt64(b *testing.B) {
	val := 4194957296
	for i := 0; i < b.N; i++ {
		Encode(val)
	}
}

func TestEncodeUint(t *testing.T) {
	var val uint
	var bytes []byte
	var dec interface{}

	val = 120 // uint8
	bytes = Encode(val)
	if len(bytes) != 2 {
		t.Errorf("Encoded length should be 2 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xcc) {
		t.Errorf("Should be encoded as uint8 0xcc but is %v", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(uint64) != 120 {
		t.Errorf("Decode should be 120 but was %v", dec)
	}

	val = 32123 // uint16
	bytes = Encode(val)
	if len(bytes) != 3 {
		t.Errorf("Encoded length should be 3 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xcd) {
		t.Errorf("Should be encoded as uint16 0xcd but is %v", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(uint64) != 32123 {
		t.Errorf("Decode should be 32123 but was %v", dec)
	}

	val = 2147483647 // uint32
	bytes = Encode(val)
	if len(bytes) != 5 {
		t.Errorf("Encoded length should be 5 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xce) {
		t.Errorf("Should be encoded as uint32 0xce but is %v", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(uint64) != 2147483647 {
		t.Errorf("Decode should be 2147483647 but was %v", dec)
	}

	val = uint(time.Now().UnixNano()) // uint64
	bytes = Encode(val)
	if len(bytes) != 9 {
		t.Errorf("Encoded length should be 9 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xcf) {
		t.Errorf("Should be encoded as uint64 0xcf but is %v", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(uint64) != uint64(val) {
		t.Errorf("Decode should be %v but was %v", val, dec)
	}
}

func BenchmarkEncodeUInt16(b *testing.B) {
	val := uint16(65432)
	for i := 0; i < b.N; i++ {
		Encode(val)
	}
}

func BenchmarkEncodeUInt32(b *testing.B) {
	val := uint32(2147483647)
	for i := 0; i < b.N; i++ {
		Encode(val)
	}
}

func BenchmarkEncodeUInt64(b *testing.B) {
	val := uint64(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		Encode(val)
	}
}

func TestEncodeString(t *testing.T) {
	var val string
	var bytes []byte
	var dec interface{}

	val = "asdf" // fixstr
	bytes = Encode(val)
	if len(bytes) != 1+len(val) {
		t.Errorf("Encoded length should be %v but is %v", 1+len(val), len(bytes))
	}
	if bytes[0] & 0xa0 != byte(0xa0) {
		t.Errorf("Should be encoded as fixstr 0xa0 but is %v", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(string) != val {
		t.Errorf("Decode should be %v but was %v", val, dec)
	}

	val = "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz" // str8
	bytes = Encode(val)
	if len(bytes) != 2+len(val) {
		t.Errorf("Encoded length should be %v but is %v", 1+len(val), len(bytes))
	}
	if bytes[0] & 0xd9 != byte(0xd9) {
		t.Errorf("Should be encoded as fixstr 0xd9 but is %v", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(string) != val {
		t.Errorf("Decode should be %v but was %v", val, dec)
	}
}

func BenchmarkEncodeStr8(b *testing.B) {
	val := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz" // str8
	for i := 0; i < b.N; i++ {
		Encode(val)
	}
}

