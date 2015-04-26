package msgpack

import (
	"sync"
	"testing"
	"time"
)

var bufpool = sync.Pool{
	New: func() interface{} {
		return make([]byte, uint64(1000))
	},
}

// returns true if the two slices are equal
func isStringSliceEqual(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for idx := range x {
		if x[idx] != y[idx] {
			return false
		}
	}
	return true
}

// returns true if the two slices are equal
func compareInterfaceStringSlice(x, y []interface{}) bool {
	if len(x) != len(y) {
		return false
	}
	for idx := range x {
		if x[idx].(string) != y[idx].(string) {
			return false
		}
	}
	return true
}

// returns true if the two slices are equal
func compareInterfaceInt64Slice(x, y []interface{}) bool {
	if len(x) != len(y) {
		return false
	}
	for idx := range x {
		if x[idx].(int64) != y[idx].(int64) {
			return false
		}
	}
	return true
}

func TestEncodeBool(t *testing.T) {
	val := true
	bytes := bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	if done != 1 {
		t.Errorf("Encoded length should be 1 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xc3) {
		t.Errorf("True should be 0xc3 but is 0x%x", bytes[0])
	}
	bufpool.Put(bytes)
	bytes = bufpool.Get().([]byte)

	val = false
	done = Encode(val, &bytes)
	if done != 1 {
		t.Errorf("Encoded length should be 1 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xc2) {
		t.Errorf("True should be 0xc2 but is 0x%x", bytes[0])
	}
	bufpool.Put(bytes)
}

func BenchmarkEncodeBool(b *testing.B) {
	val := true
	for i := 0; i < b.N; i++ {
		bytes := bufpool.Get().([]byte)
		Encode(val, &bytes)
		bufpool.Put(bytes)
	}
}

func TestEncodeFixInt(t *testing.T) {
	var val int
	var bytes []byte

	// valid pos fixint
	val = 120
	bytes = bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	if done != 1 {
		t.Errorf("Encoded length should be 1 but is %v", len(bytes))
	}
	if bytes[0] != byte(0x78) {
		t.Errorf("Int should be 0x78 but is 0x%x", bytes[0])
	}
	_, dec := Decode(&bytes, 0)
	if dec.(int64) != 120 {
		t.Errorf("Should have been decoded as 120 but was %v", dec)
	}
	bufpool.Put(bytes)

	// valid neg fixint
	val = -20
	bytes = bufpool.Get().([]byte)
	done = Encode(val, &bytes)
	if done != 1 {
		t.Errorf("Encoded length should be 1 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xf4) {
		t.Errorf("Int should be 0xf4 but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(int64) != -20 {
		t.Errorf("Should have been decoded as -20 but was %v", dec)
	}
	bufpool.Put(bytes)
}

func BenchmarkEncodeFixInt(b *testing.B) {
	val := 120
	for i := 0; i < b.N; i++ {
		bytes := bufpool.Get().([]byte)
		Encode(val, &bytes)
		bufpool.Put(bytes)
	}
}

func TestEncodeInt(t *testing.T) {
	var val int
	var dec interface{}

	val = 32123 // int16
	bytes := bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	if done != 3 {
		t.Errorf("Encoded length should be 3 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xd1) {
		t.Errorf("Should be encoded as int16 0xd1 but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(int64) != 32123 {
		t.Errorf("Decode should be 32123 but was %v", dec)
	}
	bufpool.Put(bytes)

	val = -1234 // int16
	bytes = bufpool.Get().([]byte)
	done = Encode(val, &bytes)
	if done != 3 {
		t.Errorf("Encoded length should be 3 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xd1) {
		t.Errorf("Should be encoded as int16 0xd1 but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(int64) != -1234 {
		t.Errorf("Decode should be -1234 but was %v", dec)
	}
	bufpool.Put(bytes)
}

func TestEncodeIntBigger(t *testing.T) {
	var val int
	var dec interface{}

	val = 2147483647 // int32
	bytes := bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	if done != 5 {
		t.Errorf("Encoded length should be 5 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xd2) {
		t.Errorf("Should be encoded as int32 0xd2 but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(int64) != 2147483647 {
		t.Errorf("Decode should be 2147483647 but was %v", dec)
	}
	bufpool.Put(bytes)

	val = 4194957296 // int64
	bytes = bufpool.Get().([]byte)
	done = Encode(val, &bytes)
	if done != 9 {
		t.Errorf("Encoded length should be 9 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xd3) {
		t.Errorf("Should be encoded as int64 0xd3 but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(int64) != 4194957296 {
		t.Errorf("Decode should be 4194957296 but was %v", dec)
	}
	bufpool.Put(bytes)
}

func BenchmarkEncodeInt16(b *testing.B) {
	val := 65432
	for i := 0; i < b.N; i++ {
		bytes := bufpool.Get().([]byte)
		Encode(val, &bytes)
		bufpool.Put(bytes)
	}
}

func BenchmarkEncodeInt32(b *testing.B) {
	val := 2147483647
	for i := 0; i < b.N; i++ {
		bytes := bufpool.Get().([]byte)
		Encode(val, &bytes)
		bufpool.Put(bytes)
	}
}

func BenchmarkEncodeInt64(b *testing.B) {
	val := 4194957296
	for i := 0; i < b.N; i++ {
		bytes := bufpool.Get().([]byte)
		Encode(val, &bytes)
		bufpool.Put(bytes)
	}
}

func TestEncodeUint(t *testing.T) {
	var val uint
	var bytes []byte
	var dec interface{}

	val = 120 // uint8
	bytes = bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	if done != 2 {
		t.Errorf("Encoded length should be 2 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xcc) {
		t.Errorf("Should be encoded as uint8 0xcc but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(uint64) != 120 {
		t.Errorf("Decode should be 120 but was %v", dec)
	}
	bufpool.Put(bytes)

	val = 32123 // uint16
	bytes = bufpool.Get().([]byte)
	done = Encode(val, &bytes)
	if done != 3 {
		t.Errorf("Encoded length should be 3 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xcd) {
		t.Errorf("Should be encoded as uint16 0xcd but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(uint64) != 32123 {
		t.Errorf("Decode should be 32123 but was %v", dec)
	}
	bufpool.Put(bytes)

	val = 2147483647 // uint32
	bytes = bufpool.Get().([]byte)
	done = Encode(val, &bytes)
	if done != 5 {
		t.Errorf("Encoded length should be 5 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xce) {
		t.Errorf("Should be encoded as uint32 0xce but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(uint64) != 2147483647 {
		t.Errorf("Decode should be 2147483647 but was %v", dec)
	}
	bufpool.Put(bytes)

	val = uint(time.Now().UnixNano()) // uint64
	bytes = bufpool.Get().([]byte)
	done = Encode(val, &bytes)
	if done != 9 {
		t.Errorf("Encoded length should be 9 but is %v", len(bytes))
	}
	if bytes[0] != byte(0xcf) {
		t.Errorf("Should be encoded as uint64 0xcf but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(uint64) != uint64(val) {
		t.Errorf("Decode should be %v but was %v", val, dec)
	}
	bufpool.Put(bytes)
}

func BenchmarkEncodeUInt16(b *testing.B) {
	val := uint16(65432)
	for i := 0; i < b.N; i++ {
		bytes := bufpool.Get().([]byte)
		Encode(val, &bytes)
		bufpool.Put(bytes)
	}
}

func BenchmarkEncodeUInt32(b *testing.B) {
	val := uint32(2147483647)
	for i := 0; i < b.N; i++ {
		bytes := bufpool.Get().([]byte)
		Encode(val, &bytes)
		bufpool.Put(bytes)
	}
}

func BenchmarkEncodeUInt64(b *testing.B) {
	val := uint64(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		bytes := bufpool.Get().([]byte)
		Encode(val, &bytes)
		bufpool.Put(bytes)
	}
}

func TestEncodeString(t *testing.T) {
	var val string
	var bytes []byte
	var dec interface{}

	val = "asdf" // fixstr
	bytes = bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	if done != 1+len(val) {
		t.Errorf("Encoded length should be %v but is %v", 1+len(val), len(bytes))
	}
	if bytes[0]&0xa0 != byte(0xa0) {
		t.Errorf("Should be encoded as fixstr 0xa0 but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(string) != val {
		t.Errorf("Decode should be %v but was %v", val, dec)
	}
	bufpool.Put(bytes)

	val = "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz" // str8
	bytes = bufpool.Get().([]byte)
	done = Encode(val, &bytes)
	if done != 2+len(val) {
		t.Errorf("Encoded length should be %v but is %v", 1+len(val), len(bytes))
	}
	if bytes[0]&0xd9 != byte(0xd9) {
		t.Errorf("Should be encoded as fixstr 0xd9 but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(string) != val {
		t.Errorf("Decode should be %v but was %v", val, dec)
	}
	bufpool.Put(bytes)
}

func BenchmarkEncodeStr8(b *testing.B) {
	val := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz" // str8
	for i := 0; i < b.N; i++ {
		bytes := bufpool.Get().([]byte)
		Encode(val, &bytes)
		bufpool.Put(bytes)
	}
}

func TestEncodeFixArray(t *testing.T) {
	var bytes []byte
	var dec interface{}

	val := []interface{}{"asdf", "fdsa", "four", "gabe"}
	bytes = bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	if done != 1+4+4*len(val) { // arr len + 4 strings (1 byte len + 4 string)
		t.Errorf("Encoded length should be %v but is %v", 1+4+4*len(val), done)
	}

	if bytes[0]&0x90 != byte(0x90) {
		t.Errorf("Should be encoded as fixarray 0x90 but is 0x%x", bytes[0])
	}

	_, dec = Decode(&bytes, 0)
	if !compareInterfaceStringSlice(dec.([]interface{}), val) {
		t.Errorf("Decode should be %v but was %v", val, dec)
	}
	bufpool.Put(bytes)
}

func TestEncodeArray16(t *testing.T) {
	var bytes []byte
	var dec interface{}

	val := []interface{}{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p"}
	bytes = bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	length := 3 + 16*2 // arr prefix (3 bytes) + 16 len-1 strings
	if done != length {
		t.Errorf("Encoded length should be %v but is %v", length, done)
	}

	if bytes[0] != byte(0xdc) {
		t.Errorf("Should be encoded as arr16 0xdc but is 0x%x", bytes[0])
	}

	_, dec = Decode(&bytes, 0)
	if !compareInterfaceStringSlice(dec.([]interface{}), val) {
		t.Errorf("Decode should be %v but was %v", val, dec)
	}
	bufpool.Put(bytes)
}

func BenchmarkEncodeArrayString(b *testing.B) {
	val := []string{"asdf", "fdsa", "four", "gabe"}
	for i := 0; i < b.N; i++ {
		bytes := bufpool.Get().([]byte)
		Encode(val, &bytes)
		bufpool.Put(bytes)
	}
}

func TestEncodeArrayInts(t *testing.T) {
	var bytes []byte
	var dec interface{}
	val := []interface{}{int64(1), int64(2), int64(3)}
	bytes = bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	length := 1 + 3*1 // array len + 3 fix-bit numbers
	if done != length {
		t.Errorf("Encoded length should be %v but is %v", length, done)
	}
	_, dec = Decode(&bytes, 0)
	if !compareInterfaceInt64Slice(dec.([]interface{}), val) {
		t.Errorf("Decode should be %v but was %v", val, dec)
	}
	bufpool.Put(bytes)
}

func TestEncodeFixMap(t *testing.T) {
	var bytes []byte
	var dec interface{}

	val := map[string]interface{}{"a": 1, "b": 2, "c": 3}
	bytes = bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	length := 1 + 3*2 + 3 // 1 byte for map prefix, 3 * fixstring (2 bytes) + 3 * fixint
	if done != length {
		t.Errorf("Encoded length should be %v but is %v", length, done)
	}

	if bytes[0]&0x80 != byte(0x80) {
		t.Errorf("Should be encoded as fixmap 0x80 but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	for k, _ := range dec.(map[string]interface{}) {
		if _, found := val[k]; !found {
			t.Errorf("Decode should be %v but was %v", val, dec)
		}
	}
	bufpool.Put(bytes)
}

func TestEncodeMap16(t *testing.T) {
	var bytes []byte
	var dec interface{}

	val := map[string]interface{}{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 10, "k": 11, "l": 12, "m": 13, "n": 14, "o": 15, "p": 16}
	bytes = bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	length := 3 + 16*2 + 16 // 1 byte for map prefix, 16 * fixstring (2 bytes) + 16 * fixint
	if done != length {
		t.Errorf("Encoded length should be %v but is %v", length, done)
	}

	if bytes[0] != byte(0xde) {
		t.Errorf("Should be encoded as fixmap 0xde but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	for k, _ := range dec.(map[string]interface{}) {
		if _, found := val[k]; !found {
			t.Errorf("Decode should be %v but was %v", val, dec)
		}
	}
	bufpool.Put(bytes)
}

func BenchmarkEncodeFixMap(b *testing.B) {
	val := map[string]interface{}{"a": 1, "b": 2, "c": 3}
	for i := 0; i < b.N; i++ {
		bytes := bufpool.Get().([]byte)
		Encode(val, &bytes)
		bufpool.Put(bytes)
	}
}

func BenchmarkEncodeMap16(b *testing.B) {
	val := map[string]interface{}{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 10, "k": 11, "l": 12, "m": 13, "n": 14, "o": 15, "p": 16}
	for i := 0; i < b.N; i++ {
		bytes := bufpool.Get().([]byte)
		Encode(val, &bytes)
		bufpool.Put(bytes)
	}
}

func TestEncodeFloat32(t *testing.T) {
	var bytes []byte
	var dec interface{}
	var val float32 = 2.5
	bytes = bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	length := 5 // len of float32
	if done != length {
		t.Errorf("Encoded length should be %v but is %v", length, done)
	}

	if bytes[0] != byte(0xca) {
		t.Errorf("Should be encoded as fixmap 0xca but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(float64) != float64(val) {
		t.Errorf("Decode should be %v but was %v", val, dec)
	}
	bufpool.Put(bytes)
}

func TestEncodeFloat64(t *testing.T) {
	var bytes []byte
	var dec interface{}
	var val float64 = 2000000000000.5
	bytes = bufpool.Get().([]byte)
	done := Encode(val, &bytes)
	length := 9 // len of float64
	if done != length {
		t.Errorf("Encoded length should be %v but is %v", length, done)
	}

	if bytes[0] != byte(0xcb) {
		t.Errorf("Should be encoded as fixmap 0xcb but is 0x%x", bytes[0])
	}
	_, dec = Decode(&bytes, 0)
	if dec.(float64) != float64(val) {
		t.Errorf("Decode should be %v but was %v", val, dec)
	}
	bufpool.Put(bytes)
}
